package psqldb

import (
	"context"

	"github.com/knqyf263/osbpsql/brokers"
	"github.com/knqyf263/osbpsql/db"
	"github.com/pivotal-cf/brokerapi"
	yaml "gopkg.in/yaml.v2"
)

//go:generate go-assets-builder -p psqldb -o definition.go definition.yaml

func init() {
	f, err := Assets.Open("/definition.yaml")
	if err != nil {
		panic(err)
	}

	var s brokerapi.Service
	if err = yaml.NewDecoder(f).Decode(&s); err != nil {
		panic(err)
	}

	brokers.Register(s, &PsqlDBBroker{})
}

type PsqlDBBroker struct {
	db *db.Manager
}

func (b *PsqlDBBroker) SetDB(db *db.Manager) {
	b.db = db
}

func (b *PsqlDBBroker) Provision(instanceID string, details brokerapi.ProvisionDetails) error {
	return b.db.CreateDatabase(instanceID)
}

func (b *PsqlDBBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails) error {
	return b.db.DropDatabase(instanceID)
}
