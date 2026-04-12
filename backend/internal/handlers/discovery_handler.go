package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DiscoveryHandler OAuth Discovery 端点处理
type DiscoveryHandler struct {
	issuerURL string
}

// NewDiscoveryHandler 创建 Discovery 处理器
func NewDiscoveryHandler(issuerURL string) *DiscoveryHandler {
	return &DiscoveryHandler{issuerURL: issuerURL}
}

// OAuthAuthorizationServerMetadata OAuth Authorization Server Metadata (RFC 8414)
type OAuthAuthorizationServerMetadata struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint              string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	TokenEndpointAuthSigningAlgValuesSupported []string `json:"token_endpoint_auth_signing_alg_values_supported"`
	UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
	JwksURI                           string   `json:"jwks_uri"`
	ScopesSupported                   []string `json:"scopes_supported"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	GrantTypesSupported               []string `json:"grant_types_supported"`
	PATGrantTypeSupported             bool     `json:"pat_grant_type_supported"`
}

// OpenIDConfiguration OpenID Connect Discovery 文档
type OpenIDConfiguration struct {
	Issuer                              string   `json:"issuer"`
	AuthorizationEndpoint                string   `json:"authorization_endpoint"`
	TokenEndpoint                       string   `json:"token_endpoint"`
	UserinfoEndpoint                    string   `json:"userinfo_endpoint"`
	JwksURI                             string   `json:"jwks_uri"`
	ScopesSupported                     []string `json:"scopes_supported"`
	ResponseTypesSupported              []string `json:"response_types_supported"`
	GrantTypesSupported                 []string `json:"grant_types_supported"`
	TokenEndpointAuthMethodsSupported   []string `json:"token_endpoint_auth_methods_supported"`
}

// GetOAuthAuthorizationServer .well-known/oauth-authorization-server 端点
// 返回 OAuth Authorization Server Metadata
func (h *DiscoveryHandler) GetOAuthAuthorizationServer(c *gin.Context) {
	metadata := OAuthAuthorizationServerMetadata{
		Issuer:                            h.issuerURL,
		AuthorizationEndpoint:              h.issuerURL + "/oauth/authorize",
		TokenEndpoint:                     h.issuerURL + "/oauth/token",
		TokenEndpointAuthMethodsSupported: []string{"client_secret_basic", "client_secret_post"},
		TokenEndpointAuthSigningAlgValuesSupported: []string{"HS256"},
		UserinfoEndpoint:                  h.issuerURL + "/oauth/userinfo",
		JwksURI:                           h.issuerURL + "/oauth/jwks",
		ScopesSupported: []string{
			"task:read", "task:write",
			"project:read", "project:write",
			"comment:read", "comment:write",
		},
		ResponseTypesSupported: []string{"token"},
		GrantTypesSupported:    []string{"client_credentials", "urn:ietf:params:oauth:grant-type:pat"},
		PATGrantTypeSupported: true,
	}

	c.JSON(http.StatusOK, metadata)
}

// GetOpenIDConfiguration .well-known/openid-configuration 端点
// 返回 OpenID Connect Discovery 文档
func (h *DiscoveryHandler) GetOpenIDConfiguration(c *gin.Context) {
	config := OpenIDConfiguration{
		Issuer:                       h.issuerURL,
		AuthorizationEndpoint:         h.issuerURL + "/oauth/authorize",
		TokenEndpoint:                h.issuerURL + "/oauth/token",
		UserinfoEndpoint:             h.issuerURL + "/oauth/userinfo",
		JwksURI:                      h.issuerURL + "/oauth/jwks",
		ScopesSupported: []string{
			"openid", "profile",
			"task:read", "task:write",
			"project:read", "project:write",
			"comment:read", "comment:write",
		},
		ResponseTypesSupported:        []string{"token"},
		GrantTypesSupported:           []string{"client_credentials", "urn:ietf:params:oauth:grant-type:pat"},
		TokenEndpointAuthMethodsSupported: []string{"client_secret_basic", "client_secret_post"},
	}

	c.JSON(http.StatusOK, config)
}