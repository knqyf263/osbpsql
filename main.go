package main

import (
	"net/http"
	"os"
	"strings"
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/knqyf263/osbpsql/brokers"
	_ "github.com/knqyf263/osbpsql/brokers/psql-db"
	"github.com/knqyf263/osbpsql/config"
	_ "github.com/lib/pq"
	"github.com/pivotal-cf/brokerapi"
)

var (
	instanceState sync.Map
	bindState     sync.Map

	PropertyToEnvReplacer = strings.NewReplacer(".", "_", "-", "_")
)

const (
	EnvironmentVarPrefix = "osbpsql"
)

func main() {
	logger := lager.NewLogger("osbpsql")
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	cfg := config.NewBrokerConfigFromEnv()

	// init api
	logger.Info("Serving", lager.Data{
		"port":     cfg.Port,
		"username": cfg.Credentials.Username,
	})

	serviceBroker, err := brokers.New(cfg, logger)
	if err != nil {
		logger.Fatal("Error initializing service broker: %s", err)
	}

	brokerAPI := brokerapi.New(serviceBroker, logger, cfg.Credentials)
	http.Handle("/", brokerAPI)
	http.ListenAndServe(":"+cfg.Port, nil)
}

// func run(c *cli.Context) error {

// 	e := echo.New()

// 	e.Use(middleware.Logger())
// 	e.Use(middleware.Recover())
// 	e.Use(interceptor.BasicAuth())

// 	h := NewHandler()

// 	e.GET("/v2/catalog", func(c echo.Context) error {
// 		catalog, err := h.catalog()
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get catalog information")
// 		}
// 		return c.JSON(http.StatusOK, catalog)
// 	})
// 	e.PUT("/v2/service_instances/:instance_id", func(c echo.Context) (err error) {
// 		instanceID := c.Param("instance_id")
// 		if err = h.provision(instanceID); err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, map[string]error{"description": err})
// 		}
// 		return c.JSON(http.StatusAccepted, map[string]string{"operation": osb.Provisioning})

// 	})

// 	e.DELETE("/v2/service_instances/:instance_id", func(c echo.Context) (err error) {
// 		instanceID := c.Param("instance_id")
// 		if err = h.deprovision(instanceID); err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, map[string]error{"description": err})
// 		}
// 		return c.JSON(http.StatusAccepted, map[string]string{"operation": osb.Deprovisioning})
// 	})

// 	e.GET("/v2/service_instances/:instance_id/last_operation", func(c echo.Context) (err error) {
// 		instanceID := c.Param("instance_id")
// 		state, err := h.lastOperation(instanceID)
// 		if err == osb.ErrInstanceIDNotFound {
// 			return c.JSON(http.StatusGone, map[string]string{"description": fmt.Sprintf("%s does not exist", instanceID)})
// 		}
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, map[string]error{"description": err})
// 		}
// 		return c.JSON(http.StatusOK, map[string]string{"state": state.String()})
// 	})

// 	e.PUT("/v2/service_instances/:instance_id/service_bindings/:binding_id", func(c echo.Context) (err error) {
// 		instanceID := c.Param("instance_id")
// 		bindingID := c.Param("binding_id")
// 		credentials, err := h.bind(instanceID, bindingID)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, map[string]error{"description": err})
// 		}
// 		return c.JSON(http.StatusOK, map[string]interface{}{"credentials": credentials})
// 	})

// 	e.DELETE("/v2/service_instances/:instance_id/service_bindings/:binding_id", func(c echo.Context) (err error) {
// 		instanceID := c.Param("instance_id")
// 		bindingID := c.Param("binding_id")
// 		err = h.unbind(instanceID, bindingID)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, map[string]error{"description": err})
// 		}
// 		return c.JSON(http.StatusOK, nil)
// 	})

// 	return e.Start(":8080")
// }

// type Handler struct {
// 	db *sql.DB
// }

// func NewHandler() *Handler {
// 	conninfo := "user=postgres password=postgres host=127.0.0.1 sslmode=disable"
// 	db, _ := sql.Open("postgres", conninfo)

// 	return &Handler{db: db}
// }

// func (h Handler) catalog() (*osb.Catalog, error) {
// 	currentCatalog := &osb.Catalog{
// 		Services: []*osb.Service{
// 			&osb.Service{
// 				Name:        "postgresql",
// 				ID:          "c4f353f3-8a59-437d-b4af-6a6f856248db",
// 				Description: "Test Service",
// 				Bindable:    true,
// 				Plans: []*osb.Plan{
// 					testPlan,
// 				},
// 			},
// 		},
// 	}
// 	return currentCatalog, nil
// }

// func (h Handler) provision(instanceID string) error {
// 	fmt.Println(instanceID)
// 	if state, ok := instanceState.Load(instanceID); ok {
// 		return fmt.Errorf("%s already exists (State=%s)", instanceID, state)
// 	}
// 	instanceState.Store(instanceID, osb.StateInProgress)

// 	go func() {
// 		// https://github.com/lib/pq/issues/694
// 		if _, err := h.db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, instanceID)); err != nil {
// 			instanceState.Store(instanceID, osb.StateFailed)
// 			fmt.Println(err)
// 			return
// 		}
// 		instanceState.Store(instanceID, osb.StateSucceeded)
// 	}()
// 	return nil
// }

// func (h Handler) deprovision(instanceID string) error {
// 	if _, ok := instanceState.Load(instanceID); !ok {
// 		return fmt.Errorf("%s does not exist", instanceID)
// 	}
// 	instanceState.Store(instanceID, osb.StateGone)

// 	go func() {
// 		// NOTE: https://github.com/lib/pq/issues/694
// 		if _, err := h.db.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, instanceID)); err != nil {
// 			instanceState.Store(instanceID, osb.StateFailed)
// 			fmt.Println(err)
// 			return
// 		}
// 		instanceState.Delete(instanceID)
// 	}()
// 	return nil
// }

// func (h Handler) lastOperation(instanceID string) (osb.ProvisioningState, error) {
// 	stateInterface, ok := instanceState.Load(instanceID)
// 	if !ok {
// 		return -1, osb.ErrInstanceIDNotFound
// 	}
// 	state, ok := stateInterface.(osb.ProvisioningState)
// 	if !ok {
// 		return -1, errors.New("assertion error")
// 	}
// 	return state, nil
// }

// func (h Handler) bind(instanceID, bindingID string) (credentials map[string]string, err error) {
// 	return map[string]string{
// 		"database": instanceID,
// 	}, nil
// }

// func (h Handler) unbind(instanceID, bindingID string) (err error) {
// 	return nil
// }
