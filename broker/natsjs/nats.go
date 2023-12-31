// Package nats provides a NATS broker
package natsjs

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"

	"dario.cat/mergo"
	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/codec/json"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/util/cmd"
)

func init() {
	cmd.DefaultBrokers["natsjs"] = NewBroker
}

type natsBroker struct {
	sync.Once
	sync.RWMutex

	// indicate if we're connected
	connected bool

	addrs []string
	conn  *nats.Conn
	js    jetstream.JetStream
	opts  broker.Options
	nopts nats.Options

	// should we drain the connection
	drain   bool
	closeCh chan (error)
}

type subscriber struct {
	opts broker.SubscribeOptions

	s  jetstream.Stream
	c  jetstream.Consumer
	cc jetstream.ConsumeContext
}

type publication struct {
	t   string
	err error
	m   *broker.Message
}

func (p *publication) Topic() string {
	return p.t
}

func (p *publication) Message() *broker.Message {
	return p.m
}

func (p *publication) Ack() error {
	// nats does not support acking
	return nil
}

func (p *publication) Error() error {
	log.Printf("[natsjs] %v", p.err)
	return p.err
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.c.CachedInfo().Stream
}

func (s *subscriber) Unsubscribe() error {
	s.cc.Stop()

	log.Println("[natsjs] Unsubscribe")

	return nil
}

func (n *natsBroker) Address() string {
	if n.conn != nil && n.conn.IsConnected() {
		return n.conn.ConnectedUrl()
	}

	if len(n.addrs) > 0 {
		return n.addrs[0]
	}

	return ""
}

