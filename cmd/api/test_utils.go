package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kuluruvineeth/social-go/internal/auth"
	"github.com/kuluruvineeth/social-go/internal/store"
	"github.com/kuluruvineeth/social-go/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()

	mockStore := store.NewMockStorage()
	mockCacheStore := cache.NewMockStorage()

	testAuth := &auth.TestAuthenticator{}

	return &application{
		logger:        logger,
		store:         mockStore,
		cache:         mockCacheStore,
		authenticator: testAuth,
		config:        cfg,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d, got %d", expected, actual)
	}
}
