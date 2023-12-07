package main

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/spf13/cobra"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig `mapstructure:",squash"` // Puts the base config options in the same place as the connector options

	AccessToken      string `mapstructure:"token"`
	RefreshToken     string `mapstructure:"refresh-token"`
	XeroClientId     string `mapstructure:"xero-client-id"`
	XeroClientSecret string `mapstructure:"xero-client-secret"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	isOAuthSet := (cfg.XeroClientId != "" && cfg.XeroClientSecret != "")
	isTokenSet := cfg.AccessToken != ""
	isRefreshTokenSet := cfg.RefreshToken != ""

	if !isOAuthSet && !isTokenSet {
		return fmt.Errorf("either client id and secret or a token must be set, use --help for more information")
	}

	if isRefreshTokenSet && !isOAuthSet {
		return fmt.Errorf("refresh token requires client id and secret to be set, use --help for more information")
	}

	return nil
}

// cmdFlags sets the cmdFlags required for the connector.
func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("token", "", "The Xero access token used to connect to the Xero API. ($BATON_TOKEN)")
	cmd.PersistentFlags().String("refresh-token", "", "The Xero refresh token used to exchange for a new access token. ($BATON_REFRESH_TOKEN)")
	cmd.PersistentFlags().String("xero-client-id", "", "The Xero client ID used to connect to the Xero API. ($BATON_XERO_CLIENT_ID)")
	cmd.PersistentFlags().String("xero-client-secret", "", "The Xero client secret used to connect to the Xero API. ($BATON_XERO_CLIENT_SECRET)")
}
