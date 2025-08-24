package handler

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/mapper"
	"github.com/ttrnecka/wwn_identity/webapi/internal/service"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
	"go.mongodb.org/mongo-driver/bson"
)

type RuleHandler struct {
	service           service.RuleService
	fcWWNEntryService service.FCWWNEntryService
}

func NewRuleHandler(s service.RuleService, w service.FCWWNEntryService) *RuleHandler {
	return &RuleHandler{s, w}
}

func (h *RuleHandler) GetRules(c echo.Context) error {
	customer := c.Param("name")
	rules, err := h.service.Find(c.Request().Context(), service.Filter{"customer": customer}, service.SortOption{"order": "asc"})
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
	mode := c.QueryParam("mode")
	// bulk mode with array of rules
	if mode == "bulk" {
		return h.CreateUpdateRules(c)
	}
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
	fcWWNEntries, err := h.fcWWNEntryService.All(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = h.applyRules(c.Request().Context(), fcWWNEntries)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *RuleHandler) SetupAndApplyReconcileRules(c echo.Context) error {
	fcWWNEntryId := c.Param("id")

	// get entry
	fcWWNEntry, err := h.fcWWNEntryService.Get(c.Request().Context(), fcWWNEntryId)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	// get reconciliation rules and fix entry
	var reconcileDTO dto.EntryReconcileDTO
	if err := json.NewDecoder(c.Request().Body).Decode(&reconcileDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	// save rules
	err = h.service.CreateReconcileRules(c.Request().Context(), fcWWNEntry, reconcileDTO)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// update entry in db
	_, err = h.fcWWNEntryService.Update(c.Request().Context(), fcWWNEntry.ID, fcWWNEntry)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// get all entries for given WWN

	entries, err := h.fcWWNEntryService.Find(c.Request().Context(), service.Filter{"wwn": fcWWNEntry.WWN}, service.SortOption{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = h.applyRules(c.Request().Context(), entries)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// apply rules for selected entries
	return c.NoContent(http.StatusOK)
}

func (h *RuleHandler) applyRules(ctx context.Context, fcWWNEntries []entity.FCWWNEntry) error {

	numWorkers := runtime.NumCPU() // one worker per CPU core

	var wg sync.WaitGroup
	wg.Add(len(fcWWNEntries))

	// channel to distribute indices
	idxCh := make(chan int)

	globalRules, err := h.service.Find(ctx, service.Filter{"customer": entity.GLOBAL_CUSTOMER}, service.SortOption{"order": "asc"})
	if err != nil {
		return err
	}

	ruleMap := make(map[string][]entity.Rule)

	mutex := sync.Mutex{}

	// start workers
	for range numWorkers {
		go func() {
			for i := range idxCh {
				mutex.Lock()
				rules, ok := ruleMap[fcWWNEntries[i].Customer]
				mutex.Unlock()
				if !ok {
					rules, err = h.service.Find(ctx, service.Filter{"customer": fcWWNEntries[i].Customer}, service.SortOption{"order": "asc"})
					if err != nil {
						continue
					}
					rules = append(rules, globalRules...)
					mutex.Lock()
					ruleMap[fcWWNEntries[i].Customer] = rules
					mutex.Unlock()
				}
				err = applyRules(&fcWWNEntries[i], rules)
				wg.Done()
			}
		}()
	}

	// send work
	for i := range fcWWNEntries {
		idxCh <- i
	}
	close(idxCh)

	wg.Wait()

	wwns := make([]string, 0)
	for _, e := range fcWWNEntries {
		wwns = append(wwns, e.WWN)
	}

	err = h.fcWWNEntryService.DeleteMany(ctx, service.Filter{"wwn": bson.M{"$in": wwns}})
	if err != nil {
		return err
	}

	err = h.fcWWNEntryService.InsertAll(ctx, fcWWNEntries)
	return err
}

func applyRules(entry *entity.FCWWNEntry, rules []entity.Rule) error {
	entry.Hostname = ""
	entry.Type = "Unknown"
	entry.NeedsReconcile = false
	entry.IsPrimaryCustomer = false
	entry.DuplicateRule = entity.NilObjectID()

	// RANGE rules
RANGE:
	for _, rule := range rules {
		r, err := regexp.Compile(rule.Regex)
		if err != nil {
			return err
		}
		entry.TypeRule = rule.ID
		switch rule.Type {
		case entity.WWNArrayRangeRule:
			match := r.MatchString(entry.WWN)
			if match {
				entry.Type = "Array"
				break RANGE
			}
		case entity.WWNBackupRangeRule:
			match := r.MatchString(entry.WWN)
			if match {
				entry.Type = "Backup"
				break RANGE
			}
		case entity.WWNHostRangeRule:
			match := r.MatchString(entry.WWN)
			if match {
				entry.Type = "Host"
				break RANGE
			}
		case entity.WWNOtherRangeRule:
			match := r.MatchString(entry.WWN)
			if match {
				entry.Type = "Other"
				break RANGE
			}
		}
		entry.TypeRule = entity.NilObjectID()
	}
	// do host check only for host ranges
	if entry.Type == "Array" || entry.Type == "Backup" || entry.Type == "Other" {
		return nil
	}
	// MAP rules
TOP:
	for _, rule := range rules {
		if rule.Group < 0 {
			continue
		}
		r, err := regexp.Compile(rule.Regex)
		if err != nil {
			return err
		}
		entry.HostNameRule = rule.ID
		switch rule.Type {
		case entity.ZoneRule:
			for _, zone := range entry.Zones {
				match := r.FindStringSubmatch(zone)
				if len(match) > 0 && len(match) >= rule.Group {
					entry.Hostname = match[rule.Group] // first capture group
					break TOP
				}
			}
		case entity.AliasRule:
			for _, alias := range entry.Aliases {
				match := r.FindStringSubmatch(alias)
				if len(match) > 0 && len(match) >= rule.Group {
					entry.Hostname = match[rule.Group] // first capture group
					break TOP
				}
			}
		case entity.WWNHostMapRule:
			match := r.FindStringSubmatch(entry.WWN)
			if len(match) > 0 {
				entry.Hostname = rule.Comment
				break TOP
			}
		}
		entry.HostNameRule = entity.NilObjectID()
	}

	// DUP rules
	dupReconciled := false
DUP:
	for _, rule := range rules {
		if rule.Type != entity.WWNCustomerMapRule {
			continue
		}
		if entry.WWN == rule.Regex {
			entry.DuplicateRule = rule.ID
			entry.IsPrimaryCustomer = false
			if entry.Customer == rule.Comment {
				entry.IsPrimaryCustomer = true
			}
			dupReconciled = true
			break DUP
		}
	}
	if len(entry.DuplicateCustomers) > 0 && !dupReconciled ||
		(entry.LoadedHostname != "" && entry.Hostname != entry.LoadedHostname) {
		entry.NeedsReconcile = true
	}

	return nil
}
