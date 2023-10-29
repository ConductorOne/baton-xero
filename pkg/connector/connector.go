package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/conductorone/baton-xero/pkg/xero"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Xero struct {
	client *xero.Client
}

func (x *Xero) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		orgBuilder(x.client),
		userBuilder(x.client),
	}
}

// Metadata returns metadata about the connector.
func (x *Xero) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Xero",
		Description: "Connector syncing Xero organizations and their members to Baton",
	}, nil
}

// Validate hits the Xero API to validate that the configured credentials are valid and compatible.
func (x *Xero) Validate(ctx context.Context) (annotations.Annotations, error) {
	// should be able to list users
	_, err := x.client.GetUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Provided credentials are invalid")
	}

	return nil, nil
}

// New returns the Xero connector.
func New(ctx context.Context, clientId, clientSecret, token, refreshToken string) (*Xero, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	client, err := xero.NewClient(
		ctx,
		httpClient,
		xero.NewAuth(token, refreshToken, clientId, clientSecret),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &Xero{client}, nil
}
