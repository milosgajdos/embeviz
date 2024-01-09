package qdrant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	DefaultBaseURL = "http://localhost:6333"
)

type Options struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// Option is functional graph option.
type Option func(*Options)

// HTTPClient is qdrant HTTP API client.
type HTTPClient struct {
	opts Options
}

// NewHTTPClient creates a new qdrant HTTP API client and returns it.
func NewHTTPClient(opts ...Option) *HTTPClient {
	options := Options{
		APIKey:     os.Getenv("QDRANT_API_KEY"),
		BaseURL:    DefaultBaseURL,
		HTTPClient: &http.Client{},
	}

	for _, apply := range opts {
		apply(&options)
	}

	return &HTTPClient{
		opts: options,
	}
}

// WithAPIKey sets the API key.
func WithAPIKey(apiKey string) Option {
	return func(o *Options) {
		o.APIKey = apiKey
	}
}

// WithBaseURL sets the API base URL.
func WithBaseURL(baseURL string) Option {
	return func(o *Options) {
		o.BaseURL = baseURL
	}
}

// WithHTTPClient sets the HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(o *Options) {
		o.HTTPClient = httpClient
	}
}

// UpdateAliases action (Create, Remove, Switch)
func (c *HTTPClient) UpdateAliases(ctx context.Context, actions []Action) error {
	u, err := url.Parse(c.opts.BaseURL + "/collections/aliases")
	if err != nil {
		return err
	}

	var body = &bytes.Buffer{}
	enc := json.NewEncoder(body)

	aliasReq := &AliasActionReq{
		Actions: actions,
	}

	if err := enc.Encode(aliasReq); err != nil {
		return err
	}

	options := []ReqOption{}
	if c.opts.APIKey != "" {
		options = append(options, WithBearer(c.opts.APIKey))
	}

	req, err := NewRequest(ctx, http.MethodPost, u.String(), body, options...)
	if err != nil {
		return err
	}
	resp, err := Do(c.opts.HTTPClient, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// AliasList action.
func (c *HTTPClient) AliasList(ctx context.Context, collName string) (*AliasesResp, error) {
	u, err := url.Parse(c.opts.BaseURL + "/collections/" + collName + "/aliases")
	if err != nil {
		return nil, err
	}

	options := []ReqOption{}
	if c.opts.APIKey != "" {
		options = append(options, WithBearer(c.opts.APIKey))
	}

	req, err := NewRequest(ctx, http.MethodGet, u.String(), nil, options...)
	if err != nil {
		return nil, err
	}
	resp, err := Do(c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	aliases := new(AliasesResp)
	if err := json.NewDecoder(resp.Body).Decode(aliases); err != nil {
		return nil, err
	}

	return aliases, nil
}

// NewRequest creates a new HTTP request from the provided parameters  and returns it.
// If the passed in context is nil, it creates a new background context.
// If the provided body is nil, it gets initialized to bytes.Reader.
// By default the following headers are set:
// * Accept: application/json; charset=utf-8
// If no Content-Type has been set via options it defaults to application/json.
func NewRequest(ctx context.Context, method, url string, body io.Reader, opts ...ReqOption) (*http.Request, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if body == nil {
		body = &bytes.Reader{}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	for _, setOption := range opts {
		setOption(req)
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	// if no content-type is specified we default to json
	if ct := req.Header.Get("Content-Type"); len(ct) == 0 {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	return req, nil
}

// Do sends the HTTP request req using the client and returns the response.
func Do(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest {
		return resp, nil
	}
	defer resp.Body.Close()

	var apiErr APIError
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		return nil, err
	}

	return nil, apiErr
}

// ReqOption is http requestion functional option.
type ReqOption func(*http.Request)

// WithBearer sets the Authorization header to the provided Bearer token.
func WithBearer(token string) ReqOption {
	return func(req *http.Request) {
		if req.Header == nil {
			req.Header = make(http.Header)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

// Action to take.
type Action map[string]map[string]string

// AliasActionReq is used to update collection aliases.
type AliasActionReq struct {
	// Actions to take on alias.
	Actions []Action `json:"actions"`
}

// Alias
type Alias struct {
	Name       string `json:"alias_name"`
	Collection string `json:"collection_name"`
}

// AliasesResp
type AliasesResp struct {
	Time   float64 `json:"time"`
	Status string  `json:"status"`
	Result struct {
		Aliases []Alias `json:"aliases"`
	} `json:"result"`
}

// APIError encodes qdrant API error.
type APIError struct {
	Time   float64 `json:"time"`
	Status struct {
		Error string `json:"error"`
	} `json:"status"`
	Result struct {
		OperationID int    `json:"operation_id"`
		Status      string `json:"status"`
	} `json:"result"`
}

// Error implements error interface.
func (e APIError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "unknown error"
	}
	return string(b)
}
