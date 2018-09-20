package brokers

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/go-pg/pg"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pkg/errors"

	"github.com/knqyf263/osbpsql/brokers/models"
	"github.com/knqyf263/osbpsql/config"
	"github.com/knqyf263/osbpsql/db"
)

var (
	services         []brokerapi.Service
	serviceBrokerMap = map[string]models.ServiceBroker{}
)

// Register registers ServiceBroker
func Register(service brokerapi.Service, broker models.ServiceBroker) {
	services = append(services, service)
	serviceBrokerMap[service.ID] = broker
}

// PsqlServiceBroker is a brokerapi.ServiceBroker that can be used to generate an OSB compatible service broker.
type PsqlServiceBroker struct {
	Logger lager.Logger
	db     *db.Manager
}

// New creates a PsqlServiceBroker.
// Exactly one of PsqlServiceBroker or error will be nil when returned.
func New(cfg *config.BrokerConfig, logger lager.Logger) (self *PsqlServiceBroker, err error) {

	self = &PsqlServiceBroker{}
	self.Logger = logger
	self.db = db.New(cfg.DB)
	if self.db == nil {
		return nil, errors.New("Failed to initialize DB")
	}
	if err = self.db.Migrate(); err != nil {
		return nil, err
	}

	for _, broker := range serviceBrokerMap {
		broker.SetDB(self.db)
	}

	return self, nil
}

// Services lists services in the broker's catalog.
// It is called through the `GET /v2/catalog` endpoint or the `cf marketplace` command.
func (broker *PsqlServiceBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	return services, nil
}

func (broker *PsqlServiceBroker) getPlanFromID(serviceID, planID string) error {
	var service *brokerapi.Service
	for _, s := range services {
		if s.ID == serviceID {
			service = &s
			break
		}
	}
	if service == nil {
		return fmt.Errorf("unknown service id: %q", serviceID)
	}

	for _, plan := range service.Plans {
		if plan.ID == planID {
			return nil
		}
	}

	return fmt.Errorf("unknown plan id: %q", planID)
}

// Provision creates a new instance of a service.
// It is bound to the `PUT /v2/service_instances/:instance_id` endpoint and can be called using the `cf create-service` command.
func (broker *PsqlServiceBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, clientSupportsAsync bool) (brokerapi.ProvisionedServiceSpec, error) {
	// Only support assync
	if !clientSupportsAsync {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrAsyncRequired
	}

	broker.Logger.Info("Provisioning", lager.Data{
		"instanceId":         instanceID,
		"accepts_incomplete": clientSupportsAsync,
		"details":            details,
	})
	_, err := broker.db.FindServiceInstanceByID(instanceID)
	if err == nil {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	} else if err != pg.ErrNoRows {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Database error checking for existing instance: %s", err)
	}

	serviceID := details.ServiceID
	err = broker.getPlanFromID(serviceID, details.PlanID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	if err = broker.db.CreateServiceInstanceDetails(instanceID, brokerapi.InProgress, details); err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Error saving instance details to database: %s", err)
	}

	service := serviceBrokerMap[serviceID]
	go func() {
		if err = service.Provision(instanceID, details); err != nil {
			broker.db.UpdateServiceInstanceDetails(instanceID, brokerapi.Failed)
			return
		}
		broker.db.UpdateServiceInstanceDetails(instanceID, brokerapi.Succeeded)
	}()

	return brokerapi.ProvisionedServiceSpec{IsAsync: true, DashboardURL: ""}, nil
}

// Deprovision destroys an existing instance of a service.
// It is bound to the `DELETE /v2/service_instances/:instance_id` endpoint and can be called using the `cf delete-service` command.
func (broker *PsqlServiceBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, clientSupportsAsync bool) (brokerapi.DeprovisionServiceSpec, error) {
	response := brokerapi.DeprovisionServiceSpec{IsAsync: true}

	// Only support assync
	if !clientSupportsAsync {
		return response, brokerapi.ErrAsyncRequired
	}

	_, err := broker.db.FindServiceInstanceByID(instanceID)
	if err == pg.ErrNoRows {
		return response, brokerapi.ErrInstanceDoesNotExist
	} else if err != nil {
		return response, err
	}

	service := serviceBrokerMap[details.ServiceID]
	go func() {
		if err = service.Deprovision(ctx, instanceID, details); err != nil {
			broker.db.UpdateServiceInstanceDetails(instanceID, brokerapi.Failed)
			return
		}

		if err = broker.db.DeleteServiceInstanceDetailsG(instanceID); err != nil {
			return
		}

		_, err := broker.db.FindServiceBindingByID(instanceID)
		if err == nil {
			broker.db.DeleteServiceBindingDetailsG(instanceID)
		}
	}()

	return response, nil
}

// Bind creates an account with credentials to access an instance of a service.
// It is bound to the `PUT /v2/service_instances/:instance_id/service_bindings/:binding_id` endpoint and can be called using the `cf bind-service` command.
func (broker *PsqlServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	_, err := broker.db.FindServiceBindingByID(instanceID)
	if err == nil {
		return brokerapi.Binding{}, brokerapi.ErrBindingAlreadyExists
	} else if err != pg.ErrNoRows {
		return brokerapi.Binding{}, errors.Wrap(err, "Unexpected DB error")
	}

	service := serviceBrokerMap[details.ServiceID]
	binding, err := service.Bind(ctx, instanceID, bindingID, details)
	if err != nil {
		return brokerapi.Binding{}, errors.Wrap(err, "Failed to create ServiceBinding")
	}
	if err = broker.db.CreateServiceBindingDetails(instanceID, brokerapi.Succeeded, details); err != nil {
		return brokerapi.Binding{}, errors.Wrap(err, "Failed to save ServiceBinding details to DB")
	}
	return binding, nil
}

// Unbind destroys an account and credentials with access to an instance of a service.
// It is bound to the `DELETE /v2/service_instances/:instance_id/service_bindings/:binding_id` endpoint and can be called using the `cf unbind-service` command.
func (broker *PsqlServiceBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	_, err := broker.db.FindServiceBindingByID(instanceID)
	if err == pg.ErrNoRows {
		return brokerapi.ErrBindingDoesNotExist
	} else if err != nil {
		return errors.Wrap(err, "Unexpected DB error")
	}

	service := serviceBrokerMap[details.ServiceID]
	if err := service.Unbind(ctx, instanceID, bindingID, details); err != nil {

	}
	return broker.db.DeleteServiceBindingDetailsG(instanceID)
}

// LastOperation returns the state of the last requested operation.
func (broker *PsqlServiceBroker) LastOperation(ctx context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	serviceInstance, err := broker.db.FindServiceInstanceByID(instanceID)
	if err == pg.ErrNoRows {
		return brokerapi.LastOperation{}, brokerapi.ErrInstanceDoesNotExist
	} else if err != nil {
		return brokerapi.LastOperation{}, err
	}
	return brokerapi.LastOperation{State: serviceInstance.State}, nil
}

// Update a service instance plan.
// This functionality is not implemented and will return an error indicating that plan changes are not supported.
func (broker *PsqlServiceBroker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, brokerapi.ErrPlanChangeNotSupported
}
