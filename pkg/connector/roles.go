package connector

import (
	"context"
	"fmt"
	"strings"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
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
func roleResource(role string) (*v2.Resource, error) {
	displayName := titleCase(role)

	profile := map[string]interface{}{
		"role_name": role,
	}

	resource, err := rs.NewRoleResource(
		displayName,
		resourceTypeRole,
		role,
		[]rs.RoleTraitOption{
			rs.WithRoleProfile(profile),
		},
	)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (r *roleResourceType) List(ctx context.Context, _ *v2.ResourceId, _ rs.SyncOpAttrs) ([]*v2.Resource, *rs.SyncOpResults, error) {
	var rv []*v2.Resource
	for _, r := range roles {
		rr, err := roleResource(r)
		if err != nil {
			return nil, nil, err
		}

		rv = append(rv, rr)
	}

	return rv, nil, nil
}

func (r *roleResourceType) Entitlements(_ context.Context, resource *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Entitlement, *rs.SyncOpResults, error) {
	var rv []*v2.Entitlement

	assignmentOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(resourceTypeUser),
		ent.WithDisplayName(fmt.Sprintf("%s Role", resource.DisplayName)),
		ent.WithDescription(fmt.Sprintf("%s role in Xero organization", resource.DisplayName)),
	}

	rv = append(rv, ent.NewAssignmentEntitlement(resource, resource.Id.Resource, assignmentOptions...))

	return rv, nil, nil
}

func (r *roleResourceType) Grants(ctx context.Context, resource *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Grant, *rs.SyncOpResults, error) {
	users, err := r.client.GetUsers(ctx, strings.ToUpper(resource.Id.Resource))
	if err != nil {
		return nil, nil, fmt.Errorf("xero-connector: failed to list users with role %s: %w", resource.DisplayName, err)
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

	return rv, nil, nil
}

func roleBuilder(client *xero.Client) *roleResourceType {
	return &roleResourceType{
		resourceType: resourceTypeRole,
		client:       client,
	}
}
