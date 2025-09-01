package handler

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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
		return errorWithInternal(http.StatusInternalServerError, "Failed to find rules", err)
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
		return errorWithInternal(http.StatusInternalServerError, "Failed to get rules", err)
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
		itemDTO := mapper.ToRuleDTO(item)
		writer.Write([]string{strconv.Itoa(itemDTO.Order), itemDTO.Customer, itemDTO.Regex, strconv.Itoa(itemDTO.Group), string(itemDTO.Type), itemDTO.Comment})
	}
	writer.Flush()
	return c.Attachment(f.Name(), "rules.csv")
}

func (h *RuleHandler) DeleteRule(c echo.Context) error {
	probe_id := c.Param("id")
	_, err := h.service.Get(c.Request().Context(), probe_id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	err = h.service.Delete(c.Request().Context(), probe_id)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to delete rules", err)
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
		return errorWithInternal(http.StatusInternalServerError, "Failed to update rule", err)
	}

	itemTmp, err := h.service.Get(c.Request().Context(), id.Hex())
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get rule", err)
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
			return errorWithInternal(http.StatusInternalServerError, "Failed to update rule", err)
		}
	}
	fcWWNEntries, err := h.fcWWNEntryService.All(c.Request().Context())
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get entries", err)
	}
	err = h.applyRules(c.Request().Context(), fcWWNEntries)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to apply rules", err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *RuleHandler) ApplyRules(c echo.Context) error {

	fcWWNEntries, err := h.fcWWNEntryService.All(c.Request().Context())
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get entries", err)
	}
	err = h.applyRules(c.Request().Context(), fcWWNEntries)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to apply rules", err)
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
		return errorWithInternal(http.StatusInternalServerError, "Failed to create reconcile rules", err)
	}

	// update entry in db
	_, err = h.fcWWNEntryService.Update(c.Request().Context(), fcWWNEntry.ID, fcWWNEntry)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to update entry", err)
	}

	// howeve we need to pass them again in this case as hostname reconciliation requires customer reconciliation to be done first
	// get all entries for given WWN
	entries, err := h.fcWWNEntryService.Find(c.Request().Context(), service.Filter{"wwn": fcWWNEntry.WWN}, service.SortOption{})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to file entries", err)
	}

	err = h.applyRules(c.Request().Context(), entries)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to apply rules", err)
	}

	// apply rules for selected entries
	return c.NoContent(http.StatusOK)
}

func (h *RuleHandler) ImportHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return errorWithInternal(http.StatusBadRequest, "Failed to load import file from request", err)
	}

	rules, err := h.readEntriesFromFile(file)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to read rules from import file", err)
	}

	err = h.service.DeleteAll(c.Request().Context())
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to delete rules", err)
	}

	err = h.service.InsertAll(c.Request().Context(), rules)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to insert rules", err)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Import successful"})
}

