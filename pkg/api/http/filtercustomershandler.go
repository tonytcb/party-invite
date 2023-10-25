package http

import (
	"context"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/infrastructure/config"
	"github.com/google/uuid"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/shopspring/decimal"

	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/domain"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/infrastructure/logger"
)

const (
	timeoutDefault        = 30 * time.Second
	txtFileExtension      = ".txt"
	maximumFileUploadSize = 10 << 20 // 10mb
)

//go:generate mockgen -source=filtercustomershandler.go -destination=mock_filtercustomers_test.go -package=http CustomersFileParser,FilterCustomersUsecase

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

type FilterCustomersHandler struct {
	log    logger.Logger
	cfg    *config.Config
	parser CustomersFileParser
	filter FilterCustomersUsecase
}

func NewFilterCustomersHandler(
	log logger.Logger,
	cfg *config.Config,
	parser CustomersFileParser,
	filter FilterCustomersUsecase,
) *FilterCustomersHandler {
	return &FilterCustomersHandler{log: log, cfg: cfg, parser: parser, filter: filter}
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

	customers, err := h.parser.Parse(ctx, file)
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
}
