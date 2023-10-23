package xero

import (
	"context"
	"encoding/base64"
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
	httpClient *http.Client
	baseUrl    *url.URL
	token      string
	tenant     string
}

type Auth struct {
	Token        string
	RefreshToken string
	ClientId     string
	ClientSecret string
}

func NewAuth(token, clientId, clientSecret string) *Auth {
	return &Auth{
		Token:        token,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}

func NewClient(ctx context.Context, httpClient *http.Client, auth *Auth) (*Client, error) {
	// if using custom integration with client id and secret -> exchange for token
	if auth.Token == "" {
		t, rt, err := Login(ctx, httpClient, auth.ClientId, auth.ClientSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to login: %w", err)
		}

		auth.Token = t
		auth.RefreshToken = rt
	}

	// obtain tenant id required for all requests
	tenantId, err := GetTenant(ctx, httpClient, auth.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant id: %w", err)
	}

	return &Client{
		httpClient: httpClient,
		baseUrl:    &url.URL{Scheme: "https", Host: ApiBase, Path: ApiEndpoint},
		token:      auth.Token,
		tenant:     tenantId,
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

func Login(ctx context.Context, httpClient *http.Client, clientId, clientSecret string) (string, string, error) {
	baseUrl := &url.URL{Scheme: "https", Host: IdentityBase, Path: ExchangeTokenEndpoint}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "openid email profile offline_access")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseUrl.String(),
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", constructBasicAuth(clientId, clientSecret))

	rawResponse, err := httpClient.Do(req)
	if err != nil {
		return "", "", err
	}

	defer rawResponse.Body.Close()

	if rawResponse.StatusCode >= 300 {
		return "", "", status.Error(codes.Code(rawResponse.StatusCode), "Request failed")
	}

	var res TokenResponse

	if err := json.NewDecoder(rawResponse.Body).Decode(&res); err != nil {
		return "", "", err
	}

	return res.AccessToken, res.RefreshToken, nil
}

type Connection struct {
	Id         string `json:"id"`
	TenantId   string `json:"tenantId"`
	TenantName string `json:"tenantName"`
}

func GetTenant(ctx context.Context, httpClient *http.Client, token string) (string, error) {
	conns, err := getConnections(ctx, httpClient, token)
	if err != nil {
		return "", err
	}

	if len(conns) == 0 {
		return "", fmt.Errorf("no connections found")
	}

	return conns[0].TenantId, nil
}

func getConnections(ctx context.Context, httpClient *http.Client, token string) ([]Connection, error) {
	baseUrl := &url.URL{Scheme: "https", Host: ApiBase, Path: ConnectionsEndpoint}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		baseUrl.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	rawResponse, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer rawResponse.Body.Close()

	if rawResponse.StatusCode >= 300 {
		return nil, status.Error(codes.Code(rawResponse.StatusCode), "Request failed")
	}

	var res []Connection

	if err := json.NewDecoder(rawResponse.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func constructBasicAuth(clientId, clientSecret string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
	return "Basic " + encoded
}
