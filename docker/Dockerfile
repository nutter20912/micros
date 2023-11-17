FROM golang:1.20-alpine AS base

WORKDIR /app

COPY go.* ./

RUN go mod download

###########################################################################
# local
###########################################################################
FROM base AS local

RUN apk --update --no-cache add protobuf make curl git

# install go-micro tools
#RUN go install github.com/go-micro/cli/cmd/go-micro@v1.1.4

# ! CANNOT USE IN DOCKER
# install database migrate tools
#RUN curl -sSf https://atlasgo.sh | sh
#RUN go install ariga.io/atlas-provider-gorm \
#  && go get -u ariga.io/atlas-provider-gorm/gormschema \
#  && go get -u ariga.io/atlas-go-sdk/recordriver

# debug tool - dlv
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY . /app

ENTRYPOINT ["sleep", "infinity"]

###########################################################################
# build image
###########################################################################
FROM base AS builder

ARG APP_NAME

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux \
    go build ./app/${APP_NAME}

FROM scratch AS image

ARG APP_NAME

WORKDIR /app

COPY --from=builder /app/${APP_NAME} ./main

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/.env .
COPY --from=builder /app/config ./config
COPY --from=builder /app/database/migrations/ ./database/migrations

CMD ["./main"]
