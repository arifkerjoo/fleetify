package response

import (
	"github.com/gofiber/fiber/v3"
)

type PaginationMeta struct {
	CurrentPage int   `json:"currentPage"`
	PerPage     int   `json:"perPage"`
	TotalPages  int   `json:"totalPages"`
	TotalItems  int64 `json:"totalItems"`
}

type Response struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
}

// Success — untuk response sukses tanpa pagination
func Success(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithPagination — untuk response sukses dengan pagination
func SuccessWithPagination(c fiber.Ctx, message string, data interface{}, pagination PaginationMeta) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// Created mengirimkan response saat data berhasil dibuat
func Created(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Message: message,
		Data:    data,
	})
}

// Error mengirimkan response error umum
func Error(c fiber.Ctx, code int, message string, err interface{}) error {
	return c.Status(code).JSON(Response{
		Message: message,
		Errors:  err,
	})
}

// NotFound khusus untuk data tidak ditemukan
func NotFound(c fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message, nil)
}

func InternalServerError(c fiber.Ctx, message, err string) error {
	return Error(c, fiber.StatusInternalServerError, message, err)
}

// BadRequest khusus untuk request yang salah
func BadRequest(c fiber.Ctx, message string, err interface{}) error {
	return Error(c, fiber.StatusBadRequest, message, err)
}

// InternalError khusus untuk error server
func InternalError(c fiber.Ctx, message string, err interface{}) error {
	return Error(c, fiber.StatusInternalServerError, message, err)
}

// MapSlice mengubah slice tipe A ke slice tipe B menggunakan fungsi mapper.
func MapSlice[A any, B any](input []A, mapper func(A) B) []B {
	result := make([]B, 0, len(input))
	for _, v := range input {
		result = append(result, mapper(v))
	}
	return result
}