func (n *natsBroker) setAddrs(addrs []string) []string {
	//nolint:prealloc
	var cAddrs []string
	for _, addr := range addrs {
		if len(addr) == 0 {
			continue
		}
		if !strings.HasPrefix(addr, "nats://") {
			addr = "nats://" + addr
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{nats.DefaultURL}
	}
	return cAddrs
}

func (n *natsBroker) Connect() error {
	n.Lock()
	defer n.Unlock()

	if n.connected {
		return nil
	}

	status := nats.CLOSED
	if n.conn != nil {
		status = n.conn.Status()
	}

	switch status {
	case nats.CONNECTED, nats.RECONNECTING, nats.CONNECTING:
		n.connected = true
		return nil
	default: // DISCONNECTED or CLOSED or DRAINING
		opts := n.nopts
		opts.Servers = n.addrs
		opts.Secure = n.opts.Secure
		opts.TLSConfig = n.opts.TLSConfig

		// secure might not be set
		if n.opts.TLSConfig != nil {
			opts.Secure = true
		}

		c, err := opts.Connect()
		if err != nil {
			return err
		}
		n.conn = c
		n.connected = true

		js, _ := jetstream.New(c)
		n.js = js

		return nil
	}
}

func (n *natsBroker) Disconnect() error {
	n.Lock()
	defer n.Unlock()

	// drain the connection if specified
	if n.drain {
		n.conn.Drain()
		n.closeCh <- nil
	}

	// close the client connection
	n.conn.Close()
	log.Println("[natsjs] Disconnect")

	// set not connected
	n.connected = false

	return nil
}

func (n *natsBroker) Init(opts ...broker.Option) error {
	n.setOption(opts...)
	return nil
}

func (n *natsBroker) Options() broker.Options {
	return n.opts
}

func (n *natsBroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	n.RLock()
	defer n.RUnlock()

	if n.js == nil {
		return errors.New("not connected")
	}

	b, err := n.opts.Codec.Marshal(msg)
	if err != nil {
		return err
	}

	clientOpt := &broker.PublishOptions{}

	for _, o := range opts {
		o(clientOpt)
	}

	pubOpts := []jetstream.PublishOpt{
		jetstream.WithMsgID(msg.Header["Micro-Id"]),
	}

	if val, ok := clientOpt.Context.Value(publishOptionsKey{}).([]jetstream.PublishOpt); ok {
		pubOpts = append(pubOpts, val...)
	}

	_, err = n.js.Publish(context.Background(), topic, b, pubOpts...)

	return err
}

func (n *natsBroker) stream(topic string, opt broker.SubscribeOptions) (jetstream.Stream, error) {
	name := strings.Split(topic, ".")[0]
	//jsCfg := jetstream.StreamConfig{Subjects: []string{topic}}

	//if val, ok := opt.Context.Value(streamConfigKey{}).(jetstream.StreamConfig); ok {
	//	mergo.Merge(&jsCfg, val)
	//}

	//if _, err := n.js.CreateStream(context.Background(), jsCfg); err != nil && !errors.Is(jetstream.ErrStreamNameAlreadyInUse, err) {
	//	return nil, err
	//}

	s, err := n.js.Stream(context.Background(), name)
	if err != nil {
		return nil, err
	}

	return s, nil

}

func (n *natsBroker) consumer(s jetstream.Stream, topic string, opt broker.SubscribeOptions) (jetstream.Consumer, error) {
	consumerConfig := jetstream.ConsumerConfig{FilterSubject: topic}

	if val, ok := opt.Context.Value(consumerConfigKey{}).(jetstream.ConsumerConfig); ok {
		mergo.Merge(&consumerConfig, val)
	}

	c, err := s.CreateOrUpdateConsumer(context.Background(), consumerConfig)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (n *natsBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	n.RLock()
	if n.conn == nil {
		n.RUnlock()
		return nil, errors.New("not connected")
	}
	n.RUnlock()

	opt := broker.SubscribeOptions{
		AutoAck: true,
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&opt)
	}

	fn := func(msg jetstream.Msg) {
		msg.Ack()

		var m broker.Message
		pub := &publication{t: msg.Subject()}
		eh := n.opts.ErrorHandler
		err := n.opts.Codec.Unmarshal(msg.Data(), &m)
		pub.err = err
		pub.m = &m

		if err != nil {
			m.Body = msg.Data()
			n.opts.Logger.Log(logger.ErrorLevel, err)
			if eh != nil {
				eh(pub)
			}
			return
		}
		if err := handler(pub); err != nil {
			pub.err = err
			n.opts.Logger.Log(logger.ErrorLevel, err)
			if eh != nil {
				eh(pub)
			}
		}
	}

	var err error

	n.RLock()
	s, err := n.stream(topic, opt)
	n.RUnlock()
	if err != nil {
		return nil, err
	}

	n.RLock()
	c, err := n.consumer(s, topic, opt)
	n.RUnlock()
	if err != nil {
		return nil, err
	}

	n.RLock()
	cc, err := c.Consume(fn)
	n.RUnlock()
	if err != nil {
		return nil, err
	}

	return &subscriber{opts: opt, s: s, c: c, cc: cc}, nil
}

func (n *natsBroker) String() string {
	return "natsjs"
}

func (n *natsBroker) setOption(opts ...broker.Option) {
	for _, o := range opts {
		o(&n.opts)
	}

	n.Once.Do(func() {
		n.nopts = nats.GetDefaultOptions()
	})

	if nopts, ok := n.opts.Context.Value(optionsKey{}).(nats.Options); ok {
		n.nopts = nopts
	}

	// broker.Options have higher priority than nats.Options
	// only if Addrs, Secure or TLSConfig were not set through a broker.Option
	// we read them from nats.Option
	if len(n.opts.Addrs) == 0 {
		n.opts.Addrs = n.nopts.Servers
	}

	if !n.opts.Secure {
		n.opts.Secure = n.nopts.Secure
	}

	if n.opts.TLSConfig == nil {
		n.opts.TLSConfig = n.nopts.TLSConfig
	}
	n.addrs = n.setAddrs(n.opts.Addrs)

	if n.opts.Context.Value(drainConnectionKey{}) != nil {
		n.drain = true
		n.closeCh = make(chan error)
		n.nopts.ClosedCB = n.onClose
		n.nopts.AsyncErrorCB = n.onAsyncError
		n.nopts.DisconnectedErrCB = n.onDisconnectedError
	}
}

func (n *natsBroker) onClose(conn *nats.Conn) {
	n.closeCh <- nil

	log.Println("[natsjs] onClose")
}

func (n *natsBroker) onAsyncError(conn *nats.Conn, sub *nats.Subscription, err error) {
	// There are kinds of different async error nats might callback, but we are interested
	// in ErrDrainTimeout only here.
	if err == nats.ErrDrainTimeout {
		n.closeCh <- err
	}
}

func (n *natsBroker) onDisconnectedError(conn *nats.Conn, err error) {
	n.closeCh <- err
}

func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		// Default codec
		Codec:    json.Marshaler{},
		Context:  context.Background(),
		Registry: registry.DefaultRegistry,
		Logger:   logger.DefaultLogger,
	}

	n := &natsBroker{
		opts: options,
	}
	n.setOption(opts...)

	return n
}
