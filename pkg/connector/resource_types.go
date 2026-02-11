package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
)

var (
	resourceTypeOrg = &v2.ResourceType{
		Id:          "org",
		DisplayName: "Organization",
		Annotations: annotations.New(&v2.SkipEntitlementsAndGrants{}),
	}
	resourceTypeUser = &v2.ResourceType{
		Id:          "user",
		DisplayName: "User",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_USER,
		},
		Annotations: annotations.New(&v2.SkipEntitlementsAndGrants{}),
	}
	resourceTypeRole = &v2.ResourceType{
		Id:          "role",
		DisplayName: "Role",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_ROLE,
		},
	}
)
