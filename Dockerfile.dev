FROM golang:1.11.1 as builder

# Arguments to Env variables
ARG build_env
ARG app_port
ARG app_grpc_port
ARG protobuf_release_tag

ENV BUILD_ENV $build_env
ENV APP_PORT $app_port
ENV APP_GRPC_PORT $app_grpc_port
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
    rm -rf protoc3 && \
    rm protoc-${PROTOBUF_RELEASE_TAG}-linux-x86_64.zip

RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN go get -u github.com/lyft/protoc-gen-validate

# SQLBoiler
RUN go get -u -t github.com/volatiletech/sqlboiler
# Also install the driver of your choice, there exists pqsl, mysql, mssql
RUN go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql

# Go based task runner
RUN go get -u github.com/markbates/grift

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

# Build or install hot-reload
RUN go get -u github.com/VojtechVitek/rerun/cmd/rerun

# Run binary or hot-reload
CMD rerun -watch ./ -ignore vendor bin migration -run go run ./src/cmd/grpc/app.go

EXPOSE ${APP_GRPC_PORT} ${APP_PORT}