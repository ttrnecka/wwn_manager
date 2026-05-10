package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/mapper"
	"github.com/ttrnecka/wwn_identity/webapi/internal/service"
	"github.com/ttrnecka/wwn_identity/webapi/ita"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
)

var wwnRegex = regexp.MustCompile(`^([0-9A-Fa-f]{2}:){7}[0-9A-Fa-f]{2}$`)

type FCWWNEntryHandler struct {
	service     service.FCWWNEntryService
	ruleService service.RuleService
	logger      *zerolog.Logger
}

type entryLineRecord struct {
	customer, wwn, zone, alias, loadedHostname string
	isCsvLoad                                  bool
	wwnSet                                     int
}

func NewFCWWNEntryHandler(s service.FCWWNEntryService, r service.RuleService, logger *zerolog.Logger) *FCWWNEntryHandler {
	l := logger.With().Str("component", "FCWWNEntryHandler").Logger()
	return &FCWWNEntryHandler{s, r, &l}
}

func (h *FCWWNEntryHandler) FCWWNEntries(c echo.Context) error {
	mode := c.QueryParam("softdeleted")
	if mode != "" {
		return h.FCWWNEntriesWithSoftDeleted(c)
	}
	customer := c.Param("name")

	items, err := h.service.Find(c.Request().Context(), service.Filter{"customer": customer}, service.SortOption{"wwn": "asc"})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to find entries", err)
	}

	itemDTO := make([]dto.FCWWNEntryDTO, 0)

	for _, item := range items {
		itemDTO = append(itemDTO, mapper.ToFCWWNEntryDTO(item))
	}

	return c.JSON(http.StatusOK, itemDTO)
}

func (h *FCWWNEntryHandler) FCWWNEntriesWithSoftDeleted(c echo.Context) error {
	customer := c.Param("name")

	items, err := h.service.FindWithSoftDeleted(c.Request().Context(), service.Filter{"customer": customer}, service.SortOption{"wwn": "asc"})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to find softdeleted entries", err)
	}

	itemDTO := make([]dto.FCWWNEntryDTO, 0)

	for _, item := range items {
		itemDTO = append(itemDTO, mapper.ToFCWWNEntryDTO(item))
	}

	return c.JSON(http.StatusOK, itemDTO)
}

func (h *FCWWNEntryHandler) DeleteFCWWNEntry(c echo.Context) error {
	probeID := c.Param("id")
	_, err := h.service.Get(c.Request().Context(), probeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	err = h.service.Delete(c.Request().Context(), probeID)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to delete entry", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h *FCWWNEntryHandler) SoftDeleteFCWWNEntry(c echo.Context) error {
	probeID := c.Param("id")
	entry, err := h.service.Get(c.Request().Context(), probeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	entries, err := h.service.Find(c.Request().Context(), service.Filter{"wwn": entry.WWN}, service.SortOption{})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to find common entries", err)
	}

	for _, e := range entries {
		e.DuplicateCustomers = nil
		_, err = h.service.Update(c.Request().Context(), e.ID, &e)
		if err != nil {
			return errorWithInternal(http.StatusInternalServerError, "Failed to flush duplicate customers from entry", err)
		}
	}

	err = h.service.SoftDelete(c.Request().Context(), probeID)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to soft delete entry", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h *FCWWNEntryHandler) RestoreFCWWNEntry(c echo.Context) error {
	probeID := c.Param("id")

	err := h.service.Restore(c.Request().Context(), probeID)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to restore entry", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h *FCWWNEntryHandler) CreateUpdateFCWWNEntry(c echo.Context) error {
	var itemDTO dto.FCWWNEntryDTO
	if err := json.NewDecoder(c.Request().Body).Decode(&itemDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := validate.Struct(itemDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	item := mapper.ToFCWWNEntryEntity(itemDTO)

	id, err := h.service.Update(c.Request().Context(), item.ID, &item)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to update entry", err)
	}

	itemTmp, err := h.service.Get(c.Request().Context(), id.Hex())
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get entry", err)
	}

	itemDTO = mapper.ToFCWWNEntryDTO(*itemTmp)
	return c.JSON(http.StatusOK, itemDTO)
}

func (h *FCWWNEntryHandler) ImportHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return errorWithInternal(http.StatusBadRequest, "Failed to load import file from request", err)
	}

	wwnEntries, err := h.readEntriesFromFile(file)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to read entries from import file", err)
	}

	err = h.service.DeleteAll(c.Request().Context())
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to delete entries", err)
	}

	err = h.service.InsertAll(c.Request().Context(), wwnEntries)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to insert entries", err)
	}
	err = h.service.FlagDuplicateWWNs(c.Request().Context(), service.Filter{})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to flag duplicate entries", err)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Import successful"})
}

