// Package testutils provides testing utilities for the Go Starter application.
package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// HTTPTestHelper provides utilities for HTTP testing
type HTTPTestHelper struct {
	T *testing.T
}

// NewHTTPTestHelper creates a new HTTP test helper
func NewHTTPTestHelper(t *testing.T) *HTTPTestHelper {
	return &HTTPTestHelper{T: t}
}

// MakeRequest creates an HTTP request for testing
func (h *HTTPTestHelper) MakeRequest(
	handler http.Handler,
	method, path string,
	body interface{},
) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(h.T, err)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}

// AssertJSONResponse asserts that the response contains valid JSON
func (h *HTTPTestHelper) AssertJSONResponse(rr *httptest.ResponseRecorder) {
	require.Equal(h.T, "application/json", rr.Header().Get("Content-Type"))

	var result interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	require.NoError(h.T, err)
}

// AssertStatus asserts the HTTP status code
func (h *HTTPTestHelper) AssertStatus(rr *httptest.ResponseRecorder, expectedStatus int) {
	require.Equal(h.T, expectedStatus, rr.Code)
}

// AssertJSONBody asserts the JSON response body
func (h *HTTPTestHelper) AssertJSONBody(rr *httptest.ResponseRecorder, expectedBody interface{}) {
	h.AssertJSONResponse(rr)

	var actual interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &actual)
	require.NoError(h.T, err)

	require.Equal(h.T, expectedBody, actual)
}

// ConfigHelper provides utilities for configuration testing
type ConfigHelper struct {
	T *testing.T
}

// NewConfigHelper creates a new config test helper
func NewConfigHelper(t *testing.T) *ConfigHelper {
	return &ConfigHelper{T: t}
}

// AssertConfigField asserts that a configuration field has the expected value
func (c *ConfigHelper) AssertConfigField(actual, expected interface{}, fieldName string) {
	require.Equal(c.T, expected, actual, "Field %s should have value %v", fieldName, expected)
}

// LoggerHelper provides utilities for logger testing
type LoggerHelper struct {
	T *testing.T
}

// NewLoggerHelper creates a new logger test helper
func NewLoggerHelper(t *testing.T) *LoggerHelper {
	return &LoggerHelper{T: t}
}

// AssertLogContains asserts that a log message contains expected text
func (l *LoggerHelper) AssertLogContains(logOutput, expectedText string) {
	require.Contains(l.T, logOutput, expectedText, "Log output should contain %s", expectedText)
}

// MockServer provides a mock HTTP server for testing
type MockServer struct {
	Server *httptest.Server
}

// NewMockServer creates a new mock server
func NewMockServer(handler http.Handler) *MockServer {
	return &MockServer{
		Server: httptest.NewServer(handler),
	}
}

// Close closes the mock server
func (ms *MockServer) Close() {
	ms.Server.Close()
}

// URL returns the mock server URL
func (ms *MockServer) URL() string {
	return ms.Server.URL
}

// TestDatabase provides a simple in-memory database for testing
type TestDatabase struct {
	// Add database fields here when implementing actual database tests
	// This is a placeholder for future database testing utilities
}

// NewTestDatabase creates a new test database
func NewTestDatabase() *TestDatabase {
	return &TestDatabase{}
}

// Cleanup cleans up the test database
func (td *TestDatabase) Cleanup() {
	// Add cleanup logic here when implementing actual database tests
}

// TestUser represents a test user for testing purposes
type TestUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type TestResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CreateTestResponse creates a test response for testing purposes
func CreateTestResponse(status, message string, data interface{}) TestResponse {
	return TestResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// CreateTestUser creates a test user for testing purposes
func CreateTestUser(id, email, name string) TestUser {
	return TestUser{
		ID:    id,
		Email: email,
		Name:  name,
	}
}
