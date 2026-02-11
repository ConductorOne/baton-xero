package main

import (
	"context"

	"github.com/conductorone/baton-sdk/pkg/config"
	cfg "github.com/conductorone/baton-xero/pkg/config"
	"github.com/conductorone/baton-xero/pkg/connector"
)

var version = "dev"

func main() {
	ctx := context.Background()
	config.RunConnector(
		ctx,
		"baton-xero",
		version,
		cfg.Config,
		connector.New)
}
