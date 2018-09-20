package psqldb

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"

	"github.com/knqyf263/osbpsql/brokers"
	"github.com/knqyf263/osbpsql/db"
	"github.com/knqyf263/osbpsql/util"
)

//go:generate go-assets-builder -p psqldb -o definition.go definition.yaml

var (
	bindingMap = map[string]Binding{}
)

// Binding is the struct of ServiceBinding
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

	brokers.Register(s, &Broker{})
}

// Broker is ServiceBroker for creating database
type Broker struct {
	db *db.Manager
}

// SetDB sets db
func (b *Broker) SetDB(db *db.Manager) {
	b.db = db
}

// Provision creates database automatically
func (b *Broker) Provision(instanceID string, details brokerapi.ProvisionDetails) error {
	dbName := util.RandLowerString(20)
	bindingMap[instanceID] = Binding{DBName: dbName}
	return b.db.CreateDatabaseG(dbName)
}

// Deprovision drops database
func (b *Broker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails) error {
	binding, ok := bindingMap[instanceID]
	if !ok {
		return fmt.Errorf("Unknown instance id: %s", instanceID)
	}

	err := b.db.DropDatabaseG(binding.DBName)
	if err != nil {
		return errors.Wrap(err, "Unexpected DB error")
	}
	delete(bindingMap, instanceID)

	return nil
}

// Bind returns database name
func (b *Broker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	binding, ok := bindingMap[instanceID]
	if !ok {
		return brokerapi.Binding{}, fmt.Errorf("Unknown instance id: %s", instanceID)
	}
	return brokerapi.Binding{Credentials: binding}, nil
}

// Unbind always returns nil
func (b *Broker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	return nil
}
