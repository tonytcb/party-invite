package http

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/tonytcb/party-invite/pkg/infrastructure/config"
	"github.com/tonytcb/party-invite/pkg/infrastructure/customerfile"
	"github.com/tonytcb/party-invite/pkg/infrastructure/logger"
	"github.com/tonytcb/party-invite/pkg/usecase"
)

// TestFilterCustomersAPI start up the http handler with real dependencies to assert http response.
func TestFilterCustomersAPI(t *testing.T) {
	t.Parallel()

	var (
		buf = bytes.NewBuffer([]byte(""))
		log = logger.NewLogger(buf)
	)

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("error to load configuration: %v", err)
	}

	var filterCustomersHandler = NewFilterCustomersHandler(
		log,
		cfg,
		customerfile.NewCustomersFileParser(),
		usecase.NewFilterCustomers(log),
	)

	w := httptest.NewRecorder()

	postRequestWithValidFile, err := newRequestWithFile(http.MethodPost, "localhost:8080/filter-customers", "file", "customers.txt")
	if err != nil {
		t.Fatal("failed to create valid request")
	}

	filterCustomersHandler.Handle(w, postRequestWithValidFile)

	// assert

	const expectedStatusCode = 200
	const expectedBody = `[{"id":4,"name":"Ian Kehoe"},{"id":5,"name":"Nora Dempsey"},{"id":6,"name":"Theresa Enright"},{"id":8,"name":"Eoin Ahearn"},{"id":11,"name":"Richard Finnegan"},{"id":12,"name":"Christina McArdle"},{"id":13,"name":"Olive Ahearn"},{"id":15,"name":"Michael Ahearn"},{"id":17,"name":"Patricia Cahill"},{"id":23,"name":"Eoin Gallagher"},{"id":24,"name":"Rose Enright"},{"id":26,"name":"Stephen McArdle"},{"id":29,"name":"Oliver Ahearn"},{"id":30,"name":"Nick Enright"},{"id":31,"name":"Alan Behan"},{"id":39,"name":"Lisa Ahearn"}]`

	httpResponse := w.Result()
	defer httpResponse.Body.Close()

	bytesResponse, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		t.Fatal("failed to read buffer response body")
	}

	assert.Equal(t, expectedStatusCode, httpResponse.StatusCode, "HTTP Status Code does not match")
	assert.Equal(t, expectedBody, string(bytesResponse), "HTTP Response Body does not match")
}

func loadConfig() (*config.Config, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, errors.Errorf("error to load current directory: %v", err)
	}

	rootDir := currentDir + "/../../../"

	cfg, err := config.Load(rootDir)
	if err != nil {
		return nil, errors.Errorf("error to load current directory: %v", err)
	}

	if err = cfg.IsValid(); err != nil {
		return nil, errors.Errorf("invalid configuration: %v", err)
	}

	return cfg, nil
}
