package registry

import "github.com/PMoneda/telemetry/influxdb"

type Registry interface {
	Registry(tag string, value interface{}) error
	Flush(tag string) error
	FlushAll() error
}

//Config is a general config to any datasource
type Config struct {
	Host string
	DB   string
	Port string
}

//NewInfluxClient creates a new endpoint to influx
func NewInfluxClient(config Config) Registry {
	return influxdb.New(config.Host, config.Port)
}
