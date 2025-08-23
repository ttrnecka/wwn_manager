package handler

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/mapper"
	"github.com/ttrnecka/wwn_identity/webapi/internal/service"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
	"go.mongodb.org/mongo-driver/bson"
)

type FCEntryHandler struct {
	service     service.FCEntryService
	ruleService service.RuleService
	wwnService  service.FCWWNEntryService
}

func NewFCEntryHandler(s service.FCEntryService, r service.RuleService, w service.FCWWNEntryService) *FCEntryHandler {
	return &FCEntryHandler{s, r, w}
}

func (h *FCEntryHandler) FCEntries(c echo.Context) error {
	customer := c.Param("name")

	items, err := h.service.Find(c.Request().Context(), bson.M{"customer": customer})
	if err != nil {
		return err
	}

	var iteamDTO []dto.FCEntryDTO

	for _, item := range items {
		iteamDTO = append(iteamDTO, mapper.ToFCEntryDTO(item))
	}

	return c.JSON(http.StatusOK, iteamDTO)
}

func (h *FCEntryHandler) DeleteFCEntry(c echo.Context) error {
	probe_id := c.Param("id")
	_, err := h.service.Get(c.Request().Context(), probe_id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}
	err = h.service.Delete(c.Request().Context(), probe_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return c.NoContent(http.StatusOK)
}

func (h *FCEntryHandler) CreateUpdateFCEntry(c echo.Context) error {
	var itemDTO dto.FCEntryDTO
	if err := json.NewDecoder(c.Request().Body).Decode(&itemDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := validate.Struct(itemDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	item := mapper.ToFCEntryEntity(itemDTO)

	id, err := h.service.Update(c.Request().Context(), item.ID, &item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})

	}

	itemTmp, err := h.service.Get(c.Request().Context(), id.Hex())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	itemDTO = mapper.ToFCEntryDTO(*itemTmp)
	return c.JSON(http.StatusOK, itemDTO)
}

func (h *FCEntryHandler) ImportHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	defer src.Close()

	reader := csv.NewReader(src)
	reader.Comma = ',' // switch to '\t' if TSV

	var entries []entity.FCEntry
	var wwnEntries []entity.FCWWNEntry
	wwnEntryMap := make(map[string]map[string]*entity.FCWWNEntry, 0)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
		if len(line) < 4 {
			continue
		}
		entry := entity.FCEntry{
			Customer:       line[0],
			WWN:            line[1],
			Zone:           line[2],
			Alias:          line[3],
			LoadedHostname: line[4],
		}
		entries = append(entries, entry)

		// Create or update FCWWNEntry
		if cMap, ok := wwnEntryMap[entry.Customer]; ok {
			if wwnEntry, ok := cMap[entry.WWN]; ok {
				// Update existing entry
				if entry.Zone != "" && !contains(wwnEntry.Zones, entry.Zone) {
					wwnEntry.Zones = append(wwnEntry.Zones, entry.Zone)
				}
				if entry.Alias != "" && !contains(wwnEntry.Aliases, entry.Alias) {
					wwnEntry.Aliases = append(wwnEntry.Aliases, entry.Alias)
				}
				wwnEntryMap[entry.Customer][entry.WWN] = wwnEntry
			} else {
				// Create new entry for this WWN
				newWWNEntry := entity.FCWWNEntry{
					Customer:       entry.Customer,
					WWN:            entry.WWN,
					Zones:          []string{},
					Aliases:        []string{},
					LoadedHostname: entry.LoadedHostname,
				}
				if entry.Zone != "" {
					newWWNEntry.Zones = append(newWWNEntry.Zones, entry.Zone)
				}
				if entry.Alias != "" {
					newWWNEntry.Aliases = append(newWWNEntry.Aliases, entry.Alias)
				}
				wwnEntryMap[entry.Customer][entry.WWN] = &newWWNEntry
				wwnEntries = append(wwnEntries, newWWNEntry)
			}
		} else {
			// Create new map for this customer and add the WWN entry
			newWWNEntry := entity.FCWWNEntry{
				Customer:       entry.Customer,
				WWN:            entry.WWN,
				Zones:          []string{},
				Aliases:        []string{},
				LoadedHostname: entry.LoadedHostname,
			}
			if entry.Zone != "" {
				newWWNEntry.Zones = append(newWWNEntry.Zones, entry.Zone)
			}
			if entry.Alias != "" {
				newWWNEntry.Aliases = append(newWWNEntry.Aliases, entry.Alias)
			}
			wwnEntryMap[entry.Customer] = map[string]*entity.FCWWNEntry{
				entry.WWN: &newWWNEntry,
			}
			wwnEntries = append(wwnEntries, newWWNEntry)
		}
	}

	err = h.service.DeleteAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	err = h.service.InsertAll(c.Request().Context(), entries)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	err = h.wwnService.DeleteAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	err = h.wwnService.InsertAll(c.Request().Context(), wwnEntries)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Import successful"})
}

func (h *FCEntryHandler) ListCustomers(c echo.Context) error {
	customers, err := h.service.Customers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, customers)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
