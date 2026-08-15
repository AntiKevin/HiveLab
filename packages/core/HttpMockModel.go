package core

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const defaultMockStatusCode = 200

// Response describes the HTTP response that a mock can return.
type Response struct {
	ID         string            `json:"id"`
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       []byte            `json:"body,omitempty"`
}

// Payload points a request payload to the response that should be used.
type Payload struct {
	ResponseID string `json:"responseId"`
	Body       []byte `json:"body,omitempty"`
}

// HttpMockModel describes one HTTP mock service.
type HttpMockModel struct {
	URL       string            `json:"url"`
	Port      int               `json:"port"`
	Headers   map[string]string `json:"headers,omitempty"`
	Method    string            `json:"method,omitempty"`
	Payloads  []Payload         `json:"payloads"`
	Responses []Response        `json:"responses,omitempty"`
}

// NewHttpMockModel creates a validated mock service definition.
func NewHttpMockModel(rawURL string, port int, responses ...Response) (*HttpMockModel, error) {
	mock := &HttpMockModel{
		URL:       strings.TrimSpace(rawURL),
		Port:      port,
		Payloads:  make([]Payload, 0),
		Responses: make([]Response, 0, len(responses)),
	}

	for _, response := range responses {
		if err := mock.AddResponse(response); err != nil {
			return nil, err
		}
	}

	if err := mock.Validate(); err != nil {
		return nil, err
	}

	return mock, nil
}

// AddResponse registers a response, keeping the model immutable from callers.
func (m *HttpMockModel) AddResponse(response Response) error {
	if err := response.Validate(); err != nil {
		return err
	}

	if _, found := m.ResponseByID(response.ID); found {
		return fmt.Errorf("response %q already exists", response.ID)
	}

	m.Responses = append(m.Responses, response.Clone())
	return nil
}

// AddPayload registers a payload, keeping the model immutable from callers.
func (m *HttpMockModel) AddPayload(payload Payload) error {
	if err := payload.Validate(); err != nil {
		return err
	}

	if _, found := m.ResponseByID(payload.ResponseID); !found {
		return fmt.Errorf("payload references unknown response %q", payload.ResponseID)
	}

	m.Payloads = append(m.Payloads, payload.Clone())
	return nil
}

// ResponseByID finds a response and returns a copy that callers can mutate.
func (m HttpMockModel) ResponseByID(id string) (Response, bool) {
	for _, response := range m.Responses {
		if response.ID == id {
			return response.Clone(), true
		}
	}

	return Response{}, false
}

// BaseURL returns the mock service URL with the configured port applied.
func (m HttpMockModel) BaseURL() string {
	parsedURL, err := url.Parse(m.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return ""
	}

	parsedURL.Host = net.JoinHostPort(parsedURL.Hostname(), strconv.Itoa(m.Port))
	return strings.TrimRight(parsedURL.String(), "/")
}

// EffectiveMethod returns the configured HTTP method or the default method.
func (m HttpMockModel) EffectiveMethod() string {
	method := strings.TrimSpace(m.Method)
	if method == "" {
		return http.MethodGet
	}

	return strings.ToUpper(method)
}

// Validate checks if the mock service definition is ready to run.
func (m HttpMockModel) Validate() error {
	if strings.TrimSpace(m.URL) == "" {
		return errors.New("mock URL is required")
	}

	parsedURL, err := url.ParseRequestURI(m.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid mock URL %q", m.URL)
	}

	if m.Port < 1 || m.Port > 65535 {
		return fmt.Errorf("mock port must be between 1 and 65535: %d", m.Port)
	}

	if !isSupportedHTTPMethod(m.EffectiveMethod()) {
		return fmt.Errorf("mock has invalid HTTP method: %q", m.Method)
	}

	seenResponses := make(map[string]struct{}, len(m.Responses))
	for _, response := range m.Responses {
		if err := response.Validate(); err != nil {
			return err
		}

		if _, found := seenResponses[response.ID]; found {
			return fmt.Errorf("response %q already exists", response.ID)
		}
		seenResponses[response.ID] = struct{}{}
	}

	for _, payload := range m.Payloads {
		if err := payload.Validate(); err != nil {
			return err
		}

		if _, found := seenResponses[payload.ResponseID]; !found {
			return fmt.Errorf("payload references unknown response %q", payload.ResponseID)
		}
	}

	return nil
}

// Clone returns a deep copy of the response body and headers.
func (r Response) Clone() Response {
	clone := Response{
		ID:         r.ID,
		StatusCode: r.StatusCode,
		Body:       append([]byte(nil), r.Body...),
	}

	if len(r.Headers) > 0 {
		clone.Headers = make(map[string]string, len(r.Headers))
		for key, value := range r.Headers {
			clone.Headers[key] = value
		}
	}

	return clone
}

// Validate checks if the response can be served by an HTTP mock.
func (r Response) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("response ID is required")
	}

	if r.StatusCode == 0 {
		return nil
	}

	if r.StatusCode < 100 || r.StatusCode > 599 {
		return fmt.Errorf("response %q has invalid HTTP status code: %d", r.ID, r.StatusCode)
	}

	return nil
}

// EffectiveStatusCode returns the configured status code or the default code.
func (r Response) EffectiveStatusCode() int {
	if r.StatusCode == 0 {
		return defaultMockStatusCode
	}

	return r.StatusCode
}

// Clone returns a copy of the payload body.
func (p Payload) Clone() Payload {
	return Payload{
		ResponseID: p.ResponseID,
		Body:       append([]byte(nil), p.Body...),
	}
}

// Validate checks if the payload references a response.
func (p Payload) Validate() error {
	if strings.TrimSpace(p.ResponseID) == "" {
		return errors.New("payload response ID is required")
	}

	return nil
}

func isSupportedHTTPMethod(method string) bool {
	switch method {
	case http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace:
		return true
	default:
		return false
	}
}
