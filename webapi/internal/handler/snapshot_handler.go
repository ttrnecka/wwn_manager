package handler

import (
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
