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
	"golang.org/x/exp/slices"
)

const (
	readOnly         = "readonly"
	invoiceOnly      = "invoiceonly"
	standard         = "standard"
	financialAdvisor = "financialadvisor"
	managedClient    = "managedclient"
	cashbookClient   = "cashbookclient"
)

var orgRoles = []string{readOnly, invoiceOnly, standard, financialAdvisor, managedClient, cashbookClient}

type orgResourceType struct {
	resourceType *v2.ResourceType
	client       *xero.Client
}

func (o *orgResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return o.resourceType
}

// Create a new connector resource for a Xero User.
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

func (o *orgResourceType) List(ctx context.Context, parentID *v2.ResourceId, pt *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
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

func (o *orgResourceType) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	for _, role := range orgRoles {
		assignmentOptions := []ent.EntitlementOption{
			ent.WithGrantableTo(resourceTypeUser),
			ent.WithDisplayName(fmt.Sprintf("%s Role", titleCase(role))),
			ent.WithDescription(fmt.Sprintf("%s role in %s Xero organization", titleCase(role), resource.DisplayName)),
		}

		rv = append(rv, ent.NewAssignmentEntitlement(resource, role, assignmentOptions...))
	}

	return rv, "", nil, nil
}

func (o *orgResourceType) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	users, err := o.client.GetUsers(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("xero-connector: failed to list users: %w", err)
	}

	var rv []*v2.Grant
	for _, user := range users {
		// check if role is supported
		if !slices.Contains(orgRoles, strings.ToLower(user.Role)) {
			continue
		}

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

func orgBuilder(client *xero.Client) *orgResourceType {
	return &orgResourceType{
		resourceType: resourceTypeOrg,
		client:       client,
	}
}
