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
- [ ] Docker: shared /vendor folder for improved debugging expiriance.
- [ ] Machinery: 
 - [ ] Add configuration support for all backends
 - [ ] Circuit Breaker on the Broker (Producer)
- [ ] Unit-test coverage
- [ ] Examples
- [ ] Dockerized SSL support
- [ ] Protect monitoring HTTP entrypoints (http://localhost/metrics)

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
```
git clone https://github.com/eldad87/go-boilerplate.git
cd co-boilerplate
make init // First time only!
dep ensure --vendor-only
go run src/cmd/app/app.go
 ```
 ### Verification
 To verify that you'r app is running correctly, simply browse for the following:
  - http://localhost/health/live -  Kubernetes liveness
  - http://localhost/health/ready  -  Kubernetes readiness
  - http://localhost/metrics - Prometheus instrumentation
  - http://localhost/ping  - echo `{"message":"pong"}
Or, check the logs. The app is writing logs to STDOUT in JSON format.