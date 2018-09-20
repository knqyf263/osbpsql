package psqldb

import (
	"context"
	"fmt"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"

	"github.com/knqyf263/osbpsql/brokers"
	"github.com/knqyf263/osbpsql/db"
	"github.com/knqyf263/osbpsql/util"
	"github.com/pivotal-cf/brokerapi"
	yaml "gopkg.in/yaml.v2"
)

//go:generate go-assets-builder -p psqldb -o definition.go definition.yaml

var (
	bindingMap = map[string]Binding{}
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
	dbName := util.RandString(20)
	bindingMap[instanceID] = Binding{DBName: dbName}
	return b.db.CreateDatabase(dbName)
}

func (b *PsqlDBBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails) error {
	binding, ok := bindingMap[instanceID]
	if !ok {
		return fmt.Errorf("Unknown instance id: %s", instanceID)
	}
	err := b.db.DropDatabase(binding.DBName)
	if err != nil {
		return err
	}
	delete(bindingMap, instanceID)
	return nil
}

// Bind returns database name and save the state
func (b *PsqlDBBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	_, err := b.db.FindServiceBindingByID(instanceID)
	if err == nil {
		return brokerapi.Binding{}, brokerapi.ErrBindingAlreadyExists
	} else if err != pg.ErrNoRows {
		return brokerapi.Binding{}, err
	}

	binding, ok := bindingMap[instanceID]
	if !ok {
		return brokerapi.Binding{}, fmt.Errorf("Unknown instance id: %s", instanceID)
	}

	if err = b.db.CreateServiceBindingDetails(instanceID, brokerapi.Succeeded, details); err != nil {
		return brokerapi.Binding{}, errors.Wrap(err, "Failed to save ServiceBinding details to DB")
	}

	return brokerapi.Binding{Credentials: binding}, nil
}

// Unbind deletes the ServiceBinding state from DB
func (b *PsqlDBBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	_, err := b.db.FindServiceBindingByID(instanceID)
	if err == pg.ErrNoRows {
		return brokerapi.ErrBindingDoesNotExist
	} else if err != nil {
		return err
	}
	return b.db.DeleteServiceBindingDetails(instanceID)
}