func (h *FCWWNEntryHandler) ImportAPIHandler(c echo.Context) error {
	wwnEntries, err := h.readEntriesFromAPI(c)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to read entries from api", err)
	}

	err = h.service.DeleteAll(c.Request().Context())
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to delete entries", err)
	}

	err = h.service.InsertAll(c.Request().Context(), wwnEntries)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to insert entries", err)
	}
	err = h.service.FlagDuplicateWWNs(c.Request().Context(), service.Filter{})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to flag duplicate entries", err)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Import successful"})
}

func (h *FCWWNEntryHandler) ListCustomers(c echo.Context) error {
	customers, err := h.service.Customers(c.Request().Context())
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get customers", err)
	}
	return c.JSON(http.StatusOK, customers)
}

func (h *FCWWNEntryHandler) ExportReconcileEntries(c echo.Context) error {
	items, err := h.service.Find(c.Request().Context(),
		service.Filter{
			"type":            service.Filter{"$in": []string{"Host", "Other"}},
			"wwn_set":         service.Filter{"$in": []int{1, 2}},
			"needs_reconcile": true,
			"ignore_entry":    false,
		}, service.SortOption{"wwn": "asc"})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get reconcile entries", err)
	}

	f, err := os.CreateTemp("", "reconcilecsv-")
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to create temp csv file", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	writer := csv.NewWriter(f)

	err = writer.Write([]string{"Customer", "WWN", "Zones", "Aliases", "Hostname (Generated)", "Hostname (Loaded)"})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to write csv file", err)
	}

	for _, item := range items {
		itemDTO := mapper.ToFCWWNEntryDTO(item)
		err := writer.Write([]string{itemDTO.Customer, itemDTO.WWN, strings.Join(itemDTO.Zones, ","), strings.Join(itemDTO.Aliases, ","), itemDTO.Hostname, itemDTO.LoadedHostname})
		if err != nil {
			return errorWithInternal(http.StatusInternalServerError, "Failed to write csv file", err)
		}
	}
	writer.Flush()
	return c.Attachment(f.Name(), "records_to_reconcile.csv")
}

func (h *FCWWNEntryHandler) readEntriesFromFile(file *multipart.FileHeader) ([]*entity.FCWWNEntry, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open entry file: %w", err)
	}
	defer src.Close()

	reader := csv.NewReader(src)
	reader.Comma = ','
	reader.FieldsPerRecord = -1

	var wwnEntries []*entity.FCWWNEntry
	wwnEntryMap := make(map[string]map[string]entity.FCWWNEntry, 0)

	re := regexp.MustCompile(`^([0-9A-Fa-f]{2}:){7}[0-9A-Fa-f]{2}$`)

	lineNumber := 0
	for {
		lineNumber = lineNumber + 1
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read entry file: %w", err)
		}
		if lineNumber < 3 {
			// ignore first 2 lines from ITA csv export
			continue
		}

		if len(line) < 7 {
			continue
		}

		if !re.MatchString(line[1]) {
			h.logger.Info().Msgf("Invalid WWN: %s", line[1])
			continue
		}

		customer := line[0]
		if customer == "" {
			customer = entity.UnknownCustomer
		}
		wwn := line[1]
		zone := line[2]
		alias := line[3]
		loadedHostname := line[4]
		isCsvLoad := true
		if line[5] == "N" {
			isCsvLoad = false
		}
		wwnSet, _ := strconv.Atoi(line[6])

		if loadedHostname == "No Matching Rule" {
			loadedHostname = ""
		}

		updateEntryMap(wwnEntryMap, customer, wwn, zone, alias, loadedHostname, isCsvLoad, wwnSet)
	}

	for _, v := range wwnEntryMap {
		for _, e := range v {
			wwnEntries = append(wwnEntries, &e)
		}
	}
	return wwnEntries, nil
}

func (h *FCWWNEntryHandler) readEntriesFromAPI(c echo.Context) ([]*entity.FCWWNEntry, error) {

	var wwnEntries []*entity.FCWWNEntry
	wwnEntryMap := make(map[string]map[string]entity.FCWWNEntry, 0)

	page := 0
	pageSize := 10000
	for {
		var response ita.FeedResponse
		itaClient, err := ita.NewITAClient(h.logger)
		if err != nil {
			return nil, fmt.Errorf("cannot create ITA client: %v", err)
		}
		if os.Getenv("ITA_FEED_ID") == "" {
			return nil, fmt.Errorf("ITA_FEED_ID environment variable is not set")
		}
		resp, err := itaClient.GenerateReportTemplate(c.Request().Context(), os.Getenv("ITA_FEED_ID"), page, pageSize)
		if err != nil {
			return nil, fmt.Errorf("cannot get feed report: %v", err)
		}

		if err := json.Unmarshal(resp, &response); err != nil {
			return nil, fmt.Errorf("cannot unmarshall feed report: %v", err)
		}

		for _, line := range response.Data.Report.ReportData {
			entryLine, err := h.parseEntryLine(line)
			if err != nil {
				h.logger.Error().Msgf("Failed to parse entry line: %v", err)
				continue
			}

			updateEntryMap(wwnEntryMap, entryLine.customer, entryLine.wwn, entryLine.zone, entryLine.alias, entryLine.loadedHostname, entryLine.isCsvLoad, entryLine.wwnSet)

		}

		if response.Data.Paging.Next == 0 {
			break
		}
		page = response.Data.Paging.Next
	}

	for _, v := range wwnEntryMap {
		for _, e := range v {
			wwnEntries = append(wwnEntries, &e)
		}
	}
	return wwnEntries, nil
}

