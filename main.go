package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/knqyf263/osbpsql/config"
	"github.com/knqyf263/osbpsql/interceptor"
	"github.com/knqyf263/osbpsql/osb"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	cli "gopkg.in/urfave/cli.v2"
)

var (
	appName      = "open-service-broker-framework"
	appUsage     = "An implementation of the Open Service Broker API"
	appCopyright = "Copyright (C) 2018 Teppei Fukuda."

	testPlan = &osb.Plan{
		ID:          "7992D02F-83F5-4F1D-AD2B-A31C4757D031",
		Name:        "dummy",
		Description: "DB",
		// Metadata:    &osb.Metadata{},
	}
)

func main() {
	app := &cli.App{
		Name:                  appName,
		Usage:                 appUsage,
		HelpName:              appName,
		Copyright:             appCopyright,
		EnableShellCompletion: true,
		Version:               "0.0.1",
		CommandNotFound:       cmdNotFound,
		Flags:                 config.Flags,
		Action:                run,
	}
	cli.InitCompletionFlag.Hidden = true

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func cmdNotFound(c *cli.Context, command string) {
	fmt.Fprintf(
		os.Stderr,
		"%s: '%s' is not a %s command. See '%s --help'\n",
		c.App.Name,
		command,
		c.App.Name,
		c.App.Name,
	)
	os.Exit(1)
}

func run(c *cli.Context) error {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(interceptor.BasicAuth())

	h := Handler{}
	e.GET("/v2/catalog", func(c echo.Context) error {
		catalog, err := h.catalog()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get catalog information")
		}
		return c.JSON(http.StatusOK, catalog)
	})
	e.PUT("/v2/service_instances/:instance_id", func(c echo.Context) error {
		instanceID := c.Param("instance_id")
		fmt.Println(instanceID)
		c.String(http.StatusAccepted, "")
		return h.provision()
	})

	e.DELETE("/v2/service_instances/:instance_id", func(c echo.Context) error {
		instanceID := c.Param("instance_id")
		fmt.Println(instanceID)
		return nil
	})

	return e.Start(":8080")
}

type Handler struct {
}

func (h Handler) catalog() (*osb.Catalog, error) {
	currentCatalog := &osb.Catalog{
		Services: []*osb.Service{
			&osb.Service{
				Name:        "postgresql",
				ID:          "c4f353f3-8a59-437d-b4af-6a6f856248db",
				Description: "Test Service",
				Plans: []*osb.Plan{
					testPlan,
				},
			},
		},
	}
	return currentCatalog, nil
}

func (h Handler) provision() error {
	return nil
}
