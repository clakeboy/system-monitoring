module system-monitoring

go 1.16

require (
	github.com/DataDog/zstd v1.4.5 // indirect
	github.com/Sereal/Sereal v0.0.0-20200820125258-a016b7cda3f3 // indirect
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/asdine/storm v1.1.0
	github.com/clakeboy/golib v1.3.9
	github.com/elastic/go-elasticsearch/v7 v7.12.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/shirou/gopsutil v3.21.4+incompatible
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/vmihailenco/msgpack.v2 v2.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/clakeboy/golib => ../golib
