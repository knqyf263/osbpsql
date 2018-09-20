package db

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/knqyf263/osbpsql/config"
	"github.com/pivotal-cf/brokerapi"
)

// ServiceInstance represents ServiceInstance
type ServiceInstance struct {
	ID      string
	State   brokerapi.LastOperationState
	Details brokerapi.ProvisionDetails
}

// ServiceBinding represents ServiceBinding
type ServiceBinding struct {
	ID      string
	State   brokerapi.LastOperationState
	Details brokerapi.BindDetails
}

// Manager represents database manager
type Manager struct {
	db *pg.DB
}

// New returns new manager instance
func New(conf config.DBConfig) *Manager {
	db := pg.Connect(&pg.Options{
		Addr: fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		User: conf.User,
	})
	return &Manager{db: db}
}

// Migrate migrates database
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

// FindServiceInstanceByID finds ServiceIntance by instance id
func (m Manager) FindServiceInstanceByID(instanceID string) (*ServiceInstance, error) {
	serviceInstance := &ServiceInstance{ID: instanceID}
	err := m.db.Select(serviceInstance)
	if err != nil {
		return &ServiceInstance{}, err
	}
	return serviceInstance, nil
}

// FindServiceBindingByID finds ServiceBinding by instance id
func (m Manager) FindServiceBindingByID(instanceID string) (*ServiceBinding, error) {
	serviceBinding := &ServiceBinding{ID: instanceID}
	err := m.db.Select(serviceBinding)
	if err != nil {
		return &ServiceBinding{}, err
	}
	return serviceBinding, nil
}

// CreateDatabase creates database
func (m Manager) CreateDatabase(dbName string) error {
	// NOTE: https://github.com/lib/pq/issues/694
	// dbName is generated automatically
	_, err := m.db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
	return err
}

// DropDatabase drops database
func (m Manager) DropDatabase(dbName string) error {
	// NOTE: https://github.com/lib/pq/issues/694
	_, err := m.db.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, dbName))
	return err
}

// CreateServiceInstanceDetails saves ServiceInstance details to DB
func (m Manager) CreateServiceInstanceDetails(instanceID string, state brokerapi.LastOperationState,
	details brokerapi.ProvisionDetails) error {
	return m.db.Insert(&ServiceInstance{ID: instanceID, State: state, Details: details})
}

// UpdateServiceInstanceDetails updates ServiceInstance details
func (m Manager) UpdateServiceInstanceDetails(instanceID string, state brokerapi.LastOperationState) error {
	return m.db.Update(&ServiceInstance{ID: instanceID, State: state})
}

// DeleteServiceInstanceDetails deletes ServiceInstance details from DB
func (m Manager) DeleteServiceInstanceDetails(instanceID string) error {
	return m.db.Delete(&ServiceInstance{ID: instanceID})
}

// CreateServiceBindingDetails saves ServiceBinding details to DB
func (m Manager) CreateServiceBindingDetails(instanceID string, state brokerapi.LastOperationState,
	details brokerapi.BindDetails) error {
	return m.db.Insert(&ServiceBinding{ID: instanceID, State: state, Details: details})
}

// UpdateServiceBindingDetails saves ServiceBinding details to DB
func (m Manager) UpdateServiceBindingDetails(instanceID string, state brokerapi.LastOperationState) error {
	return m.db.Update(&ServiceBinding{ID: instanceID, State: state})
}

// DeleteServiceBindingDetails deletes ServiceBinding details from DB
func (m Manager) DeleteServiceBindingDetails(instanceID string) error {
	return m.db.Delete(&ServiceBinding{ID: instanceID})
}
