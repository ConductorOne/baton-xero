package connector

import (
	"context"
	"fmt"
	"strings"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-xero/pkg/xero"
)

const (
	readOnly         = "readonly"
	invoiceOnly      = "invoiceonly"
	standard         = "standard"
	financialAdvisor = "financialadvisor"
	managedClient    = "managedclient"
	cashbookClient   = "cashbookclient"
)

var roles = []string{readOnly, invoiceOnly, standard, financialAdvisor, managedClient, cashbookClient}

type roleResourceType struct {
	resourceType *v2.ResourceType
	client       *xero.Client
}

func (r *roleResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return r.resourceType
}

// Create a new connector resource for a Xero Role.
func roleResource(ctx context.Context, role string) (*v2.Resource, error) {
	displayName := titleCase(role)

	profile := map[string]interface{}{
		"role_name": role,
	}

	resource, err := resource.NewRoleResource(
		displayName,
		resourceTypeRole,
		role,
		[]resource.RoleTraitOption{
			resource.WithRoleProfile(profile),
		},
	)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (r *roleResourceType) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var rv []*v2.Resource
	for _, r := range roles {
		rr, err := roleResource(ctx, r)
		if err != nil {
			return nil, "", nil, err
		}

		rv = append(rv, rr)
	}

	return rv, "", nil, nil
}

func (r *roleResourceType) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	assignmentOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(resourceTypeUser),
		ent.WithDisplayName(fmt.Sprintf("%s Role", resource.DisplayName)),
		ent.WithDescription(fmt.Sprintf("%s role in Xero organization", resource.DisplayName)),
	}

	rv = append(rv, ent.NewAssignmentEntitlement(resource, resource.Id.Resource, assignmentOptions...))

	return rv, "", nil, nil
}

func (r *roleResourceType) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	users, err := r.client.GetUsers(ctx, strings.ToUpper(resource.Id.Resource))
	if err != nil {
		return nil, "", nil, fmt.Errorf("xero-connector: failed to list users with role %s: %w", resource.DisplayName, err)
	}

	var rv []*v2.Grant
	for _, user := range users {
		rv = append(rv, grant.NewGrant(
			resource,
			strings.ToLower(user.Role),
			&v2.ResourceId{
				ResourceType: resourceTypeUser.Id,
				Resource:     user.Id,
			},
		))
	}

	return rv, "", nil, nil
}

func roleBuilder(client *xero.Client) *roleResourceType {
	return &roleResourceType{
		resourceType: resourceTypeRole,
		client:       client,
	}
}
