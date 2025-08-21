package handler

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/wwn_identity/webapi/internal/mapper"
	"github.com/ttrnecka/wwn_identity/webapi/internal/service"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RuleHandler struct {
	service service.RuleService
}

func NewRuleHandler(s service.RuleService) *RuleHandler {
	return &RuleHandler{s}
}

func (h *RuleHandler) GetRules(c echo.Context) error {
	customer := c.Param("name")
	rules, err := h.service.Find(c.Request().Context(), bson.M{"customer": customer}, options.Find().SetSort(bson.M{"order": 1}))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	itemsDTO := []dto.RuleDTO{}
	for _, item := range rules {
		itemsDTO = append(itemsDTO, mapper.ToRuleDTO(item))
	}
	return c.JSON(http.StatusOK, itemsDTO)
}

func (h *RuleHandler) Rules(c echo.Context) error {

	items, err := h.service.All(c.Request().Context())
	if err != nil {
		return err
	}
	itemsDTO := []dto.RuleDTO{}
	for _, item := range items {
		itemsDTO = append(itemsDTO, mapper.ToRuleDTO(item))
	}
	return c.JSON(http.StatusOK, itemsDTO)
}

func (h *RuleHandler) ExportRules(c echo.Context) error {

	items, err := h.service.All(c.Request().Context())
	if err != nil {
		return err
	}

	f, err := os.CreateTemp("", "exportcsv-")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	writer := csv.NewWriter(f)

	for _, item := range items {
		itemDTO := mapper.ToRuleDTO(item)
		writer.Write([]string{strconv.Itoa(itemDTO.Order), itemDTO.Customer, itemDTO.Regex, strconv.Itoa(itemDTO.Group), string(itemDTO.Type)})
	}
	writer.Flush()
	return c.Attachment(f.Name(), "rules.csv")
}

func (h *RuleHandler) DeleteRule(c echo.Context) error {
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

func (h *RuleHandler) CreateUpdateRule(c echo.Context) error {
	customer := c.Param("name")
	var itemDTO dto.RuleDTO
	if err := json.NewDecoder(c.Request().Body).Decode(&itemDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	itemDTO.Customer = customer

	if err := validate.Struct(itemDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	item := mapper.ToRuleEntity(itemDTO)

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

	itemDTO = mapper.ToRuleDTO(*itemTmp)
	return c.JSON(http.StatusOK, itemDTO)
}

func (h *RuleHandler) CreateUpdateRules(c echo.Context) error {
	customer := c.Param("name")
	var itemsDTO []dto.RuleDTO
	if err := json.NewDecoder(c.Request().Body).Decode(&itemsDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	for _, item := range itemsDTO {
		item.Customer = customer
		if err := validate.Struct(item); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		item := mapper.ToRuleEntity(item)

		_, err := h.service.Update(c.Request().Context(), item.ID, &item)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
	}

	return c.NoContent(http.StatusOK)
}