func (h *RuleHandler) applyRules(ctx context.Context, fcWWNEntries []entity.FCWWNEntry) error {

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU() // one worker per CPU core

	globalRules, err := h.service.Find(ctx, service.Filter{"customer": entity.GLOBAL_CUSTOMER, "type": service.Filter{"$in": entity.RangeRules}}, service.SortOption{"order": "asc"})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get GLOBAL rules", err)
	}

	// RANGE RULES START
	wg.Add(len(fcWWNEntries))
	idxCh := make(chan int)

	for range numWorkers {
		go func() {
			for i := range idxCh {
				err = applyRangeRules(&fcWWNEntries[i], globalRules)
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
	// RANGE RULES STOP

	// HOST RULES START
	wg.Add(len(fcWWNEntries))
	idxCh = make(chan int)

	ruleMap := make(map[string][]entity.Rule)

	globalHostRules, err := h.service.Find(ctx, service.Filter{"customer": entity.GLOBAL_CUSTOMER, "type": service.Filter{"$in": entity.HostRules}}, service.SortOption{"order": "asc"})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get GLOBAL rules", err)
	}

	mutex := sync.Mutex{}

	// start workers
	for range numWorkers {
		go func() {
			for i := range idxCh {
				mutex.Lock()
				rules, ok := ruleMap[fcWWNEntries[i].Customer]
				mutex.Unlock()
				if !ok {
					rules, err = h.service.Find(ctx, service.Filter{"customer": fcWWNEntries[i].Customer, "type": service.Filter{"$in": entity.HostRules}}, service.SortOption{"order": "asc"})
					if err != nil {
						continue
					}
					rules = append(rules, globalHostRules...)
					mutex.Lock()
					ruleMap[fcWWNEntries[i].Customer] = rules
					mutex.Unlock()
				}
				err = applyHostRules(&fcWWNEntries[i], rules)
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
	// HOST RULES STOP

	// update entries in db so we can run FlagDuplicateWWNs to get data for reconcile rules
	wwns := make([]string, 0)
	for _, e := range fcWWNEntries {
		wwns = append(wwns, e.WWN)
	}

	ids := make([]entity.ID, 0)
	for _, e := range fcWWNEntries {
		ids = append(ids, e.ID)
	}

	err = h.fcWWNEntryService.DeleteMany(ctx, service.Filter{"_id": bson.M{"$in": ids}})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to delete entries", err)
	}

	err = h.fcWWNEntryService.InsertAll(ctx, fcWWNEntries)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to insert entries", err)
	}

	err = h.fcWWNEntryService.FlagDuplicateWWNs(ctx, service.Filter{"wwn": wwns})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to flag duplicate entries", err)
	}

	fcWWNEntries, err = h.fcWWNEntryService.Find(ctx, service.Filter{"wwn": bson.M{"$in": wwns}}, service.SortOption{})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to reload entries", err)
	}

	// RECONCILE RULES START
	wg.Add(len(fcWWNEntries))
	idxCh = make(chan int)

	ruleMap = make(map[string][]entity.Rule)

	globalReconcileRules, err := h.service.Find(ctx, service.Filter{"customer": entity.GLOBAL_CUSTOMER, "type": service.Filter{"$in": entity.ReconcileRules}}, service.SortOption{"order": "asc"})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get GLOBAL rules", err)
	}

	// start workers
	for range numWorkers {
		go func() {
			for i := range idxCh {
				mutex.Lock()
				rules, ok := ruleMap[fcWWNEntries[i].Customer]
				mutex.Unlock()
				if !ok {
					rules, err = h.service.Find(ctx, service.Filter{"customer": fcWWNEntries[i].Customer, "type": service.Filter{"$in": entity.ReconcileRules}}, service.SortOption{"order": "asc"})
					if err != nil {
						continue
					}
					rules = append(rules, globalReconcileRules...)
					mutex.Lock()
					ruleMap[fcWWNEntries[i].Customer] = rules
					mutex.Unlock()
				}
				err = applyReconcileRules(&fcWWNEntries[i], rules)
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
	// RECONCILE RULES STOP

	// update fc wwn entries in db again
	err = h.fcWWNEntryService.DeleteMany(ctx, service.Filter{"wwn": bson.M{"$in": wwns}})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to delete entries", err)
	}

	err = h.fcWWNEntryService.InsertAll(ctx, fcWWNEntries)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to insert entries", err)
	}

	return nil
}

func (h *RuleHandler) readEntriesFromFile(file *multipart.FileHeader) ([]entity.Rule, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open entry file: %w", err)
	}
	defer src.Close()

	reader := csv.NewReader(src)
	reader.Comma = ','
	reader.FieldsPerRecord = -1

	rules := make([]entity.Rule, 0)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read rule file: %w", err)
		}
		if len(line) < 6 {
			continue
		}

		order, _ := strconv.Atoi(line[0])
		customer := line[1]
		regexp := line[2]
		group, _ := strconv.Atoi(line[3])
		rtype := line[4]
		comment := line[5]

		newRule := entity.Rule{
			Order:    order,
			Customer: customer,
			Regex:    regexp,
			Group:    group,
			Type:     entity.RuleType(rtype),
			Comment:  comment,
		}
		rules = append(rules, newRule)

	}
	return rules, nil
}

