module github.com/eldad87/go-boilerplate

go 1.14

require (
	cloud.google.com/go/pubsub v1.3.1
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
	github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/friendsofgo/errors v0.9.2
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-playground/validator/v10 v10.3.0
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-stack/stack v1.8.0
	github.com/gobuffalo/envy v1.9.0 // indirect
	github.com/gobuffalo/packd v1.0.0 // indirect
	github.com/gobuffalo/packr v1.30.1
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.1
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.5.4
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.9.5
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.2.0
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/ibm-developer/generator-ibm-core-golang-gin v1.0.4
	github.com/jaegertracing/jaeger-client-go v2.23.1+incompatible
	github.com/jaegertracing/jaeger-lib v2.2.0+incompatible
	github.com/jmattheis/go-packr-swagger-ui v3.20.5+incompatible
	github.com/jmespath/go-jmespath v0.3.0
	github.com/kat-co/vala v0.0.0-20170210184112-42e1d8b61f12
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.10.2
	github.com/magefile/mage v1.9.0
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/mitchellh/mapstructure v1.4.1
	github.com/ngrok/sqlmw v0.0.0-20200129213757-d5c93a81bec6
	github.com/opentracing-contrib/go-stdlib v1.0.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.3.0
	github.com/prometheus/client_model v0.1.0
	github.com/prometheus/common v0.7.0
	github.com/prometheus/procfs v0.0.8
	github.com/rogpeppe/go-internal v1.6.0
	github.com/rubenv/sql-migrate v0.0.0-20200429072036-ae26b214fa43
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cast v1.3.1
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	github.com/uber/jaeger-client-go v2.23.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible
	github.com/volatiletech/null/v8 v8.1.0
	github.com/volatiletech/randomize v0.0.1
	github.com/volatiletech/sqlboiler/v4 v4.4.0
	github.com/volatiletech/strmangle v0.0.1
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c
	github.com/xdg/stringprep v1.0.0
	go.mongodb.org/mongo-driver v1.3.3
	go.opencensus.io v0.22.4
	go.uber.org/atomic v1.5.0
	go.uber.org/zap v1.13.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/mod v0.3.0
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	golang.org/x/oauth2 v0.0.0-20210201163806-010130855d6c
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c
	golang.org/x/text v0.3.5
	golang.org/x/tools v0.0.0-20200828161849-5deb26317202
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	google.golang.org/api v0.30.0
	google.golang.org/appengine v1.6.6
	google.golang.org/genproto v0.0.0-20210207032614-bba0dbe2a9ea
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	gopkg.in/guregu/null.v4 v4.0.0
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	honnef.co/go/tools v0.0.1-2020.1.4
)
