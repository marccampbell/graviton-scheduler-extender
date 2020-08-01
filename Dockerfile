FROM golang:1.14 as builder
ARG VERSION=0.0.1

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# build
WORKDIR /go/src/marccampbell/graviton-scheduler-extender
COPY go.mod .
COPY go.sum .
RUN GO111MODULE=on go mod download
COPY . .
RUN make bin

# runtime image
FROM alpine:3.12.0
RUN apk add ca-certificates
COPY --from=builder /go/src/marccampbell/graviton-scheduler-extender/bin//graviton-scheduler-extender /usr/bin/graviton-scheduler-extender
ENTRYPOINT ["/usr/bin/graviton-scheduler-extender"]
