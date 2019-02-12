FROM golang:1.11.1 as builder

# Arguments to Env variables
ARG build_env
ARG app_port
ARG protobuf_release_tag

ENV BUILD_ENV $build_env
ENV APP_PORT $app_port
ENV PROTOBUF_RELEASE_TAG $protobuf_release_tag

# Path
ENV GOBIN=$GOPATH/bin
ENV PATH=$PATH:$GOBIN
WORKDIR $GOPATH/src/github.com/eldad87/go-boilerplate

# Protobuf
RUN apt-get update && \
    apt-get -y install unzip
RUN curl -OL "https://github.com/google/protobuf/releases/download/v${PROTOBUF_RELEASE_TAG}/protoc-${PROTOBUF_RELEASE_TAG}-linux-x86_64.zip" && \
    unzip "protoc-${PROTOBUF_RELEASE_TAG}-linux-x86_64.zip" -d protoc3 && \
    mv protoc3/bin/* /usr/local/bin/ && \
    mv protoc3/include/* /usr/local/include/ && \
    rm -rf $TEMP_DIR/protoc_inst

RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

# Dep
ADD Gopkg.toml Gopkg.toml
ADD Gopkg.lock Gopkg.lock

Run go get -u github.com/golang/dep/cmd/dep
Run dep ensure --vendor-only

# Install the correct protoc-gen-go in the correct version
RUN go install ./vendor/github.com/golang/protobuf/protoc-gen-go/
RUN go install ./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/

# Copy our code
Add src/ src/

# Packr
RUN go get -u github.com/gobuffalo/packr/packr
RUN packr

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /app ./src/cmd/grpc/app.go

# Remove embeded .go files
# RUN packr clean - irrelevant as we start from scratch

# From scratch
FROM scratch

COPY --from=builder /go/src/app/github.com/eldad87/go-boilerplate/protoc3/bin/* /usr/local/bin/
COPY --from=builder /go/src/app/github.com/eldad87/go-boilerplate/protoc3/include/* /usr/local/include/
# COPY --from=builder /go/src/app/github.com/eldad87/go-boilerplate/config/${BUILD_ENV} ./config/src/${BUILD_ENV}

EXPOSE ${APP_GRPC_PORT} ${APP_PORT}
CMD /app