package testutil

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupTestDB(t *testing.T) {
	db := SetupTestDB(t)
	assert.NotNil(t, db)
	err := db.Ping()
	assert.NoError(t, err)
	db.Close()
}

func TestMakeRequest_Get(t *testing.T) {
	req := MakeRequest(http.MethodGet, "/test/path", nil)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "/test/path", req.URL.Path)
}

func TestMakeRequest_PostWithBody(t *testing.T) {
	body := bytes.NewBufferString(`{"key":"value"}`)
	req := MakeRequest(http.MethodPost, "/api/test", body)

	assert.Equal(t, http.MethodPost, req.Method)
	assert.Equal(t, "/api/test", req.URL.Path)
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
}

func TestMakeRequest_EmptyBody(t *testing.T) {
	req := MakeRequest(http.MethodDelete, "/api/resource/1", nil)

	assert.Equal(t, http.MethodDelete, req.Method)
	assert.Equal(t, "/api/resource/1", req.URL.Path)
}

func TestExecuteRequest_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	req := MakeRequest(http.MethodGet, "/test", nil)
	recorder := ExecuteRequest(handler, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"status":"ok"}`, recorder.Body.String())
}

func TestExecuteRequest_PostWithData(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"created":true}`))
	})

	body := bytes.NewBufferString(`{"name":"test"}`)
	req := MakeRequest(http.MethodPost, "/api/create", body)
	recorder := ExecuteRequest(handler, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.JSONEq(t, `{"created":true}`, recorder.Body.String())
}

func TestExecuteRequest_ErrorResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid request"}`))
	})

	req := MakeRequest(http.MethodPost, "/api/error", nil)
	recorder := ExecuteRequest(handler, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
