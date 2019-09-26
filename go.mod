module github.com/vulcanize/mcd_transformers

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/apilayer/freegeoip v3.5.0+incompatible // indirect
	github.com/btcsuite/btcd v0.0.0-20190115013929-ed77733ec07d // indirect
	github.com/cespare/cp v1.1.1 // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/elastic/gosigar v0.10.5 // indirect
	github.com/ethereum/go-ethereum v1.9.5
	github.com/google/uuid v1.1.0 // indirect
	github.com/influxdata/influxdb v1.7.7 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/oschwald/maxminddb-golang v1.5.0 // indirect
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/afero v1.2.1 // indirect
	github.com/spf13/viper v1.3.2
	github.com/vulcanize/vulcanizedb v0.0.7
	golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5
	google.golang.org/appengine v1.6.2 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
)

replace github.com/ethereum/go-ethereum => github.com/vulcanize/go-ethereum v1.5.10-0.20190910005838-ca79f6ef9877

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/vulcanize/vulcanizedb => github.com/vulcanize/vulcanizedb v0.0.8-0.20190925215242-5c0e5592abd2
