package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	AccessToken = field.StringField(
		"token",
		field.WithDescription("The Xero access token used to connect to the Xero API."),
		field.WithIsSecret(true),
		field.WithDisplayName("Access Token"),
		field.WithExportTarget(field.ExportTargetCLIOnly),
	)
	RefreshToken = field.StringField(
		"refresh-token",
		field.WithDescription("The Xero refresh token used to exchange for a new access token."),
		field.WithIsSecret(true),
		field.WithDisplayName("Refresh Token"),
		field.WithExportTarget(field.ExportTargetCLIOnly),
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
)

//go:generate go run ./gen
var Config = field.NewConfiguration([]field.SchemaField{
	AccessToken,
	RefreshToken,
	XeroClientId,
	XeroClientSecret,
},
	field.WithConnectorDisplayName("Xero"),
	field.WithIconUrl("/static/app-icons/xero.svg"),
	field.WithHelpUrl("/docs/baton/xero"),
)
