module github.com/eldad87/go-boilerplate

go 1.14

require (
	cloud.google.com/go/pubsub v1.2.0
	github.com/BurntSushi/toml v0.3.1
	github.com/RichardKnop/logging v0.0.0-20190827224416-1a693bdd4fae
	github.com/RichardKnop/machinery v1.8.2
	github.com/RichardKnop/redsync v1.2.0
	github.com/TheZeroSlave/zapsentry v1.4.0
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/aws/aws-sdk-go v1.31.3
	github.com/beorn7/perks v1.0.1
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/eapache/go-resiliency v1.1.0
	github.com/envoyproxy/protoc-gen-validate v0.3.0
	github.com/friendsofgo/errors v0.9.2
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-playground/validator/v10 v10.3.0
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-stack/stack v1.8.0
	github.com/gobuffalo/envy v1.9.0 // indirect
	github.com/gobuffalo/packd v1.0.0 // indirect
	github.com/gobuffalo/packr v1.30.1
	github.com/golang/protobuf v1.4.2
	github.com/golang/snappy v0.0.1
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.4.0
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.5
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/ibm-developer/generator-ibm-core-golang-gin v1.0.4
	github.com/jaegertracing/jaeger-client-go v2.23.1+incompatible
	github.com/jaegertracing/jaeger-lib v2.2.0+incompatible
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/jmattheis/go-packr-swagger-ui v3.20.5+incompatible
	github.com/jmespath/go-jmespath v0.3.0
	github.com/kat-co/vala v0.0.0-20170210184112-42e1d8b61f12
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.10.2
	github.com/magefile/mage v1.9.0
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/mitchellh/mapstructure v1.3.0
	github.com/ngrok/sqlmw v0.0.0-20200129213757-d5c93a81bec6
	github.com/opentracing-contrib/go-stdlib v0.0.0-20190519235532-cf7a6c988dc9
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.3.0
	github.com/prometheus/client_model v0.1.0
	github.com/prometheus/common v0.7.0
	github.com/prometheus/procfs v0.0.8
	github.com/rogpeppe/go-internal v1.6.0
	github.com/rubenv/sql-migrate v0.0.0-20200429072036-ae26b214fa43
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cast v1.3.1
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.0
	github.com/srikrsna/protoc-gen-gotag v0.5.0
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	github.com/uber/jaeger-client-go v2.23.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible
	github.com/volatiletech/null/v8 v8.1.0
	github.com/volatiletech/randomize v0.0.1
	github.com/volatiletech/sqlboiler/v4 v4.1.2
	github.com/volatiletech/strmangle v0.0.1
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c
	github.com/xdg/stringprep v1.0.0
	go.mongodb.org/mongo-driver v1.3.3
	go.opencensus.io v0.22.3
	go.uber.org/atomic v1.5.0
	go.uber.org/zap v1.13.0
	golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073
	golang.org/x/mod v0.2.0
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	golang.org/x/sys v0.0.0-20200519105757-fe76b779f299
	golang.org/x/text v0.3.2
	golang.org/x/tools v0.0.0-20200521155704-91d71f6c2f04
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
	google.golang.org/api v0.25.0
	google.golang.org/appengine v1.6.5
	google.golang.org/genproto v0.0.0-20200521103424-e9a78aa275b7
	google.golang.org/grpc v1.28.0
	google.golang.org/protobuf v1.23.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	gopkg.in/guregu/null.v4 v4.0.0
	gopkg.in/ini.v1 v1.56.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	honnef.co/go/tools v0.0.1-2020.1.4
)
