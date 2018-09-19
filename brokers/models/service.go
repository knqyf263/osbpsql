package models

import (
	"context"

	"github.com/knqyf263/osbpsql/db"
	"github.com/pivotal-cf/brokerapi"
)

type ServiceBroker interface {
	SetDB(db *db.Manager)
	Provision(instanceID string, details brokerapi.ProvisionDetails) error
	Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error)
	// BuildInstanceCredentials(bindRecord ServiceBindingCredentials, instanceRecord ServiceInstanceDetails) (map[string]string, error)
	Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error
	// Deprovision(instance ServiceInstanceDetails, details brokerapi.DeprovisionDetails) error
	Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails) error
	// PollInstance(instanceID string) (bool, error)
	// LastOperationWasDelete(instanceID string) (bool, error)
	// ProvisionsAsync() bool
	// DeprovisionsAsync() bool
}
