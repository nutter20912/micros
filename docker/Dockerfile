FROM golang:1.20-alpine

WORKDIR /app

###########################################################################
# local
###########################################################################
RUN apk --update --no-cache add protobuf make curl git

# install go-micro tools
RUN go install github.com/go-micro/cli/cmd/go-micro@v1.1.4

# ! CANNOT USE IN DOCKER
# install database migrate tools
#RUN curl -sSf https://atlasgo.sh | sh
#RUN go install ariga.io/atlas-provider-gorm \
#  && go get -u ariga.io/atlas-provider-gorm/gormschema \
#  && go get -u ariga.io/atlas-go-sdk/recordriver

# debug tool - dlv
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.* ./

RUN go mod download

COPY . /app
