FROM golang:1.15-alpine as builder
ARG VERSION=0.0.1

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# config
WORKDIR /go/src/k8s-admission-webhook
COPY go.mod .
COPY go.sum .
RUN GO111MODULE=on go mod download
COPY . .
RUN go install -ldflags "-s -w -X main.version=$VERSION" ./

# runtime image
FROM gcr.io/google_containers/ubuntu-slim:0.14
COPY --from=builder /go/bin/* /usr/bin/k8s-admission-webhook
ENTRYPOINT ["k8s-admission-webhook"]