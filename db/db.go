package db

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/knqyf263/osbpsql/config"
	"github.com/pivotal-cf/brokerapi"
)

type ServiceInstance struct {
	ID      string
	State   brokerapi.LastOperationState
	Details brokerapi.ProvisionDetails
}

type ServiceBinding struct {
	Id    string
	State string
}

type Manager struct {
	db *pg.DB
}

func New(conf config.DBConfig) *Manager {
	db := pg.Connect(&pg.Options{
		Addr: fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		User: conf.User,
	})
	return &Manager{db: db}
}

func (m Manager) Migrate() error {
	for _, model := range []interface{}{(*ServiceInstance)(nil), (*ServiceBinding)(nil)} {
		err := m.db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m Manager) FindServiceInstanceById(instanceID string) error {
	return m.db.Select(&ServiceInstance{ID: instanceID})
}

func (m Manager) CreateDatabase(dbName string) error {
	// NOTE: https://github.com/lib/pq/issues/694
	// dbName is generated automatically
	_, err := m.db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
	return err
}

func (m Manager) DropDatabase(dbName string) error {
	// NOTE: https://github.com/lib/pq/issues/694
	_, err := m.db.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, dbName))
	return err
}

func (m Manager) CreateServiceInstanceDetails(instanceID string, details brokerapi.ProvisionDetails) error {
	return m.db.Insert(&ServiceInstance{ID: instanceID, State: brokerapi.InProgress, Details: details})
}

func (m Manager) UpdateServiceInstanceDetails(instanceID string, state brokerapi.LastOperationState) error {
	return m.db.Update(&ServiceInstance{ID: instanceID, State: state})
}
