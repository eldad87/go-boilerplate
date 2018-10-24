FROM golang:1.11.1 as builder

# Arguments to Env variables
ARG build_env
ARG app_port

ENV BUILD_ENV $build_env
ENV APP_PORT $app_port

# Path
ENV GOBIN=$GOPATH/bin
ENV PATH=$PATH:$GOBIN
WORKDIR $GOPATH/src/github.com/eldad87/go-boilerplate

# Dep
ADD Gopkg.toml Gopkg.toml
ADD Gopkg.lock Gopkg.lock

Run go get -u github.com/golang/dep/cmd/dep
Run dep ensure --vendor-only

# Copy our code
Add src/ src/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /app ./src/cmd/app/app.go

# From scratch
FROM scratch

COPY --from=builder /app /app
# COPY --from=builder /go/src/app/src/config/${BUILD_ENV} ./config/src/${BUILD_ENV}

EXPOSE ${APP_PORT}
CMD /app