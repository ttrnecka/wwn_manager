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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FCWWNEntryHandler struct {
	service     service.FCWWNEntryService
	ruleService service.RuleService
}

func NewFCWWNEntryHandler(s service.FCWWNEntryService, r service.RuleService) *FCWWNEntryHandler {
	return &FCWWNEntryHandler{s, r}
}

func (h *FCWWNEntryHandler) FCWWNEntries(c echo.Context) error {
	customer := c.Param("name")

	items, err := h.service.Find(c.Request().Context(), bson.M{"customer": customer}, options.Find().SetSort(bson.M{"wwn": 1}))
	if err != nil {
		return err
	}

	var itemDTO []dto.FCWWNEntryDTO

	for _, item := range items {
		itemDTO = append(itemDTO, mapper.ToFCWWNEntryDTO(item))
	}

	return c.JSON(http.StatusOK, itemDTO)
}

func (h *FCWWNEntryHandler) DeleteFCWWNEntry(c echo.Context) error {
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

	itemDTO = mapper.ToFCWWNEntryDTO(*itemTmp)
	return c.JSON(http.StatusOK, itemDTO)
}

func (h *FCWWNEntryHandler) ImportHandler(c echo.Context) error {
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
	reader.Comma = ','

	var wwnEntries []entity.FCWWNEntry
	wwnEntryMap := make(map[string]map[string]entity.FCWWNEntry, 0)

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

		customer := line[0]
		wwn := line[1]
		zone := line[2]
		alias := line[3]
		loadedHostname := line[4]

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

	err = h.service.DeleteAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	err = h.service.InsertAll(c.Request().Context(), wwnEntries)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Import successful"})
}

func (h *FCWWNEntryHandler) ListCustomers(c echo.Context) error {
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
