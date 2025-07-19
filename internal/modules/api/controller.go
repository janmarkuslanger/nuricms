package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/middleware"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/internal/utils"
)

type Controller struct {
	services *service.Set
}

func NewController(services *service.Set) *Controller {
	return &Controller{services: services}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (ct Controller) RegisterRoutes(s *server.Server) {
	s.Handle("GET /api/collections/{alias}/content", ct.listContents,
		middleware.ApikeyAuth(ct.services.Apikey),
	)

	s.Handle("GET /api/content/{id}", ct.findContentById,
		middleware.ApikeyAuth(ct.services.Apikey),
	)

	s.Handle("GET /api/collections/{alias}/content/filter", ct.listContentsByFieldValue,
		middleware.ApikeyAuth(ct.services.Apikey),
	)
}

func (ct Controller) findContentById(ctx server.Context) {
	idStr := ctx.Request.PathValue("id")
	id, _ := utils.StringToUint(idStr)

	data, _ := ct.services.Api.FindContentByID(id)

	writeJSON(ctx.Writer, http.StatusOK, dto.ApiResponse{
		Data:    data,
		Success: true,
		Meta: &dto.MetaData{
			Timestamp: time.Now().UTC(),
		},
	})
}

func (ct Controller) listContents(ctx server.Context) {
	alias := ctx.Request.PathValue("alias")

	page := 1
	if p := ctx.Request.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	const perPage = 100
	offset := (page - 1) * perPage

	data, err := ct.services.Api.FindContentByCollectionAlias(alias, offset, perPage)
	if err != nil {
		writeJSON(ctx.Writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(ctx.Writer, http.StatusOK, dto.ApiResponse{
		Data:    data,
		Success: true,
		Meta: &dto.MetaData{
			Timestamp: time.Now().UTC(),
		},
		Pagination: &dto.Pagination{
			PerPage: perPage,
			Page:    page,
		},
	})
}

func (ct Controller) listContentsByFieldValue(ctx server.Context) {
	req := ctx.Request
	w := ctx.Writer

	alias := req.PathValue("alias")

	fieldAlias := req.URL.Query().Get("field")
	value := req.URL.Query().Get("value")
	if fieldAlias == "" || value == "" {
		writeJSON(w, http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Meta:    &dto.MetaData{Timestamp: time.Now().UTC()},
		})
		return
	}

	page := 1
	if p := req.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	const perPage = 100
	offset := (page - 1) * perPage

	items, _ := ct.services.Api.FindContentByCollectionAndFieldValue(alias, fieldAlias, value, offset, perPage)

	writeJSON(w, http.StatusOK, dto.ApiResponse{
		Data:    items,
		Success: true,
		Meta: &dto.MetaData{
			Timestamp: time.Now().UTC(),
		},
		Pagination: &dto.Pagination{
			PerPage: perPage,
			Page:    page,
		},
	})
}
