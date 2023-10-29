package xero

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ApiBase      = "api.xero.com"
	IdentityBase = "identity.xero.com"

	ApiEndpoint           = "/api.xro/2.0"
	ExchangeTokenEndpoint = "/connect/token"
	ConnectionsEndpoint   = "/connections"

	UsersEndpoint = "/Users"
	UserEndpoint  = "/Users/%s"
	OrgsEndpoint  = "/Organisations"
)

type Client struct {
	httpClient   *http.Client
	baseUrl      *url.URL
	token        string
	refreshToken string
	tenant       string
}

func NewClient(ctx context.Context, httpClient *http.Client, auth *Auth) (*Client, error) {
	// login if token is not present
	if auth.Token == "" {
		err := auth.Login(ctx, httpClient)
		if err != nil {
			return nil, fmt.Errorf("failed to login: %w", err)
		}
	}

	// obtain tenant id required for all requests
	tenantId, err := GetTenant(ctx, httpClient, auth.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant id: %w", err)
	}

	return &Client{
		httpClient:   httpClient,
		baseUrl:      &url.URL{Scheme: "https", Host: ApiBase, Path: ApiEndpoint},
		token:        auth.Token,
		refreshToken: auth.RefreshToken,
		tenant:       tenantId,
	}, nil
}

func (c *Client) joinURL(path string) *url.URL {
	newURL := *c.baseUrl
	newURL.Path += path

	return &newURL
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UsersResponse struct {
	Users []User `json:"users"`
}

// GetUsers returns all users under the team account.
func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	var usersResponse UsersResponse

	err := c.get(
		ctx,
		c.joinURL(UsersEndpoint),
		&usersResponse,
	)

	if err != nil {
		return nil, err
	}

	return usersResponse.Users, nil
}

type OrgResponse struct {
	Orgs []Organization `json:"Organisations"`
}

// GetOrganizations returns all organizations available.
func (c *Client) GetOrganizations(ctx context.Context) ([]Organization, error) {
	var orgsResponse OrgResponse

	err := c.get(
		ctx,
		c.joinURL(OrgsEndpoint),
		&orgsResponse,
	)

	if err != nil {
		return nil, err
	}

	return orgsResponse.Orgs, nil
}

func (c *Client) get(ctx context.Context, urlAddress *url.URL, resourceResponse interface{}) error {
	return c.doRequest(ctx, urlAddress, http.MethodGet, nil, resourceResponse)
}

func (c *Client) doRequest(
	ctx context.Context,
	urlAddress *url.URL,
	method string,
	data url.Values,
	resourceResponse interface{},
) error {
	var body strings.Reader

	if data != nil {
		encodedData := data.Encode()
		bodyReader := strings.NewReader(encodedData)
		body = *bodyReader
	}

	req, err := http.NewRequestWithContext(ctx, method, urlAddress.String(), &body)
	if err != nil {
		return err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("xero-tenant-id", c.tenant)

	rawResponse, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer rawResponse.Body.Close()

	if rawResponse.StatusCode >= 300 {
		return status.Error(codes.Code(rawResponse.StatusCode), "Request failed")
	}

	if err := json.NewDecoder(rawResponse.Body).Decode(resourceResponse); err != nil {
		return err
	}

	return nil
}
