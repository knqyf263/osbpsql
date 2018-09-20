package main

import (
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	_ "github.com/lib/pq"
	"github.com/pivotal-cf/brokerapi"

	"github.com/knqyf263/osbpsql/brokers"
	_ "github.com/knqyf263/osbpsql/brokers/psql-db"
	_ "github.com/knqyf263/osbpsql/brokers/psql-user"
	"github.com/knqyf263/osbpsql/config"
)

func main() {
	logger := lager.NewLogger("osbpsql")
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	cfg := config.NewBrokerConfigFromEnv()

	// init api
	logger.Info("Serving", lager.Data{
		"port":     cfg.Port,
		"username": cfg.Credentials.Username,
	})

	serviceBroker, err := brokers.New(cfg, logger)
	if err != nil {
		logger.Fatal("Error initializing service broker: %s", err)
	}

	brokerAPI := brokerapi.New(serviceBroker, logger, cfg.Credentials)
	http.Handle("/", brokerAPI)
	http.ListenAndServe(":"+cfg.Port, nil)
}
