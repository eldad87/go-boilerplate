# Go-Connect
An easy to use, extensible boilerplate for Go applications.

# Motivation
It is very important to write applications that follows the latest standards, Re-using libraries (DRY) and most important; keeping you focused on the actual product.
Which in turn, leads to a more robust and resilient result.

# Concept
Dockeriezed, Production grate, easy to (re)use boilerplate for Go applications. Based on Ashley [McNamara + Brian Ketelsen. Go best practices.](https://www.youtube.com/watch?v=MzTcsI6tn-0 "McNamara + Brian Ketelsen. Go best practices"), [Ben Johnson. Standard Package Layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1 "Ben Johnson. Standard Package Layout"),[ golang-standards]( https://github.com/golang-standards/project-layout " golang-standard") and much more.

# Features
- [Gin](https://github.com/gin-gonic/gin "Gin") - HTTP web framework with smashing performance.
- [Machinery](https://github.com/RichardKnop/machinery "Machinery") -  asynchronous task queue/job queue based on distributed message passing.
- [Logrus](https://github.com/sirupsen/logrus "Logrus") - Structured, pluggable logging for Go.
- [Prometheus](github.com/prometheus/client_golang "Prometheus") -  Instrumentation
- [Health Check](https://github.com/heptiolabs/healthcheck "Health Check") - Implementing Kubernetes liveness and readiness probe handlers
- And much more!

# File Strucutre (src/) (TBD/WIP)
- app/
- cmd/
- config/
- config/development/
- internal/
- pkg/

# TODO
- [x] Docker: MySQL + RabbitMQ setup
- [x] Docker-compose: Jaeger all-in-one
- [x] Gin: Open tracing
- [x] Examples
- [x] Machinery: Open tracing
- [x] Redis Commander
- [x] Logrus Prometheus message type count
- [x] Gin rate limit
- [x] Jaeger/OT + Prometheus
- [ ] Machinery: 
 - [ ] Producer: Hystrix (Conn, CB, TO), throttling/Rate Limit (Juju)
 - [ ] Machinery Redis result backend limits/config (MaxActive, MaxIdle, MaxConnLifetime etc.)
 - [ ] Consumer/AsyncRes: Hystrix (Conn, CB, TO), throttling/Rate Limit (Juju). KIM: Redis already ~protected-ish
 - [ ] Easier produce/consume pattern
 - [ ] Add configuration support for all backends
- [ ] Logrus [Slack report](https://github.com/johntdyer/slackrus)
- [x] Docker: shared /vendor folder for improved debugging expiriance.
- [ ] Consumer throttling
- [ ] Producer throttling
- [ ] Easy tasks registration
- [ ] Protect monitoring HTTP entrypoints (http://localhost/metrics)
- [ ] Unit-test coverage
- [ ] Prometheus server
- [ ] ELK

- [ ] Dockerized SSL support

# Installing
### Docker
Run your project
```
git clone https://github.com/eldad87/go-boilerplate.git
cd co-boilerplate
make init // First time only!
make up
```
##### Dependecnies
In order to manage the project's dependencies, enter project's shell and continue as usual, for example:
```
make shell
dep ensure -add "github.com/username/repo"
exit
```
##### Commands
For all available commands, please checkout the [Makefile](Makefile "Makefile").
### Linux
##### Project
```
git clone https://github.com/eldad87/go-boilerplate.git
cd co-boilerplate
make init // First time only!
dep ensure --vendor-only
go run src/cmd/app/app.go
 ```
##### Instrumentation
###### Jaeger
This step is not required if you already 
Setup a [Jaeger](https://sematext.com/blog/opentracing-jaeger-as-distributed-tracer/) container:
```
sudo docker run -d -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5578 \
 -p 16686:16686 -p 14268:14268 --name jaeger jaegertracing/all-in-one:latest
```
To explore the traces, navigate to http://localhost:16686
Next, check Jaeger (OpenTracing) at http://localhost:16686/ and Redis-Commander at http://localhost:8082/

### Verification
 To verify that you'r app is running correctly, simply browse for the following:
  - http://localhost/health/live -  Kubernetes liveness
  - http://localhost/health/ready  -  Kubernetes readiness
  - http://localhost/metrics - Prometheus instrumentation
  - http://localhost/ping  - echo `{"message":"pong"}
Or, check the logs. The app is writing logs to STDOUT in JSON format.





 ### Configuration
 TBD
  - File structure / env
  - Env > File > Default