package http

import (
	"bytes"
	"go.uber.org/mock/gomock"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/domain"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/infrastructure/config"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/infrastructure/logger"
)

func TestFilterCustomersHandler_Handle(t *testing.T) {
	t.Parallel()

	defaultConfig := &config.Config{
		BaseLocation:   "dublin",
		LocationNearTo: 100,
	}

	saoPaulo, err := domain.NewCoordinate("-23.533773", "-46.625290")
	if err != nil {
		t.Fatal("failed to build coordinate")
	}
	montreal, err := domain.NewCoordinate("45.508888", "-73.561668")
	if err != nil {
		t.Fatal("failed to build coordinate")
	}

	var (
		customer1 = domain.NewCustomer(1, "User name 1", domain.DublinLocation)
		customer2 = domain.NewCustomer(2, "User name 2", saoPaulo)
		customer3 = domain.NewCustomer(3, "User name 3", montreal)

		customersList1 = []domain.Customer{customer1, customer2, customer3}
	)

	postRequestWithValidFile, err := newRequestWithFile(http.MethodPost, "localhost:8080", "file", "customers.txt")
	if err != nil {
		t.Fatal("failed to create valid request")
	}

	putRequest, _ := newRequestWithFile(http.MethodPut, "localhost:8080", "file", "customers.txt")
	postRequestWithInvalidRequestParamName, _ := newRequestWithFile(http.MethodPost, "localhost:8080", "another_name", "customers.txt")
	postRequestWithInvalidFile, _ := newRequestWithFile(http.MethodPost, "localhost:8080", "file", "invalid-ext.sql")

	var log = logger.NewLogger(&bytes.Buffer{})

	type fields struct {
		parser func(*testing.T, *gomock.Controller) CustomersFileParser
		filter func(*testing.T, *gomock.Controller) FilterCustomersUsecase
	}
	type args struct {
		responseWriter http.ResponseWriter
		request        *http.Request
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantStatusCode   int
		wantResponseBody string
	}{
		{
			name: "should handle successfully a file containing valid customers",
			fields: fields{
				parser: func(t *testing.T, ctrl *gomock.Controller) CustomersFileParser {
					parser := NewMockCustomersFileParser(ctrl)
					parser.EXPECT().
						Parse(gomock.Any(), gomock.Any()).
						Return(customersList1, nil).
						Times(1)

					return parser
				},
				filter: func(t *testing.T, ctrl *gomock.Controller) FilterCustomersUsecase {
					filter := NewMockFilterCustomersUsecase(ctrl)

					customers := []domain.Customer{customer1, customer2}

					filter.EXPECT().
						ByNearLocation(gomock.Any(), customersList1, domain.DublinLocation, decimal.NewFromInt32(100), domain.OrderByCustomerID).
						Return(customers, nil).
						Times(1)

					return filter
				},
			},
			args: args{
				responseWriter: httptest.NewRecorder(),
				request:        postRequestWithValidFile,
			},
			wantStatusCode:   http.StatusOK,
			wantResponseBody: `[{"id":1,"name":"User name 1"},{"id":2,"name":"User name 2"}]`,
		},

		{
			name: "should error on http method not allowed error",
			fields: fields{
				parser: func(t *testing.T, ctrl *gomock.Controller) CustomersFileParser {
					return NewMockCustomersFileParser(ctrl)
				},
				filter: func(t *testing.T, ctrl *gomock.Controller) FilterCustomersUsecase {
					return NewMockFilterCustomersUsecase(ctrl)
				},
			},
			args: args{
				responseWriter: httptest.NewRecorder(),
				request:        putRequest,
			},
			wantStatusCode:   http.StatusMethodNotAllowed,
			wantResponseBody: `{}`,
		},
		{
			name: "should error on invalid request field name",
			fields: fields{
				parser: func(t *testing.T, ctrl *gomock.Controller) CustomersFileParser {
					return NewMockCustomersFileParser(ctrl)
				},
				filter: func(t *testing.T, ctrl *gomock.Controller) FilterCustomersUsecase {
					return NewMockFilterCustomersUsecase(ctrl)
				},
			},
			args: args{
				responseWriter: httptest.NewRecorder(),
				request:        postRequestWithInvalidRequestParamName,
			},
			wantStatusCode:   http.StatusBadRequest,
			wantResponseBody: `{"error":"error to read uploaded file: http: no such file"}`,
		},
		{
			name: "should error on invalid file extension",
			fields: fields{
				parser: func(t *testing.T, ctrl *gomock.Controller) CustomersFileParser {
					return NewMockCustomersFileParser(ctrl)
				},
				filter: func(t *testing.T, ctrl *gomock.Controller) FilterCustomersUsecase {
					return NewMockFilterCustomersUsecase(ctrl)
				},
			},
			args: args{
				responseWriter: httptest.NewRecorder(),
				request:        postRequestWithInvalidFile,
			},
			wantStatusCode:   http.StatusBadRequest,
			wantResponseBody: `{"error":"invalid '.sql' file extension"}`,
		},
		{
			name: "should error on parser due to some domain validation",
			fields: fields{
				parser: func(t *testing.T, ctrl *gomock.Controller) CustomersFileParser {
					parser := NewMockCustomersFileParser(ctrl)
					parser.EXPECT().
						Parse(gomock.Any(), gomock.Any()).
						Return(nil, domain.NewErrInvalidArgument("root cause", "failed some domain validation"))

					return parser
				},
				filter: func(t *testing.T, ctrl *gomock.Controller) FilterCustomersUsecase {
					return NewMockFilterCustomersUsecase(ctrl)
				},
			},
			args: args{
				responseWriter: httptest.NewRecorder(),
				request:        postRequestWithValidFile,
			},
			wantStatusCode:   http.StatusUnprocessableEntity,
			wantResponseBody: `{"error":"error to parse input file: failed some domain validation: root cause"}`,
		},
		{
			name: "should error on filter customers usecase",
			fields: fields{
				parser: func(t *testing.T, ctrl *gomock.Controller) CustomersFileParser {
					parser := NewMockCustomersFileParser(ctrl)
					parser.EXPECT().
						Parse(gomock.Any(), gomock.Any()).
						Return(customersList1, nil).
						Times(1)

					return parser
				},
				filter: func(t *testing.T, ctrl *gomock.Controller) FilterCustomersUsecase {
					filter := NewMockFilterCustomersUsecase(ctrl)
					filter.EXPECT().
						ByNearLocation(gomock.Any(), customersList1, domain.DublinLocation, decimal.NewFromInt32(100), domain.OrderByCustomerID).
						Return(nil, errors.New("some error on calculation")).
						Times(1)

					return filter
				},
			},
			args: args{
				responseWriter: httptest.NewRecorder(),
				request:        postRequestWithValidFile,
			},
			wantStatusCode:   http.StatusInternalServerError,
			wantResponseBody: `{"error":"error to filter customers by location: some error on calculation"}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			h := &FilterCustomersHandler{
				log:    log,
				cfg:    defaultConfig,
				parser: tt.fields.parser(t, mockCtrl),
				filter: tt.fields.filter(t, mockCtrl),
			}
			h.Handle(tt.args.responseWriter, tt.args.request)

			if v, ok := tt.args.responseWriter.(*httptest.ResponseRecorder); ok {
				httpResponse := v.Result()

				bytesResponse, err := io.ReadAll(httpResponse.Body)
				if err != nil {
					t.Fatal("failed to read buffer response body")
				}

				assert.Equal(t, tt.wantStatusCode, httpResponse.StatusCode, "HTTP Status Code does not match")
				assert.Equal(t, tt.wantResponseBody, string(bytesResponse), "HTTP Response Body does not match")
			}
		})
	}
}

func newRequestWithFile(method string, endpoint string, fieldName string, fileName string) (*http.Request, error) {
	currentDir, _ := os.Getwd()
	fileDir := currentDir + "/../../../Data"
	filePath := path.Join(fileDir, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "error to open file")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filepath.Base(file.Name()))
	if err != nil {
		return nil, errors.Wrap(err, "error to create form file")
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, errors.Wrap(err, "error to read multipart file")
	}
	writer.Close()

	r, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, errors.Wrap(err, "error to build request")
	}

	r.Header.Add("Content-Type", writer.FormDataContentType())

	return r, nil
}