func (h *FCWWNEntryHandler) parseEntryLine(line ita.ReportData) (entryLineRecord, error) {
	entryLine := entryLineRecord{}
	wwn, ok := line["wwn"].Value.(string)
	if !ok {
		return entryLine, fmt.Errorf("WWN type is not string: %T", line["wwn"].Value)
	}
	if !wwnRegex.MatchString(wwn) {
		return entryLine, fmt.Errorf("invalid WWN: %s", wwn)
	}
	entryLine.wwn = wwn

	customer, ok := line["customer"].Value.(string)
	if !ok {
		return entryLine, fmt.Errorf("customer type is not string: %T", line["customer"].Value)
	}
	if customer == "" {
		customer = entity.UnknownCustomer
	}
	entryLine.customer = customer
	zone, ok := line["element_name"].Value.(string)
	if !ok {
		return entryLine, fmt.Errorf("zone type is not string: %T", line["element_name"].Value)
	}
	entryLine.zone = zone
	alias, ok := line["alias"].Value.(string)
	if !ok {
		return entryLine, fmt.Errorf("alias type is not string: %T", line["alias"].Value)
	}
	entryLine.alias = alias
	loadedHostname, ok := line["loaded_host"].Value.(string)
	if !ok {
		return entryLine, fmt.Errorf("loaded_host type is not string: %T", line["loaded_host"].Value)
	}
	entryLine.loadedHostname = loadedHostname
	isCsvLoad := true
	csvLoad, ok := line["is_csv_load"].Value.(string)
	if !ok {
		return entryLine, fmt.Errorf("is_csv_load type is not string: %T", line["is_csv_load"].Value)
	}
	if csvLoad == "N" {
		isCsvLoad = false
	}
	entryLine.isCsvLoad = isCsvLoad
	wwnSetF, ok := line["wwn_set"].Value.(float64)
	if !ok {
		return entryLine, fmt.Errorf("wwn_set type is not numberic: %T", line["wwn_set"].Value)
	}
	wwnSet := int(wwnSetF)
	entryLine.wwnSet = wwnSet

	if loadedHostname == "No Matching Rule" {
		entryLine.loadedHostname = ""
	}
	return entryLine, nil
}
func updateEntryMap(wwnEntryMap map[string]map[string]entity.FCWWNEntry, customer, wwn, zone, alias, loadedHostname string, isCsvLoad bool, wwnSet int) {

	// Create or update FCWWNEntry
	if cMap, ok := wwnEntryMap[customer]; ok {
		if wwnEntry, ok := cMap[wwn]; ok {
			// Update existing entry
			if zone != "" && !contains(wwnEntry.Zones, zone) {
				wwnEntry.Zones = append(wwnEntry.Zones, zone)
			}
			if alias != "" && !contains(wwnEntry.Aliases, alias) {
				wwnEntry.Aliases = append(wwnEntry.Aliases, alias)
			}
			wwnEntryMap[customer][wwn] = wwnEntry
		} else {
			// Create new entry for this WWN
			newWWNEntry := entity.FCWWNEntry{
				Customer:       customer,
				WWN:            wwn,
				Zones:          []string{},
				Aliases:        []string{},
				LoadedHostname: loadedHostname,
				IsCSVLoad:      isCsvLoad,
				WWNSet:         wwnSet,
			}
			if zone != "" {
				newWWNEntry.Zones = append(newWWNEntry.Zones, zone)
			}
			if alias != "" {
				newWWNEntry.Aliases = append(newWWNEntry.Aliases, alias)
			}
			wwnEntryMap[customer][wwn] = newWWNEntry
		}
	} else {
		// Create new map for this customer and add the WWN entry
		newWWNEntry := entity.FCWWNEntry{
			Customer:       customer,
			WWN:            wwn,
			Zones:          []string{},
			Aliases:        []string{},
			LoadedHostname: loadedHostname,
			IsCSVLoad:      isCsvLoad,
			WWNSet:         wwnSet,
		}
		if zone != "" {
			newWWNEntry.Zones = append(newWWNEntry.Zones, zone)
		}
		if alias != "" {
			newWWNEntry.Aliases = append(newWWNEntry.Aliases, alias)
		}
		wwnEntryMap[customer] = map[string]entity.FCWWNEntry{
			wwn: newWWNEntry,
		}
	}
}
