package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/dto"
	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/models"
	"github.com/gin-gonic/gin"
)

type URLService interface {
	GetURL(shortCode string) (*models.URL, error)
	CreateURL(url string, expiresInDays *uint32) (*models.URL, error)
	Delete(shortCode string) error
}

type AnalyticsService interface {
	RegisterClick(click *models.Click)
	GetStats(shortCode string) (*models.Stats, error)
	GetDayStats(shortCode string, date string) (*models.DayStats, error)
}

type Handler struct {
	urlService       URLService
	analyticsService AnalyticsService
}

func NewHandler(urlService URLService, analyticsService AnalyticsService) *Handler {
	return &Handler{
		urlService:       urlService,
		analyticsService: analyticsService,
	}
}

// GET /:shortCode
func (h *Handler) GetOriginalURL(c *gin.Context) {
	shortCode, ok := c.Params.Get("shortCode")
	if !ok {
		err := errors.New("missing shortCode")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := dto.ValidateShortCode(shortCode); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	u, err := h.urlService.GetURL(shortCode)
	if err != nil {
		if err.Error() == "url not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "URL not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		return
	}
	if u.ExpiresAt != nil && u.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusGone, dto.ErrorResponse{Error: "URL expired"})
		return
	}
	if !u.IsActive {
		c.JSON(http.StatusGone, dto.ErrorResponse{Error: "URL deleted"})
		return
	}

	// Регистрируем клик асинхронно
	go func() {
		click := &models.Click{
			ShortCode: u.ShortCode,
			IPAdress:  c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Referer:   c.Request.Referer(),
			ClickedAt: time.Now(),
		}

		h.analyticsService.RegisterClick(click)
	}()

	c.Header("Location", u.URL)
	c.Status(http.StatusMovedPermanently)
}

// POST /shorten
func (h *Handler) CreateShortURL(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if len(body) <= 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "empty body"})
		return
	}

	req := &dto.CreateURLRequest{}
	err = json.Unmarshal(body, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	u, err := h.urlService.CreateURL(req.URL, req.ExpiresInDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		return
	}

	resp := &dto.CreateURLResponse{
		ShortCode: u.ShortCode,
		ExpiresAt: u.ExpiresAt,
	}

	c.JSON(http.StatusOK, resp)
}

// GET /stats/:shortCode
func (h *Handler) GetStats(c *gin.Context) {
	shortCode, ok := c.Params.Get("shortCode")
	if !ok {
		err := errors.New("missing shortCode")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := dto.ValidateShortCode(shortCode); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	stats, err := h.analyticsService.GetStats(shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		return
	}

	origURL, err := h.urlService.GetURL(shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		return
	}

	resp := &dto.GetStatsResponse{
		ShortCode:      stats.ShortCode,
		URL:            origURL.URL,
		TotalClicks:    stats.TotalClicks,
		UniqueVisitors: stats.UniqueVisitors,
		LastClickedAt:  stats.LastClickedAt,
		Referers:       stats.Referers,
		IsActive:       origURL.IsActive,
		ExpiresAt:      origURL.ExpiresAt,
	}

	c.JSON(http.StatusOK, resp)
}

// DELETE /:shortCode
func (h *Handler) DeleteURL(c *gin.Context) {
	shortCode, ok := c.Params.Get("shortCode")
	if !ok {
		err := errors.New("missing shortCode")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := dto.ValidateShortCode(shortCode); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	err := h.urlService.Delete(shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		return
	}

	c.Status(http.StatusOK)
}

// GET /stats/:shortCode/:date
func (h *Handler) GetDayStats(c *gin.Context) {
	shortCode, ok := c.Params.Get("shortCode")
	if !ok {
		err := errors.New("missing shortCode")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := dto.ValidateShortCode(shortCode); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	date, ok := c.Params.Get("date")
	if !ok {
		err := errors.New("missing date")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	stats, err := h.analyticsService.GetDayStats(shortCode, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
