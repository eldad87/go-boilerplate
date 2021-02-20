FROM golang:1.14.2 as builder

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

# Protobuf, gRPC and Gateway
RUN set -e && \
    GO111MODULE=on go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v${grpc_gateway_version} && \
    cd /go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v${grpc_gateway_version}/protoc-gen-grpc-gateway && \
    go install .

RUN set -e && \
    GO111MODULE=on go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v${grpc_gateway_version} && \
    cd /go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v${grpc_gateway_version}/protoc-gen-openapiv2 && \
    go install .

RUN go get -u google.golang.org/protobuf/cmd/protoc-gen-go
RUN go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

RUN go get -u github.com/srikrsna/protoc-gen-gotag
RUN go get -d github.com/envoyproxy/protoc-gen-validate && make build -C /go/src/github.com/envoyproxy/protoc-gen-validate/

# SQLBoiler
RUN GO111MODULE=off go get -u -t github.com/volatiletech/sqlboiler
# Also install the driver of your choice, there exists pqsl, mysql, mssql
RUN GO111MODULE=off go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql

# Go based task runner
RUN git clone https://github.com/magefile/mage
RUN cd mage && go run bootstrap.go
RUN cd $GOPATH/src/github.com/eldad87/go-boilerplate && rm -rf mage

# Dep
ADD go.mod go.mod
ADD go.sum go.sum
Run go mod download

# Install the correct protoc-gen-go in the correct version
# RUN go install ./vendor/github.com/golang/protobuf/protoc-gen-go/

# Copy our code
Add src/ src/

# Packr
RUN go get -u github.com/gobuffalo/packr/packr
RUN packr

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /app ./src/cmd/grpc/app.go

# Remove embeded .go files, used in case and tested locally
RUN packr clean

# From scratch
FROM scratch

COPY --from=builder /app /app
COPY --from=builder /usr/local/bin/* /usr/local/bin/
COPY --from=builder /usr/local/include/* /usr/local/include/
COPY --from=builder /go/bin/mage /go/bin/mage
# COPY --from=builder /go/src/app/github.com/eldad87/go-boilerplate/config/${BUILD_ENV} ./config/src/${BUILD_ENV}

EXPOSE ${APP_GRPC_PORT} ${APP_PORT}
CMD /app