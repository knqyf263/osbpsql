package broker

import (
	"context"
	"errors"

	"code.cloudfoundry.org/lager"
	"github.com/knqyf263/osbpsql/config"
	"github.com/knqyf263/osbpsql/db"
	"github.com/pivotal-cf/brokerapi"
)

// PsqlServiceBroker is a brokerapi.ServiceBroker that can be used to generate an OSB compatible service broker.
type PsqlServiceBroker struct {
	Logger lager.Logger
	db     *db.Manager
}

// New creates a PsqlServiceBroker.
// Exactly one of PsqlServiceBroker or error will be nil when returned.
func New(cfg *config.BrokerConfig, logger lager.Logger) (*PsqlServiceBroker, error) {

	self := PsqlServiceBroker{}
	self.Logger = logger
	self.db = db.New(cfg.DB)
	if self.db == nil {
		return nil, errors.New("Failed to initialize DB")
	}
	if err := self.db.Migrate(); err != nil {
		return nil, err
	}
	if err := self.db.CreateDatabase("foo"); err != nil {
		return nil, err
	}
	// self.Catalog = cfg.Catalog

	return &self, nil
}

// Services lists services in the broker's catalog.
// It is called through the `GET /v2/catalog` endpoint or the `cf marketplace` command.
func (broker *PsqlServiceBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	svcs := []brokerapi.Service{}

	return svcs, nil
}

// Provision creates a new instance of a service.
// It is bound to the `PUT /v2/service_instances/:instance_id` endpoint and can be called using the `cf create-service` command.
func (broker *PsqlServiceBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, clientSupportsAsync bool) (brokerapi.ProvisionedServiceSpec, error) {
	return brokerapi.ProvisionedServiceSpec{IsAsync: true, DashboardURL: ""}, nil
}

// Deprovision destroys an existing instance of a service.
// It is bound to the `DELETE /v2/service_instances/:instance_id` endpoint and can be called using the `cf delete-service` command.
func (broker *PsqlServiceBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, clientSupportsAsync bool) (brokerapi.DeprovisionServiceSpec, error) {
	response := brokerapi.DeprovisionServiceSpec{IsAsync: true}
	return response, nil
}

// Bind creates an account with credentials to access an instance of a service.
// It is bound to the `PUT /v2/service_instances/:instance_id/service_bindings/:binding_id` endpoint and can be called using the `cf bind-service` command.
func (broker *PsqlServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	return brokerapi.Binding{}, nil
}

// Unbind destroys an account and credentials with access to an instance of a service.
// It is bound to the `DELETE /v2/service_instances/:instance_id/service_bindings/:binding_id` endpoint and can be called using the `cf unbind-service` command.
func (broker *PsqlServiceBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	return nil
}

// Unbind destroys an account and credentials with access to an instance of a service.
// It is bound to the `GET /v2/service_instances/:instance_id/last_operation` endpoint.
// It is called by `cf create-service` or `cf delete-service` if the operation was asynchronous.
func (broker *PsqlServiceBroker) LastOperation(ctx context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{State: brokerapi.InProgress}, nil
}

// Update a service instance plan.
// This functionality is not implemented and will return an error indicating that plan changes are not supported.
func (broker *PsqlServiceBroker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, brokerapi.ErrPlanChangeNotSupported
}
