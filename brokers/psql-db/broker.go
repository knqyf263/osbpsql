package psqldb

import (
	"context"
	"fmt"

	"github.com/knqyf263/osbpsql/brokers"
	"github.com/knqyf263/osbpsql/db"
	"github.com/pivotal-cf/brokerapi"
	yaml "gopkg.in/yaml.v2"
)

//go:generate go-assets-builder -p psqldb -o definition.go definition.yaml

var (
	bindingMap map[string]Binding
)

type Binding struct {
	DBName string
}

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
	bindingMap[instanceID] = Binding{DBName: instanceID}
	return b.db.CreateDatabase(instanceID)
}

func (b *PsqlDBBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails) error {
	delete(bindingMap, instanceID)
	return b.db.DropDatabase(instanceID)
}

func (b *PsqlDBBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	binding, ok := bindingMap[instanceID]
	if !ok {
		return brokerapi.Binding{}, fmt.Errorf("Unknown instance id: %s", instanceID)
	}
	return brokerapi.Binding{Credentials: binding}, nil
}

func (b *PsqlDBBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	return nil
}
