# Go-Boilerplate
An easy to use, extensible boilerplate for Go applications.

# Motivation
It is very important to write applications that follows the latest standards, Re-using libraries (DRY) and most important; keeping you focused on the actual logic.
Which in turn, leads to a robust and resilient product.

# Concept
Dockeriezed, Production grade, easy to (re)use boilerplate for Go applications. Based on Ashley [McNamara + Brian Ketelsen. Go best practices](https://www.youtube.com/watch?v=MzTcsI6tn-0 "McNamara + Brian Ketelsen. Go best practices"), [Ben Johnson. Standard Package Layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1 "Ben Johnson. Standard Package Layout"),[ golang-standards]( https://github.com/golang-standards/project-layout " golang-standard") and much more.

# Features
- [gRPC](https://grpc.io/ "gRPC") - A high-performance, open-source universal RPC framework.
- [gRPC request validator](https://github.com/mwitkow/go-proto-validators "gRPC request validator") - A protoc plugin that generates Validate() error functions on Go proto structs based on field options inside.
- [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway/ "gRPC-Gateway") - gRPC to JSON proxy generator following the gRPC HTTP spec.
- [OpenAPI](https://github.com/grpc-ecosystem/grpc-gateway/ "gRPC-Gateway") - Online Documentation for our gRPC-Gateway APIs.
- [Machinery](https://github.com/RichardKnop/machinery "Machinery") -  An asynchronous task queue/job queue based on distributed message passing.
- [Zap](https://github.com/uber-go/zap "Zap") - Blazing fast, structured, leveled logging in Go
- [Grift](https://github.com/markbates/grift "Grift") - Go based task runner
- [Prometheus](https://github.com/prometheus/client_golang "Prometheus") -  Instrumentation
- [Health Check](https://github.com/heptiolabs/healthcheck "Health Check") - Implementing Kubernetes liveness and readiness probe handlers
- [SQL-Migrate](https://github.com/rubenv/sql-migrate "SQL-Migrate") - SQL schema migration tool for Go
- [SQLBoiler](https://github.com/volatiletech/sqlboiler "SQLBoiler") - Generate a Go ORM tailored to your database schema.
- And much more!

# File Strucutre (src/) (TBD/WIP)
- app/
- cmd/
- config/
- grifts/
- config/development/
- internal/
- migration/
- pkg/

# Examples (TBD/WIP) 
### Migration / Seed
- Add your migration to the `src/migration` folder (1549889122_init_schema.sql), Make sure to follow the [standards](https://github.com/rubenv/sql-migrate#writing-migrations "standards"): 
```sql
-- +migrate Up
CREATE TABLE visits (
    id int NOT NULL AUTO_INCREMENT,
    first_name varchar(255),
    last_name varchar(255),
    created_at timestamp default NOW(),
    updated_at timestamp default NOW(),
    PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS visits;

```
- Run ```make grift db:migrate```
- You're done!
- For additional information, make sure to visit the official [repository](https://github.com/rubenv/sql-migrate "repository"): 

### ORM / SQLBoiler Usage
SQLBoiler is a tool to generate a Go ORM tailored to your database schema. Which result in a fully featured, ActiveRecord like ORM without any [performance](https://github.com/volatiletech/sqlboiler#benchmarks "performance") hit.
- Make sure `sqlboiler.yaml` is configured correctly and points to the relevant databse: `cat src/config/development/sqlboiler.yaml`
- Run `make sqlboiler`
- The auto generated ORM is now located in `src/app/mysql*or-any-other-engine*/models`
- Example (Based on the migration above):
```go
import (
    "fmt"
    "context"
    "database/sql"
    "github.com/go-sql-driver/mysql"
)

func main() {
    visitId = 1
    
    db, err := sql.Open("mysql", "database.dsn")
    if err != nil {
        panic("Failed to open DB connection")
    }
    
    if err := db.Ping(); err != nil {
        panic("Failed to ping DB")
    }
    
    if visit, err := models.FindVisit(context.Background(), db, visitId); err != nil {
        fmt.Print("")
    } else {
        fmt.Print(visit.FirstName, visit.LastName)
    }
}

```
- For additional information, make sure to visit the official [repository](https://github.com/volatiletech/sqlboiler "repository"): 

### Add CLI task
Grift is a very simple library that allows you to write simple "task" scripts in Go and run them by name without having to write big main type of wrappers. Grift is similar to, and inspired by, Rake.
- Add your TASK to the `src/grifts` folder (migrate.go), Make sure to follow the [standards](https://godoc.org/github.com/markbates/grift/grift "standards"): 
```go
package grifts

import (
	"database/sql"
	"fmt"
	"github.com/gobuffalo/packr"
	. "github.com/markbates/grift/grift"
	"github.com/rubenv/sql-migrate"
)

var _ = Namespace("db", func() {
    Desc("migrate", "Migrates the databases")
    Set("migrate", func(c *Context) error {
        db, err := sql.Open("mysql", "database.dsn")
        if err != nil {
            panic("Failed to open DB connection")
        }
        
        if err := db.Ping(); err != nil {
            panic("Failed to ping DB")
        }
        
        migrations := &migrate.PackrMigrationSource{
            Box: packr.NewBox("../../src/migration"),
        }
        appliedMigrations, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
        if err != nil {
            fmt.Print(err)
        }
        fmt.Printf("Applied %v migrations", appliedMigrations)
        
        return nil
    })
})
```
- Now you can run your task ```make grift db:migrate```
- You're done!
- For additional information, make sure to visit the official [repository](https://godoc.org/github.com/markbates/grift/grift  "repository"): 
### Machinery
- Consumer
```go
TODO: Register "repeat(str string) string { return str }" as "repeat"
```
- Producer, Using machinery
```go
    import (
    	"github.com/afex/hystrix-go/hystrix"
    	reHystrix "github.com/eldad87/go-boilerplate/src/pkg/concurrency/hystrix"
    	machineryProducer "github.com/eldad87/go-boilerplate/src/pkg/task/producer/machinery"
        "github.com/eldad87/go-boilerplate/src/pkg/task/producer"
        "log"
        "tiume"
    )
    
    macineryServer := .. Your machinery instance ..
    
    // Initiate our producer, configure its limitations using Hystrix abd retrier
	producer := machineryProducer.NewProducer(server,
	    // Uses Hystrix as circuit breaker
		hystrix.CommandConfig{
			Timeout:                conf.GetInt("machinery.broker.timeout"),
			MaxConcurrentRequests:  conf.GetInt("machinery.broker.max_conn"),
			RequestVolumeThreshold: conf.GetInt("machinery.broker.vol_threshold"),
			SleepWindow:            conf.GetInt("machinery.broker.sleep_window"),
			ErrorPercentThreshold:  conf.GetInt("machinery.broker.err_per_threshold"),
		},
		// Retry X attepts with Y delay between each failure
		reHystrix.RetryConfig{Attempts: conf.GetInt("machinery.broker.retries"), Delay: conf.GetDuration("machinery.broker.retry_delay")})

    // Lets send a task throught the wire
    // First thing, define our options
    //  You can pass things like eta, routingKey, onSuccess and more.
    //  Or any other attribute in request.Request
    options := map[string]interface{}{} 
    // Next, lets actually send the task, note that "repeat" is already registered on the consumer side
    if request, err := producer.NewRequest("repeat", options, "Hello"); err == nil {
        asyncResult, err := producer.Produce(request, nil)
        err := asyncResult.Subscribe(time.Duration(1000) // Wait 1sec
        
        if err != nil {
            // You can check if we exceded our subscription duration (1sec)
            // if e, ok := err.(producer.ErrTimeoutReached); ok == true {
            //     log.error("Timeout reached")
            // }
            log.error(e)
        } else {
            // We got a response, it doesn't mean that everything went smooth.
            //  we still need to check if the response was succesfull:
            // asyncResult.IsSuccess()
            // asyncResult.IsFailure()
            // results := asyncResult.Result() // []producer.Result / results[0].Type == "string" / results[0].Value == "Hello"
        }
        
        return asyncResult, nil
    } else {
        log.error("Unable to call \"repeat\")
   }
```

# TODO
- [x] Docker: MySQL + RabbitMQ setup
- [x] Docker-compose: Jaeger all-in-one
- [x] Gin: Open tracing
- [x] Examples
- [x] Machinery: Open tracing
- [x] Redis Commander
- [x] Zap Prometheus message type count
- [x] Jaeger/OT + Prometheus
- [x] gRPC, HTTP Gateway and Swagger-UI
- [x] gRPC and HTTP Gateway validation error sync
- [x] gRPC opentracing and instrumentation
- [x] Task utility similar to Rake
- [x] Auto migration run during development / Manual task
- [X] Embed OpenAPI using packr
- [ ] Run health checks and metrics on a different port then gRPC-gateway
- [ ] Machinery: 
  - [x] Producer and Result interface/wrapper
  - [x] Producer: Hystrix (Conn, CB, TO)
  - [ ] Consumer: task registration pattern
  - [ ] Machinery Redis result backend limits/config (MaxActive, MaxIdle, MaxConnLifetime etc.)
  - [ ] Add configuration support for all backends (including healthchecks)
- [x] Docker: shared /vendor folder for improved debugging expiriance.
- [x] Healtcheck for Redis, AMQP and Goroutine Threshold
- [ ] Protect monitoring HTTP entrypoints (http://localhost/metrics)
- [ ] Unit-test coverage
- [ ] Prometheus server
- [ ] Log shipping
- [ ] Dockerized SSL support
- [ ] Hystrix turbine/dashboard 

# Installing
### Docker
Run your project
```bash
git clone https://github.com/eldad87/go-boilerplate.git
cd co-boilerplate
make init // First time only!
make up
```
##### Dependecnies
In order to manage the project's dependencies, enter project's shell and continue as usual, for example:
```bash
make shell
dep ensure -add "github.com/username/repo"
exit
```
##### Commands
For all available commands, please checkout the [Makefile](Makefile "Makefile").
### Linux
##### Project
```bash
git clone https://github.com/eldad87/go-boilerplate.git
cd co-boilerplate
make init // First time only!
dep ensure --vendor-only
go run src/cmd/app/app.go
 ```
##### Instrumentation
###### Jaeger
You can skip this step if you already gave a running instance of [Jaeger](https://sematext.com/blog/opentracing-jaeger-as-distributed-tracer/):
```bash
sudo docker run -d -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5578 \
 -p 16686:16686 -p 14268:14268 --name jaeger jaegertracing/all-in-one:latest
```
To explore the traces, navigate to http://localhost:16686
Next, check Jaeger (OpenTracing) at http://localhost:16686/ and Redis-Commander at http://localhost:8083/

### Verification
 To verify that your project is running correctly, simply browse the following:
  - http://localhost/health/live - Kubernetes liveness
  - http://localhost/health/ready - Kubernetes readiness
  - http://localhost/metrics - Prometheus instrumentation
  - http://localhost/ping - echo `{"message":"pong"}
  - http://localhost:8080/swaggerui/ - Swagger UI
  - http://localhost:8080/v1/visit/__INT__ - gRPC Gateway, replace __INT__ with any numeric value
Or, check the logs. Logs are writing STDOUT in a JSON format.
http://localhost:16686/search

 ### Configuration
 TBD
  - File structure / env
  - Env > File > Default