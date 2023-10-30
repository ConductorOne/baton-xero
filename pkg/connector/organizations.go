package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-xero/pkg/xero"
)

type orgResourceType struct {
	resourceType *v2.ResourceType
	client       *xero.Client
}

func (o *orgResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return o.resourceType
}

// Create a new connector resource for a Xero Organization.
func orgResource(ctx context.Context, org *xero.Organization) (*v2.Resource, error) {
	resource, err := resource.NewResource(
		org.Name,
		resourceTypeOrg,
		org.Id,
	)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (o *orgResourceType) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	orgs, err := o.client.GetOrganizations(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("xero-connector: failed to list orgs: %w", err)
	}

	var rv []*v2.Resource
	for _, org := range orgs {
		orgCopy := org

		or, err := orgResource(ctx, &orgCopy)
		if err != nil {
			return nil, "", nil, err
		}

		rv = append(rv, or)
	}

	return rv, "", nil, nil
}

func (o *orgResourceType) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *orgResourceType) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func orgBuilder(client *xero.Client) *orgResourceType {
	return &orgResourceType{
		resourceType: resourceTypeOrg,
		client:       client,
	}
}