func applyRangeRules(entry *entity.FCWWNEntry, rules []entity.Rule) error {
	entry.Type = "Unknown"

	// RANGE rules
RANGE:
	for _, rule := range rules {
		r, err := regexp.Compile(rule.Regex)
		if err != nil {
			return fmt.Errorf("apply-rules - regex %s won't compile: %w", rule.Regex, err)
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
	return nil
}

func applyHostRules(entry *entity.FCWWNEntry, rules []entity.Rule) error {
	entry.Hostname = ""

	// do host check only for host ranges
	if entry.Type == "Array" || entry.Type == "Backup" || entry.Type == "Other" {
		return nil
	}

	// HOST rules
TOP:
	for _, rule := range rules {
		if rule.Group < 0 {
			continue
		}
		r, err := regexp.Compile(rule.Regex)
		if err != nil {
			return fmt.Errorf("apply-rules - regex %s won't compile: %w", rule.Regex, err)
		}
		entry.HostNameRule = rule.ID
		switch rule.Type {
		case entity.ZoneRule:
			for _, zone := range entry.Zones {
				match := r.FindStringSubmatch(zone)
				if len(match) > 0 && len(match) >= rule.Group {
					entry.Hostname = strings.ToLower(match[rule.Group]) // first capture group
					break TOP
				}
			}
		case entity.AliasRule:
			for _, alias := range entry.Aliases {
				match := r.FindStringSubmatch(alias)
				if len(match) > 0 && len(match) >= rule.Group {
					entry.Hostname = strings.ToLower(match[rule.Group]) // first capture group
					break TOP
				}
			}
		case entity.WWNHostMapRule:
			match := r.FindStringSubmatch(entry.WWN)
			if len(match) > 0 {
				entry.Hostname = strings.ToLower(rule.Comment)
				break TOP
			}
		}
		entry.HostNameRule = entity.NilObjectID()
	}

	// for set 2 and 3 there is no zone or alias so we just take loaded hostname as is
	if entry.WWNSet == entity.WWNSetManual || entry.WWNSet == entity.WWNSetAuto {
		entry.Hostname = strings.ToLower(entry.LoadedHostname)
	}

	return nil
}

func applyReconcileRules(entry *entity.FCWWNEntry, rules []entity.Rule) error {
	entry.NeedsReconcile = false
	entry.IsPrimaryCustomer = true
	entry.ReconcileRules = entity.NilOjectIdSlice()

	// do host check only for host ranges
	if entry.Type == "Array" || entry.Type == "Backup" || entry.Type == "Other" {
		return nil
	}

	dupReconciled := true
	//autoreconciliation for entries with auto set -> auto is always primary, rest secondary
	// for secondary hosts we make sure the decode and loaded hostname will be the smae
	if len(entry.DuplicateCustomers) > 0 {

		// check if the customer translate to unique host names
		unique := true
		seen := make(map[string]struct{})
		for _, dc := range entry.DuplicateCustomers {
			if _, exists := seen[strings.ToLower(dc.Hostname)]; exists {
				unique = false
			}
			seen[strings.ToLower(dc.Hostname)] = struct{}{}
		}
		dupReconciled = false
		entry.IsPrimaryCustomer = false
		// Auto set it automatically primary customer
		if entry.WWNSet == entity.WWNSetAuto {
			dupReconciled = true
			entry.IsPrimaryCustomer = true
		} else if entry.WWNSet == entity.WWNSetManual {
			if !unique {
				dupReconciled = true
				entry.IsPrimaryCustomer = true
			}
		} else {
			if !unique {
				entry.LoadedHostname = ""
				dupReconciled = true
			}
			// otherwise check of some other customer is auto
			// if it is we reconcile it as it should be not primary if auto set exists
			// as well the loaded hostname belongs to primary so we just flush it form secondary
			// there will be only one type 2 or type 3 set, not both
			for _, c := range entry.DuplicateCustomers {
				if c.WWNSet == entity.WWNSetAuto {
					dupReconciled = true
					entry.LoadedHostname = ""
					if strings.EqualFold(entry.Hostname, c.Hostname) {
						entry.IgnoreEntry = true
					}
				}
				// case where there is manually inserted wwn for but the customer is different but decoded hostname is same
				if c.WWNSet == entity.WWNSetManual &&
					entry.Hostname != "" &&
					strings.EqualFold(entry.Hostname, c.Hostname) {
					entry.IgnoreEntry = true
				}
			}
		}
	}
	loadedReconciled := true
	if entry.LoadedHostname != "" && !strings.EqualFold(entry.Hostname, entry.LoadedHostname) {
		loadedReconciled = false
	}
REC:
	for _, rule := range rules {
		switch rule.Type {
		case entity.WWNCustomerMapRule:
			if !dupReconciled && entry.WWN == rule.Regex {
				entry.ReconcileRules = append(entry.ReconcileRules, rule.ID)
				entry.IsPrimaryCustomer = false
				if entry.Customer == rule.Comment {
					entry.IsPrimaryCustomer = true
				} else {
					entry.LoadedHostname = ""
					// need to update loadedReconcile if we flushed the name
					loadedReconciled = true
				}
				dupReconciled = true
			}
		case entity.IgnoreLoaded:
			if !loadedReconciled {
				r, err := regexp.Compile(rule.Regex)
				if err != nil {
					return fmt.Errorf("apply-rec-rules - regex %s won't compile: %w", rule.Regex, err)
				}
				match := r.MatchString(entry.LoadedHostname)
				if match {
					entry.ReconcileRules = append(entry.ReconcileRules, rule.ID)
					entry.IgnoreLoaded = true
					loadedReconciled = true
				}
			}
		}
		if loadedReconciled && dupReconciled {
			break REC
		}
	}
	if len(entry.DuplicateCustomers) > 0 && !dupReconciled ||
		(!loadedReconciled) {
		entry.NeedsReconcile = true
	}

	return nil
}
