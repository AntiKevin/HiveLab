package core

import "testing"

func TestNewHttpMockModelValidatesURLAndPort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		rawURL string
		port   int
	}{
		{name: "empty URL", rawURL: "", port: 8080},
		{name: "relative URL", rawURL: "localhost", port: 8080},
		{name: "zero port", rawURL: "http://localhost", port: 0},
		{name: "port too high", rawURL: "http://localhost", port: 65536},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if _, err := NewHttpMockModel(tt.rawURL, tt.port); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestNewHttpMockModelAcceptsValidConfiguration(t *testing.T) {
	t.Parallel()

	mock, err := NewHttpMockModel("http://localhost", 8080, Response{
		ID:         "success",
		StatusCode: 201,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       []byte(`{"ok":true}`),
	})
	if err != nil {
		t.Fatalf("expected valid mock: %v", err)
	}

	response, found := mock.ResponseByID("success")
	if !found {
		t.Fatal("expected response to be registered")
	}

	if response.EffectiveStatusCode() != 201 {
		t.Fatalf("expected status 201, got %d", response.EffectiveStatusCode())
	}

	if mock.BaseURL() != "http://localhost:8080" {
		t.Fatalf("expected base URL with port, got %q", mock.BaseURL())
	}

	if mock.EffectiveMethod() != "GET" {
		t.Fatalf("expected default method GET, got %q", mock.EffectiveMethod())
	}
}

func TestAddResponseRejectsDuplicateID(t *testing.T) {
	t.Parallel()

	mock, err := NewHttpMockModel("http://localhost", 8080, Response{ID: "same"})
	if err != nil {
		t.Fatalf("expected valid mock: %v", err)
	}

	err = mock.AddResponse(Response{ID: "same"})
	if err == nil {
		t.Fatal("expected duplicate response error")
	}
}

func TestResponseValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		response Response
		wantErr  bool
	}{
		{name: "empty ID", response: Response{}, wantErr: true},
		{name: "default status", response: Response{ID: "ok"}, wantErr: false},
		{name: "lower than HTTP range", response: Response{ID: "bad", StatusCode: 99}, wantErr: true},
		{name: "higher than HTTP range", response: Response{ID: "bad", StatusCode: 600}, wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.response.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpMockModelMethodValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		method  string
		wantErr bool
	}{
		{name: "default method", method: "", wantErr: false},
		{name: "lowercase method", method: "post", wantErr: false},
		{name: "invalid method", method: "FETCH", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := HttpMockModel{
				URL:    "http://localhost",
				Port:   8080,
				Method: tt.method,
			}

			err := mock.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResponseByIDReturnsClone(t *testing.T) {
	t.Parallel()

	mock, err := NewHttpMockModel("http://localhost", 8080, Response{
		ID:      "success",
		Headers: map[string]string{"X-Test": "original"},
		Body:    []byte("original"),
	})
	if err != nil {
		t.Fatalf("expected valid mock: %v", err)
	}

	response, found := mock.ResponseByID("success")
	if !found {
		t.Fatal("expected response to be found")
	}

	response.Headers["X-Test"] = "changed"
	response.Body[0] = 'X'

	unchanged, found := mock.ResponseByID("success")
	if !found {
		t.Fatal("expected response to still be found")
	}

	if unchanged.Headers["X-Test"] != "original" {
		t.Fatalf("expected header clone, got %q", unchanged.Headers["X-Test"])
	}

	if string(unchanged.Body) != "original" {
		t.Fatalf("expected body clone, got %q", string(unchanged.Body))
	}
}

func TestAddPayloadRequiresKnownResponse(t *testing.T) {
	t.Parallel()

	mock, err := NewHttpMockModel("http://localhost", 8080, Response{ID: "success"})
	if err != nil {
		t.Fatalf("expected valid mock: %v", err)
	}

	if err := mock.AddPayload(Payload{ResponseID: "missing"}); err == nil {
		t.Fatal("expected unknown response error")
	}

	if err := mock.AddPayload(Payload{ResponseID: "success", Body: []byte("request")}); err != nil {
		t.Fatalf("expected payload to be valid: %v", err)
	}
}

func TestPayloadValidation(t *testing.T) {
	t.Parallel()

	if err := (Payload{}).Validate(); err == nil {
		t.Fatal("expected empty response ID to be invalid")
	}

	if err := (Payload{ResponseID: "success"}).Validate(); err != nil {
		t.Fatalf("expected payload to be valid: %v", err)
	}
}
