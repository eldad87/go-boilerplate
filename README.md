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
- [Viper](https://github.com/spf13/viper "Viper") - Go configuration with fangs.
- And much more!

# Boilerplate
### File Strucutre (src/)
- `app/` - Where we define our service(s) and common data structure(s) for better inter-service communication and decoupling. For example:
  - `app/visit.go` - contain our `VisitService interface` and `Visit Struct`
  - `app/mysql/visit.go` - Implements our `VisitService` using `MySQL` as our DB.  
  In case and you want to use a different DB, Simply add a new folder (e.g. `mongodb`) and implement the `VisitService` accordingly.
  - From now on, any service that will relay on the `VisitService` will be decoupled from it's implementation. 
- `cmd/` - Our App can compile into different executables (multiple `main()` functions), each run a different variation of our App.
  - For example a Worker app (`cmd/worker/app.go`) can be built using the same (bulletproof & heavily tested) source code as our main cmd (`cmd/grpc/app.go`).  
  In addition, excluding a lot of unnecessary dependency will lead to a better memory consumption.
- `config/` - Configuration management is implemented using [Viper](https://github.com/spf13/viper "Viper"), It loads configuration based on the following priority:
  - Environment variables
  - `config/ENV-NAME/ENgit puV-NAME.yaml` - Yaml configuration file, where `ENV-NAME` is `development` etc.  
  In order to use a configuration file with a compiled App - Simply preserve the same file structure. Anyway, you better using Env variables (or defaults) instad
  - Default values - hard coded into `config.go`
- `grifts/` - CLI task, more information is listed below.
- `migration/` - Migrations, more information is listed below.
- `pkg/` - public packages that can be used on different projects. Each package is a stand-alone unit and can function out-side of your source code.
- `transport/` - Where we define our we transport our data (HTTP handlers, gRPC etc) 

#### Make file
The `make` file is mainly used as a "shortcut" to highly used tools such as docker, auto genration etc.  

# Document by Example (TBD/WIP) 
In order to build a rock-solid application you usually use the same set of tools over and over again, either in a form of a different language or a framework. It's all the same:
- Configuration management
- DB Migrations
- ORM
- CLI tasks
- Handlers/Controllers
- Async Job processing
- Health Check/Instrumentation
- Logs
- Documentation
- and finally Tests

The following examples will demonstrate how each functionality can be used.

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

### Service / Data layer
This is probably the most important part in the boilerplate, the goal is to create a separation between the data layer and the way other component interacts with it.

- To begin with, we need to define our basic data struct, pay attention to the tags as they'll become handy later on:
```go
package app

type Visit struct {
	ID        uint      `validate:"gte=0"`
	FirstName string    `validate:"required,gte=3,lte=254"`
	LastName  string    `validate:"required,gte=3,lte=254"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
```
- Next, define our service `VisitService`:
```go
// ..

type VisitService interface {
	Get(c context.Context, id *uint) (*Visit, error)
}
```
- Save both `Visit struct` and `VisitService interface` as `app/visit.go`
- From now on, to avoid coupling, any Component that relay on the `VisitService` will use it as a dependency so it won't be worried on how its been implemented, for example:
```go
type MyController struct {
	VisitService app.VisitService
}
```
- Assuming you followed the examples above and that you generated your SQLBoiler code, the database is migrated etc you can start implementing your service.
- Create a new folder under `app/.` that represent your data layer (e.g `app/mysql`)
- Implement your service `app/mysql/visit.go`:  
```go
package mysql

import (
	"context"
	"database/sql"
	"github.com/eldad87/go-boilerplate/src/app"
	"github.com/eldad87/go-boilerplate/src/app/mysql/models"
	"github.com/eldad87/go-boilerplate/src/pkg/validator"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func NewVisitService(db *sql.DB, sv validator.StructValidator) *VisitService {
    return &VisitService{db, sv}
}

type VisitService struct {
    db *sql.DB
    sv validator.StructValidator
}

func (vs *VisitService) Get(c context.Context, id *uint) (*app.Visit, error) {
    bVisit, err := models.FindVisit(c, vs.db, *id)
    if err != nil {
        return nil, err
    }
    
    return sqlBoilerToVisit(bVisit), nil
}

func (vs *VisitService) Set(c context.Context, v *app.Visit) (*app.Visit, error) {
    bVisit := models.Visit{
        ID:        v.ID,
        FirstName: null.StringFrom(v.FirstName),
        LastName:  null.StringFrom(v.LastName),
    }
    
    // Validate our Struct using the "validate" tags
    err := vs.sv.StructCtx(c, v)
    if err != nil {
        return nil, err
    }
    
    if bVisit.ID == 0 {
        err = bVisit.Insert(c, vs.db, boil.Infer())
    } else {
        _, err = bVisit.Update(c, vs.db, boil.Infer())
    }
    
    if err != nil {
        return nil, err
    }
    
    return sqlBoilerToVisit(&bVisit), nil
}

func sqlBoilerToVisit(bVisit *models.Visit) *app.Visit {
    return &app.Visit{
        ID:        bVisit.ID,
        FirstName: bVisit.FirstName.String,
        LastName:  bVisit.LastName.String,
        CreatedAt: bVisit.CreatedAt,
        UpdatedAt: bVisit.UpdatedAt,
    }
}
```


### Handlers/Controllers - gRPC + gRPC-Gateway
In this example, We will
 - Define a simple gRPC handler
 - Expose it as a RESTful API
 - Set the request constrains (Validations)
 - Document it using OpenAPI (Swagger)

So lets start:
- Define our `VisitTransport` using protobuf in `transport/grpc/proto/visit_transport.proto`, pay attention to the constrain on `ID`: 
```proto
// ...

service VisitTransport {
    // Simple return a Visit record by ID
    rpc Get(ID) returns (VisitResponse) {
        option (google.api.http) = {
          get: "/v1/visit/{id}" // RESTful route
        };
    }
    // Update/Create a device
    rpc Set(VisitRequest) returns (VisitResponse) {
        option (google.api.http) = {
          post: "/v1/visit"
          body: "*"
          additional_bindings {
            put: "/v1/visit"
            body: "*"
          }
        };
    }
}

// Define a request and constrains, ID must be > 0
message ID {
    uint32 id = 1 [(validate.rules).uint32.gte = 0];
};

message VisitResponse {
    uint32 id = 1;
    string first_name = 2;
    string last_name = 3;
    google.protobuf.Timestamp created_at = 4;
    google.protobuf.Timestamp updated_at = 5;
};

message VisitResponse {
    uint32 id = 1;
    string first_name = 2;
    string last_name = 3;
    google.protobuf.Timestamp created_at = 4;
    google.protobuf.Timestamp updated_at = 5;
};
```
- Generate our gRPC handler, validators, RESTful API and documentation (Swagger):  
  - `make protobuf`  
  \* The auto-generated code is located under the same folder as our `visit_transport.proto` file.
- Implement the auto-generated `VisitTransportServer interface` (can be found in in `transport/grpc/proto/visit_transport.pb.go`).
  - Important! - Use your services to persist your data (`app/*.go`), Inject them as dependency when needed.
  - Store your implementation in a different folder (e.g `transport/grpc/visit_transport.go`) for better separation.
```go
package grpc

import (
	"context"
	"github.com/eldad87/go-boilerplate/src/app"
	pb "github.com/eldad87/go-boilerplate/src/transport/grpc/proto"
	"github.com/golang/protobuf/ptypes"
)

type VisitTransport struct {
	VisitService app.VisitService
}

func (vs *VisitTransport) Get(c context.Context, id *pb.ID) (*pb.VisitResponse, error) {
	i := uint(id.GetId())
	v, err := vs.VisitService.Get(c, &i)
	if err != nil {
		return nil, err
	}

	return vs.visitToProto(v)
}

// Update/Create a device
func (vs *VisitTransport) Set(c context.Context, v *pb.VisitRequest) (*pb.VisitResponse, error) {
	aVis, err := vs.protoToVisit(v)
	if err != nil {
		return nil, err
	}

	gVis, err := vs.VisitService.Set(c, aVis)
	if err != nil {
		return nil, err
	}

	return vs.visitToProto(gVis)
}

func (vs *VisitTransport) visitToProto(visit *app.Visit) (*pb.VisitResponse, error) {
	created, err := ptypes.TimestampProto(visit.CreatedAt)
	if err != nil {
		return nil, err
	}

	updated, err := ptypes.TimestampProto(visit.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &pb.VisitResponse{Id: uint32(visit.ID), FirstName: visit.FirstName, LastName: visit.LastName, CreatedAt: created, UpdatedAt: updated}, nil
}

func (vs *VisitTransport) protoToVisit(visit *pb.VisitRequest) (*app.Visit, error) {
	return &app.Visit{
		ID:        uint(visit.Id),
		FirstName: visit.FirstName,
		LastName:  visit.LastName,
	}, nil
}

```
- Register your gRPC handler
```go
import(
    // ..
    
    "database/sql"
    v9validator "gopkg.in/go-playground/validator.v9"
    service "github.com/eldad87/go-boilerplate/src/app/mysql"
    grpcService "github.com/eldad87/go-boilerplate/src/transport/grpc"
    pb "github.com/eldad87/go-boilerplate/src/transport/grpc/proto"
)

func main() {
    // ..
	
    grpcServer := grpc.NewServer(
        // ..
    )
    defer grpcServer.GracefulStop()
    
    db, err := sql.Open("mysql", "DSN...")
    if err != nil {
    	log.Fatal("Cannot connect to DB")
    }
	
    mySQLVisitService := service.NewVisitService(db, v9validator.New())
    visitTrans := grpcService.VisitTransport{VisitService: mySQLVisitService}
    pb.RegisterVisitServiceServer(grpcServer, &visitTrans)
    
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatal("Cannot start gRPC server")
    }
}
```
- Now you can query your service using cURL (e.g `localhost/visit/1`) or a gRPC client 
- You're done!

### Logger
Logger is a simple tool, yet probably the most useful one. Developers tend to consider logs as a "Free" resource, but in reality it can become quite pricey. Especially when your log-shipper need to do some parsing and aggregations.   
In this implementation I used [Zap](https://github.com/uber-go/zap "Zap") as it provide everything that an equivalent library does (e.g Logrush) and yet 10x faster.
```go
import(
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

zapConfig := zap.NewProductionConfig()
zapConfig.Level.UnmarshalText([]byte("debug"))
zapConfig.Development = true
// Log our time, set the time as a string field called "time"
logger.Info("Hi there, the time is", zap.String("time", time.Now().Format(time.RFC850)))

```

### Async Job processing
##### Machinery
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