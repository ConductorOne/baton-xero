package xero

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var DefaultScopes = []string{"openid", "email", "profile", "offline_access", "accounting.settings", "accounting.transactions"}

type Auth struct {
	Token        string
	RefreshToken string
	ClientId     string
	ClientSecret string
}

func NewAuth(token, refreshToken, clientId, clientSecret string) *Auth {
	return &Auth{
		Token:        token,
		RefreshToken: refreshToken,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}

func (a *Auth) Login(ctx context.Context, httpClient *http.Client) error {
	if a.RefreshToken == "" {
		// login to obtain new token and refresh token
		t, rt, err := ClientCredentialsFlow(ctx, httpClient, a.ClientId, a.ClientSecret)
		if err != nil {
			return fmt.Errorf("failed to login: %w", err)
		}

		a.Token = t
		a.RefreshToken = rt
	} else {
		// use refresh token to obtain new token if present
		t, rt, err := RefreshTokenFlow(ctx, httpClient, a.RefreshToken, a.ClientId, a.ClientSecret)
		if err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}

		a.Token = t
		a.RefreshToken = rt
	}

	return nil
}

func ClientCredentialsFlow(ctx context.Context, httpClient *http.Client, clientId, clientSecret string) (string, string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("scope", strings.Join(DefaultScopes, " "))

	t, rt, err := exchangeToken(ctx, httpClient, &data, &Auth{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	})
	if err != nil {
		return "", "", fmt.Errorf("client_credentials flow failed: %w", err)
	}

	return t, rt, nil
}

func RefreshTokenFlow(ctx context.Context, httpClient *http.Client, refreshToken, clientId, clientSecret string) (string, string, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("scope", strings.Join(DefaultScopes, " "))

	t, rt, err := exchangeToken(ctx, httpClient, &data, &Auth{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	})
	if err != nil {
		return "", "", fmt.Errorf("refresh_token flow failed: %w", err)
	}

	return t, rt, nil
}

func exchangeToken(ctx context.Context, httpClient *http.Client, data *url.Values, auth *Auth) (string, string, error) {
	baseUrl := &url.URL{Scheme: "https", Host: IdentityBase, Path: ExchangeTokenEndpoint}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseUrl.String(),
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(auth.ClientId, auth.ClientSecret)

	rawResponse, err := httpClient.Do(req)
	if err != nil {
		return "", "", err
	}

	defer rawResponse.Body.Close()

	if rawResponse.StatusCode >= 300 {
		body := ""
		if rawResponse.Body != nil {
			b, _ := io.ReadAll(rawResponse.Body)
			if b != nil {
				body = string(b)
			}
		}
		if body == "" {
			body = "no error body"
		}
		return "", "", status.Error(codes.Code(rawResponse.StatusCode), fmt.Sprintf("Request failed: %s", body)) //nolint:gosec // disable G115
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
		return nil, status.Error(codes.Code(rawResponse.StatusCode), "Request failed") //nolint:gosec // disable G115
	}

	var res []Connection

	if err := json.NewDecoder(rawResponse.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}
