package handler

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/mapper"
	"github.com/ttrnecka/wwn_identity/webapi/internal/service"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FCEntryHandler struct {
	service     service.FCEntryService
	ruleService service.RuleService
}

func NewFCEntryHandler(s service.FCEntryService, r service.RuleService) *FCEntryHandler {
	return &FCEntryHandler{s, r}
}

func (h *FCEntryHandler) FCEntries(c echo.Context) error {
	customer := c.Param("name")

	items, err := h.service.Find(c.Request().Context(), bson.M{"customer": customer})
	if err != nil {
		return err
	}

	rules, err := h.ruleService.Find(c.Request().Context(), bson.M{"customer": customer}, options.Find().SetSort(bson.M{"order": 1}))
	if err != nil {
		return err
	}

	itemsDTO := []dto.FCEntryDTO{}
	for i, item := range items {
		itemsDTO = append(itemsDTO, mapper.ToFCEntryDTO(item))
		itemsDTO[i].Hostname = applyRules(item, rules)
	}
	return c.JSON(http.StatusOK, itemsDTO)
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
			Customer: line[0],
			WWN:      line[1],
			Zone:     line[2],
			Alias:    line[3],
		}
		entries = append(entries, entry)
	}

	err = h.service.DeleteAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	for _, entry := range entries {
		_, err := h.service.Update(c.Request().Context(), entry.ID, &entry)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
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

func applyRules(entry entity.FCEntry, rules []entity.Rule) string {
	hostname := ""
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Order < rules[j].Order
	})

TOP:
	for _, rule := range rules {
		r, err := regexp.Compile(rule.Regex)
		if err != nil {
			continue
		}
		switch rule.Type {
		case entity.ZoneRule:
			match := r.FindStringSubmatch(entry.Zone)
			if len(match) > 1 {
				hostname = match[1] // first capture group
				break TOP
			}
		case entity.AliasRule:
			match := r.FindStringSubmatch(entry.Alias)
			if len(match) > 1 {
				hostname = match[1] // first capture group
				break TOP
			}
		}
	}
	return hostname
}
