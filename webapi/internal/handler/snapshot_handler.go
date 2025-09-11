package handler

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ttrnecka/wwn_identity/webapi/internal/mapper"
	"github.com/ttrnecka/wwn_identity/webapi/internal/service"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
)

type SnapshotHandler struct {
	service      service.SnapshotService
	entryService service.FCWWNEntryService
}

func NewSnapshotHandler(s service.SnapshotService, e service.FCWWNEntryService) *SnapshotHandler {
	return &SnapshotHandler{s, e}
}

func (h *SnapshotHandler) Snapshots(c echo.Context) error {

	items, err := h.service.All(c.Request().Context())
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get snapshots", err)
	}
	itemsDTO := []dto.SnapshotDTO{}
	for _, item := range items {
		itemsDTO = append(itemsDTO, mapper.ToSnapshotDTO(item))
	}
	return c.JSON(http.StatusOK, itemsDTO)
}

func (h *SnapshotHandler) CreateSnapshot(c echo.Context) error {
	entries, err := h.entryService.Find(c.Request().Context(), service.Filter{"type": "Unknown"}, service.SortOption{})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to check unknown type entries", err)
	}
	if len(entries) > 0 {
		return errorWithInternal(http.StatusUnprocessableEntity, "Cannot make snapshot, Unknown type entries exist", err)
	}

	// TODO: uncomment once the snapshot work is done
	entries, err = h.entryService.Find(c.Request().Context(), service.Filter{"needs_reconcile": true}, service.SortOption{})
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to check entries requiring reconciliation", err)
	}
	if len(entries) > 0 {
		return errorWithInternal(http.StatusUnprocessableEntity, "Cannot make snapshot, reconcile all entries first", err)
	}

	var snapshotDTO dto.SnapshotDTO
	if err := json.NewDecoder(c.Request().Body).Decode(&snapshotDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := validate.Struct(snapshotDTO); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	snapshot, err := h.service.MakeSnapshot(c.Request().Context(), snapshotDTO.Comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
		// return errorWithInternal(http.StatusInternalServerError, "Failed to create snapshot", err)
	}
	itemDTO := mapper.ToSnapshotDTO(*snapshot)

	return c.JSON(http.StatusOK, itemDTO)
}

func (h *SnapshotHandler) GetSnapshotEntries(c echo.Context) error {
	snapshot_id := c.Param("id")
	snapshot, err := h.service.Get(c.Request().Context(), snapshot_id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	entries, err := h.service.GetSnapshotEntries(c.Request().Context(), *snapshot)
	if err != nil {
		return errorWithInternal(http.StatusInternalServerError, "Failed to get snapshots entries", err)
	}
	itemsDTO := []dto.FCWWNEntryDTO{}
	for _, item := range entries {
		itemsDTO = append(itemsDTO, mapper.ToFCWWNEntryDTO(item))
	}
	return c.JSON(http.StatusOK, itemsDTO)
}
