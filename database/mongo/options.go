package mongo

import (
	"context"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Paginatior struct {
	CurrentPage int64 `json:"current_page,omitempty"`
	LastPage    int64 `json:"last_page,omitempty"`
	PerPage     int64 `json:"per_page,omitempty"`
	Total       int64 `json:"total,omitempty"`
}

type Pagination struct {
	options     *options.FindOptions
	coll        *mongoDriver.Collection
	currentPage int64
	perPage     int64
	filter      interface{}
}

var (
	DEFAULT_PAGE  int64 = 1
	DEFAULT_LIMIT int64 = 10
)

func NewPagination(coll *mongoDriver.Collection) *Pagination {
	return &Pagination{
		coll:        coll,
		options:     options.Find(),
		currentPage: DEFAULT_PAGE,
		perPage:     DEFAULT_LIMIT,
	}
}

func (p *Pagination) Where(filter interface{}) *Pagination {
	p.filter = filter
	return p
}

func (p *Pagination) Page(page *int64) *Pagination {
	if page != nil {
		p.currentPage = *page
	}

	return p
}

func (p *Pagination) Limit(limit *int64) *Pagination {
	if limit != nil {
		p.perPage = *limit
	}

	return p
}

func (p *Pagination) Desc(key string) *Pagination {
	p.options.SetSort(bson.D{{Key: key, Value: -1}})
	return p
}

func (p *Pagination) Asc(key string) *Pagination {
	p.options.SetSort(bson.D{{Key: key, Value: 1}})

	return p
}

func (p *Pagination) Find(ctx context.Context, results interface{}) (*Paginatior, error) {
	skip := (p.currentPage - 1) * p.perPage

	cur, err := p.coll.Find(ctx, p.filter, p.options.SetSkip(skip).SetLimit(p.perPage))
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, results); err != nil {
		return nil, err
	}

	count, err := p.count(ctx)
	if err != nil {
		return nil, err
	}

	np := &Paginatior{
		CurrentPage: p.currentPage,
		PerPage:     p.perPage,
		LastPage:    int64(math.Ceil(float64(*count) / float64(p.perPage))),
		Total:       *count,
	}

	return np, nil
}

func (p *Pagination) count(ctx context.Context) (*int64, error) {
	opts := options.Count().SetHint("_id_")

	count, err := p.coll.CountDocuments(ctx, p.filter, opts)
	if err != nil {
		return nil, err
	}

	return &count, err
}
