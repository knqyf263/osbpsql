package config

import (
	cli "gopkg.in/urfave/cli.v2"
)

type cliConfig struct {
	BasicAuthUsername string
	BasicAuthPassword string
}

// Config represetns configurations
var Config = &cliConfig{}

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "basic-auth-username",
		Usage:       "BASIC auth username",
		EnvVars:     []string{"BASIC_AUTH_USERNAME"},
		Destination: &Config.BasicAuthUsername,
	},
	&cli.StringFlag{
		Name:        "basic-auth-password",
		Usage:       "BASIC auth password",
		EnvVars:     []string{"BASIC_AUTH_PASSWORD"},
		Destination: &Config.BasicAuthPassword,
	},
}
