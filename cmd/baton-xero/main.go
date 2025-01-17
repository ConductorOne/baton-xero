package main

import (
	"context"
	"fmt"
	"os"

	configSchema "github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/conductorone/baton-xero/pkg/connector"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	version       = "dev"
	connectorName = "baton-xero"
)

var (
	AccessToken = field.StringField(
		"token",
		field.WithDescription("The Xero access token used to connect to the Xero API"),
	)
	RefreshToken = field.StringField(
		"refresh-token",
		field.WithDescription("The Xero refresh token used to exchange for a new access token"),
	)
	XeroClientId = field.StringField(
		"xero-client-id",
		field.WithDescription("The Xero client ID used to connect to the Xero API"),
	)
	XeroClientSecret = field.StringField(
		"xero-client-secret",
		field.WithDescription("The Xero client secret used to connect to the Xero API"),
	)
	configurationFields = []field.SchemaField{AccessToken, RefreshToken, XeroClientId, XeroClientSecret}
	fieldRelationships  = []field.SchemaFieldRelationship{
		field.FieldsRequiredTogether(XeroClientId, XeroClientSecret),
		field.FieldsDependentOn([]field.SchemaField{RefreshToken}, []field.SchemaField{XeroClientId}),
		field.FieldsMutuallyExclusive(AccessToken, XeroClientId),
		field.FieldsAtLeastOneUsed(AccessToken, XeroClientId),
	}
)

func main() {
	ctx := context.Background()

	_, cmd, err := configSchema.DefineConfiguration(ctx,
		connectorName,
		getConnector,
		field.NewConfiguration(configurationFields, fieldRelationships...),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, cfg *viper.Viper) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	xeroConnector, err := connector.New(ctx, cfg)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	connector, err := connectorbuilder.NewConnector(ctx, xeroConnector)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	return connector, nil
}
