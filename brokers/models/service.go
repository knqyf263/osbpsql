package models

import (
	"context"

	"github.com/knqyf263/osbpsql/db"
	"github.com/pivotal-cf/brokerapi"
)

type ServiceBroker interface {
	SetDB(db *db.Manager)
	Provision(instanceID string, details brokerapi.ProvisionDetails) error
	// Bind(instanceID, bindingID string, details brokerapi.BindDetails) (ServiceBindingCredentials, error)
	// BuildInstanceCredentials(bindRecord ServiceBindingCredentials, instanceRecord ServiceInstanceDetails) (map[string]string, error)
	// Unbind(details ServiceBindingCredentials) error
	// Deprovision(instance ServiceInstanceDetails, details brokerapi.DeprovisionDetails) error
	Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails) error
	// PollInstance(instanceID string) (bool, error)
	// LastOperationWasDelete(instanceID string) (bool, error)
	// ProvisionsAsync() bool
	// DeprovisionsAsync() bool
}
