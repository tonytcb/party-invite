package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/tonytcb/party-invite/pkg/domain"
	"github.com/tonytcb/party-invite/pkg/infrastructure/config"
	"github.com/tonytcb/party-invite/pkg/infrastructure/logger"
)

const (
	timeoutDefault        = 30 * time.Second
	txtFileExtension      = ".txt"
	maximumFileUploadSize = 10 << 20 // 10mb
)

//go:generate mockgen -source=filtercustomershandler.go -destination=mock_filtercustomers_test.go -package=http CustomersFileParser,FilterCustomersUsecase,FilterCustomersCache

type CustomersFileParser interface {
	Parse(context.Context, io.Reader) (domain.Customers, error)
}

type FilterCustomersUsecase interface {
	ByNearLocation(
		ctx context.Context,
		customers domain.Customers,
		baseLocation *domain.Coordinate,
		nearDistanceFilter decimal.Decimal,
		orderBy domain.OrderBy,
	) (domain.Customers, error)
}

type FilterCustomersCache interface {
	Get(context.Context, []byte) ([]byte, error)
	Save(context.Context, []byte, []byte) error
}

type FilterCustomersHandler struct {
	log logger.Logger
	cfg *config.Config

	parser CustomersFileParser
	filter FilterCustomersUsecase
	cache  FilterCustomersCache
}

func NewFilterCustomersHandler(
	log logger.Logger,
	cfg *config.Config,
	parser CustomersFileParser,
	filter FilterCustomersUsecase,
	cache FilterCustomersCache,
) *FilterCustomersHandler {
	return &FilterCustomersHandler{log: log, cfg: cfg, parser: parser, filter: filter, cache: cache}
}

// Handle filters a list of customer given the input file.
func (h *FilterCustomersHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var (
		correlationID = uuid.NewString()
		log           = h.log.WithCorrelationID(correlationID)
		ctx           = context.WithValue(r.Context(), config.CorrelationIDKeyName, correlationID)
	)

	ctx, cancel := context.WithTimeout(ctx, timeoutDefault)
	defer cancel()

	if err := r.ParseMultipartForm(maximumFileUploadSize); err != nil {
		newHTTPError(err, "error to set max upload file", http.StatusInternalServerError).json(w)
		return
	}

	if r.Method != http.MethodPost {
		newHTTPError(nil, "", http.StatusMethodNotAllowed).empty(w)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		newHTTPError(err, "error to read uploaded file", http.StatusBadRequest).json(w)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Errorf("Error to close uploaded file, err=%v", err)
		}
	}()

	log.Infof("Filtering customers, filename=%s filesize=%d", header.Filename, header.Size)

	if ext := filepath.Ext(header.Filename); ext != txtFileExtension {
		newHTTPError(nil, "invalid '"+ext+"' file extension", http.StatusBadRequest).json(w)
		return
	}

	// to read the file twice, we need to duplicate it since it's a buffer
	var tempBuf = &bytes.Buffer{}
	var tempFile = io.TeeReader(file, tempBuf)

	fileContents, _ := io.ReadAll(tempBuf)
	//fileReader := io.NopCloser(bytes.NewReader(fileContents))

	cachedResponse, err := h.cache.Get(ctx, fileContents)
	if err != nil {
		newHTTPError(err, "error to load cache", errToStatusCode(err)).json(w)
		return
	}

	if cachedResponse != nil {
		w.Write(cachedResponse) //nolint:errcheck
		return
	}

	customers, err := h.parser.Parse(ctx, tempFile)
	if err != nil {
		newHTTPError(err, "error to parse input file", errToStatusCode(err)).json(w)
		return
	}

	filteredCustomers, err := h.filter.ByNearLocation(
		ctx,
		customers,
		h.cfg.GetBaseLocation(),
		decimal.NewFromInt32(h.cfg.LocationNearTo),
		domain.OrderByCustomerID,
	)
	if err != nil {
		newHTTPError(err, "error to filter customers by location", errToStatusCode(err)).json(w)
		return
	}

	response, err := customersToJSONOutput(filteredCustomers)
	if err != nil {
		newHTTPError(err, "error to build response output", http.StatusServiceUnavailable).json(w)
		return
	}

	log.Infof("Filtered customers length response: input=%d output=%d", len(customers), len(filteredCustomers))

	if _, err = w.Write(response); err != nil {
		newHTTPError(err, "", http.StatusInternalServerError).empty(w)
	}

	if err = h.cache.Save(ctx, fileContents, response); err != nil {
		log.Errorf("Error to store response on cache: %v", err)
	}
}
