package handlers

import (
	"backend/internal/interfaces/middleware"
	"backend/internal/usecase"
	"backend/utils"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
	jwtUtil     *utils.JWTUtil
}

func NewUserHandler(userUsecase usecase.UserUsecase, jwtUtil *utils.JWTUtil) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		jwtUtil:     jwtUtil,
	}
}

func (h *UserHandler) RegisterRoutes(api fiber.Router) {
	users := api.Group("/users")
	users.Use(middleware.JWTMiddleware(h.jwtUtil))

	users.Get("", h.GetAllUsers)
	users.Get("/:id", h.GetUserByID)
	users.Put("/:id", h.UpdateUser)
	users.Post("/:id/profile-image", h.UploadProfileImage)
	users.Delete("/:id", h.DeleteUser)
	users.Post("/change-password", h.ChangePassword)
}

// GetAllUsers retrieves all users with pagination and search
func (h *UserHandler) GetAllUsers(c fiber.Ctx) error {

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	search := c.Query("search", "")

	if limit > 100 {
		limit = 100
	}
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	users, total, err := h.userUsecase.GetAllUsers(limit, offset, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	return c.JSON(fiber.Map{
		"success": true,
		"data":    users,
		"pagination": fiber.Map{
			"total":        total,
			"page":         page,
			"limit":        limit,
			"total_pages":  totalPages,
			"has_next":     page < totalPages,
			"has_previous": page > 1,
		},
	})
}

// GetUserByID retrieves a single user by ID
func (h *UserHandler) GetUserByID(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	user, err := h.userUsecase.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// UpdateUser updates user information
func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	var req usecase.UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.userUsecase.UpdateUser(userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User updated successfully",
		"data":    user,
	})
}

// UploadProfileImage handles profile image upload
func (h *UserHandler) UploadProfileImage(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Verify user can only upload their own image or is admin
	requestUserID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	role := middleware.GetRoleFromContext(c)
	if userID != requestUserID && role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only upload your own profile image",
		})
	}

	// Parse multipart form
	file, err := c.FormFile("profile_image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Profile image is required",
		})
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size must not exceed 5MB",
		})
	}

	// Validate file type
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedTypes[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only image files (jpg, jpeg, png, gif, webp) are allowed",
		})
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d%s", userID.String(), timestamp, ext)
	uploadPath := filepath.Join("uploads", "users", filename)

	// Save file
	if err := c.SaveFile(file, uploadPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save image",
		})
	}

	// Update user profile image in database
	relativePath := filepath.Join("users", filename)
	if err := h.userUsecase.UpdateProfileImage(userID, relativePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	imageURL := fmt.Sprintf("/uploads/%s", relativePath)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile image uploaded successfully",
		"data": fiber.Map{
			"image_url": imageURL,
		},
	})
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	requestUserID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if user is admin
	role := middleware.GetRoleFromContext(c)
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only admins can delete users",
		})
	}

	if err := h.userUsecase.DeleteUser(userID, requestUserID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User deleted successfully",
	})
}

// ChangePassword changes user password
func (h *UserHandler) ChangePassword(c fiber.Ctx) error {
	type ChangePasswordRequest struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8"`
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var req ChangePasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.userUsecase.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password changed successfully",
	})
}
