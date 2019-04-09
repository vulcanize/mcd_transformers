FROM golang:alpine as builder
RUN apk --update --no-cache add make git g++

# Build statically linked vDB binary
RUN go get -u -d github.com/vulcanize/vulcanizedb
WORKDIR /go/src/github.com/vulcanize/vulcanizedb
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' .

# Build migration tool. Locked to 2.6.0 until #158 is fixed
RUN go get -u -d github.com/pressly/goose/ github.com/lib/pq
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN cd ../.. && git checkout v2.6.0 && go get
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -tags='no_mysql no_sqlite' -o goose


# Second stage
FROM alpine
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/vulcanizedb /app/
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose /app/goose
# COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/db/migrations/* /app/migrations/

ADD ./dockerfiles/startup_script.sh /app/
ADD ./environments/staging.toml /app/environments/
# Collision between core and plugin?
ADD ./db/migrations/* /app/migrations/

WORKDIR /app
CMD ["./startup_script.sh"]
