// Package dto
// Структуры для http и валидация
package dto

type ErrorResponse struct {
	Error string `json:"error"`
}
