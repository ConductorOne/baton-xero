package config

import (
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	AccessToken = field.StringField(
		"token",
		field.WithDescription("The Xero access token used to connect to the Xero API."),
		field.WithIsSecret(true),
		field.WithDisplayName("Access Token"),
	)
	RefreshToken = field.StringField(
		"refresh-token",
		field.WithDescription("The Xero refresh token used to exchange for a new access token."),
		field.WithIsSecret(true),
		field.WithDisplayName("Refresh Token"),
	)
	XeroClientId = field.StringField(
		"xero-client-id",
		field.WithDescription("The Xero client ID used to connect to the Xero API."),
		field.WithDisplayName("Xero Client ID"),
	)
	XeroClientSecret = field.StringField(
		"xero-client-secret",
		field.WithDescription("The Xero client secret used to connect to the Xero API."),
		field.WithIsSecret(true),
		field.WithDisplayName("Xero Client Secret"),
	)

	// FieldRelationships defines relationships between the fields listed in
	// Config that can be automatically validated.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run ./gen
var Config = field.NewConfiguration([]field.SchemaField{
	AccessToken,
	RefreshToken,
	XeroClientId,
	XeroClientSecret,
})

// ValidateConfig is run after the configuration is loaded, and should return an
// error if it isn't valid.
func ValidateConfig(cfg *Xero) error {
	isOAuthSet := (cfg.XeroClientId != "" && cfg.XeroClientSecret != "")
	isTokenSet := cfg.Token != ""
	isRefreshTokenSet := cfg.RefreshToken != ""

	if !isOAuthSet && !isTokenSet {
		return fmt.Errorf("either client id and secret or a token must be set, use --help for more information")
	}

	if isRefreshTokenSet && !isOAuthSet {
		return fmt.Errorf("refresh token requires client id and secret to be set, use --help for more information")
	}

	return nil
}
