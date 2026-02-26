// Package handler
// Обработка http запросов
package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.StaticFile("/", "./index.html")

	// Редирект
	// Редирект на оригинальный URL; 301
	// URL с таким кодом не существует; 404
	// URL не активен, вышло время или удалён; 410
	// Внутреняя ошибка; 500
	r.GET("/:shortCode", h.GetOriginalURL)

	// Создание URL
	// Успешное создание; 200
	// Данные не прошли валидацию; 400
	// Внутреняя ошибка; 500
	r.POST("/shorten", h.CreateShortURL)

	// Статистика по коду
	// Статистика найдена; 200
	// Ошибка; 500
	r.GET("/stats/:shortCode", h.GetStats)

	// Статистика по коду за день
	// Статистика найдена; 200
	// Ошибка; 500
	r.GET("/stats/:shortCode/:date", h.GetDayStats)

	// Удаление URL
	// Успешное удаление; 200
	// Ошибка; 500
	r.DELETE("/:shortCode", h.DeleteURL)
}
