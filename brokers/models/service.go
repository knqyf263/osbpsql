package models

import (
	"context"

	"github.com/knqyf263/osbpsql/db"
	"github.com/pivotal-cf/brokerapi"
)

// ServiceBroker : ServiceBroker
type ServiceBroker interface {
	SetDB(db *db.Manager)
	Provision(instanceID string, details brokerapi.ProvisionDetails) error
	Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error)
	Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error
	Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails) error
}
