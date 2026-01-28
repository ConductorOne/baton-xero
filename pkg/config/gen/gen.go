package main

import (
	cfg "github.com/conductorone/baton-xero/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("xero", cfg.Config)
}
