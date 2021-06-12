FROM golang:1.16-alpine as builder

# To fix go get and build with cgo
RUN apk add --no-cache --virtual .build-deps \
    bash \
    gcc \
    musl-dev \
    ca-certificates

# Install SSL ca certificates.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

RUN mkdir build
COPY . /build
WORKDIR /build

RUN go mod download

ARG GIT_TAG=dev
ARG GIT_COMMIT=none
ARG RDATE=unknown

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-X main.version=$GIT_TAG -X main.commit=$GIT_COMMIT -X main.date=$RDATE -extldflags '-static'" \
    -o event-indexer .
RUN adduser -S -D -H -h /build event-indexer

USER event-indexer

FROM scratch

WORKDIR /app
ENTRYPOINT ["./event-indexer"]
CMD ["run"]

COPY --from=builder /build/event-indexer /app/
