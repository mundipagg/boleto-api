package healthcheck

import (
	"os"
	"time"

	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
)

func EnsureDependencies() {
	ensureMongo()
}

func ensureMongo() {
	l := log.CreateLog()
	_, err := db.CreateMongo()
	if err != nil {
		l.Error(err.Error(), "healthcheck.ensureMongo - Error creating mongo connection")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
}
