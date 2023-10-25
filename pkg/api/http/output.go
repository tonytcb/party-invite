package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"

	"github.com/tonytcb/party-invite/pkg/domain"
)

type customer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func customersToJSONOutput(input domain.Customers) ([]byte, error) {
	var customers = make([]customer, 0)

	for _, v := range input {
		customers = append(customers, customer{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	bytes, err := json.Marshal(customers)
	if err != nil {
		return nil, errors.Wrap(err, "error to encode customers output")
	}

	return bytes, err
}

type httpError struct {
	err     error
	details string
	code    int
	Error   string `json:"error"`
}

func newHTTPError(err error, details string, code int) *httpError {
	return &httpError{err: err, details: details, code: code}
}

func (e httpError) json(w http.ResponseWriter) {
	errorContent := e.details
	if e.err != nil {
		errorContent = fmt.Sprintf("%s: %s", e.details, e.err.Error())
	}

	content := struct {
		Error string `json:"error"`
	}{
		Error: errorContent,
	}

	w.WriteHeader(e.code)

	output, err := json.Marshal(content)
	if err != nil {
		return
	}

	_, _ = w.Write(output)
}

func (e httpError) empty(w http.ResponseWriter) {
	w.WriteHeader(e.code)
	_, _ = w.Write([]byte("{}"))
}

func errToStatusCode(err error) int {
	var invalidArgumentErr *domain.ErrInvalidArgument

	switch {
	case errors.As(err, &invalidArgumentErr):
		return http.StatusUnprocessableEntity

	case errors.Is(err, context.Canceled):
		return http.StatusGatewayTimeout

	default:
		return http.StatusInternalServerError
	}
}
