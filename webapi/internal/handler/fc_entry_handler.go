package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/mapper"
	"github.com/ttrnecka/wwn_identity/webapi/internal/service"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
)

type FCWWNEntryHandler struct {
	service     service.FCWWNEntryService
	ruleService service.RuleService
}

func NewFCWWNEntryHandler(s service.FCWWNEntryService, r service.RuleService) *FCWWNEntryHandler {
	return &FCWWNEntryHandler{s, r}
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
	probe_id := c.Param("id")
	_, err := h.service.Get(c.Request().Context(), probe_id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	err = h.service.Delete(c.Request().Context(), probe_id)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to delete entry", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h *FCWWNEntryHandler) SoftDeleteFCWWNEntry(c echo.Context) error {
	probe_id := c.Param("id")
	_, err := h.service.Get(c.Request().Context(), probe_id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	err = h.service.SoftDelete(c.Request().Context(), probe_id)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to soft delete entry", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h *FCWWNEntryHandler) RestoreFCWWNEntry(c echo.Context) error {
	probe_id := c.Param("id")

	err := h.service.Restore(c.Request().Context(), probe_id)
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

func (h *FCWWNEntryHandler) ListCustomers(c echo.Context) error {
	customers, err := h.service.Customers(c.Request().Context())
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get customers", err)
	}
	return c.JSON(http.StatusOK, customers)
}

func (h *FCWWNEntryHandler) ExportHostWWNMap(c echo.Context) error {
	items, err := h.service.All(c.Request().Context())
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get rules", err)
	}

	f, err := os.CreateTemp("", "exportcsv-")
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to create temp csv file", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	writer := csv.NewWriter(f)

	for _, item := range items {
		itemDTO := mapper.ToFCWWNEntryDTO(item)
		if itemDTO.IsPrimaryCustomer {
			writer.Write([]string{itemDTO.Hostname, itemDTO.WWN, itemDTO.WWN})
		}
	}
	writer.Flush()
	return c.Attachment(f.Name(), "host_wwn.csv")
}

func (h *FCWWNEntryHandler) ExportCustomerWWNMap(c echo.Context) error {
	items, err := h.service.Find(c.Request().Context(), service.Filter{"is_primary_customer": false}, service.SortOption{"wwn": "asc"})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get rules", err)
	}

	f, err := os.CreateTemp("", "exportcsv-")
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to create temp csv file", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	writer := csv.NewWriter(f)

	for _, item := range items {
		itemDTO := mapper.ToFCWWNEntryDTO(item)
		writer.Write([]string{itemDTO.WWN, itemDTO.Customer, itemDTO.Hostname})
	}
	writer.Flush()
	return c.Attachment(f.Name(), "customer_wwn_host_override.csv")
}

func (h *FCWWNEntryHandler) readEntriesFromFile(file *multipart.FileHeader) ([]entity.FCWWNEntry, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open entry file: %w", err)
	}
	defer src.Close()

	reader := csv.NewReader(src)
	reader.Comma = ','
	reader.FieldsPerRecord = -1

	var wwnEntries []entity.FCWWNEntry
	wwnEntryMap := make(map[string]map[string]entity.FCWWNEntry, 0)

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

		customer := line[0]
		if customer == "" {
			customer = entity.UNKNOWN_CUSTOMER
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

	for _, v := range wwnEntryMap {
		for _, e := range v {
			wwnEntries = append(wwnEntries, e)
		}
	}
	return wwnEntries, nil
}
