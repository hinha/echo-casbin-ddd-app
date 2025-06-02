package handler

import (
	"net/http"
	"strconv"

	"github.com/hinha/echo-casbin-ddd-app/internal/application/dto"
	"github.com/hinha/echo-casbin-ddd-app/internal/application/interfaces"
	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
	"github.com/labstack/echo/v4"
)

// @title User API
// @version 1.0
// @description API for managing users
// @BasePath /api/users

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userUseCase interfaces.UserUseCase
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userUseCase interfaces.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// RegisterRequest represents the request for user registration
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required"`
}

// RegisterResponse represents the response for user registration
type RegisterResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

// UserResponse represents a user in the response
type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// toUserResponse converts a user entity to a user response
func toUserResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Active:    user.Active,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration request"
// @Success 201 {object} RegisterResponse "Registered user with token"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	input := dto.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	output, err := h.userUseCase.Register(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	resp := RegisterResponse{
		User:  toUserResponse(output.User),
		Token: output.Token,
	}

	return c.JSON(http.StatusCreated, resp)
}

// LoginRequest represents the request for user login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents the response for user login
type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

// Login handles user login
// @Summary User login
// @Description Authenticate a user and return a token
// @Tags users
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login request"
// @Success 200 {object} LoginResponse "Logged in user with token"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	input := dto.LoginInput{
		Username: req.Username,
		Password: req.Password,
	}

	output, err := h.userUseCase.Login(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	resp := LoginResponse{
		User:  toUserResponse(output.User),
		Token: output.Token,
	}

	return c.JSON(http.StatusOK, resp)
}

// GetUser handles getting a user by ID
// @Summary Get a user by ID
// @Description Retrieve a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse "User details"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	user, err := h.userUseCase.GetUserByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, toUserResponse(user))
}

// UpdateUserRequest represents the request for updating a user
type UpdateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

// UpdateUser handles updating a user
// @Summary Update a user
// @Description Update an existing user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "User update request"
// @Success 200 {object} UserResponse "Updated user"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /{id} [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	input := dto.UpdateUserInput{
		ID:       uint(id),
		Username: req.Username,
		Email:    req.Email,
	}

	user, err := h.userUseCase.UpdateUser(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toUserResponse(user))
}

// ChangePasswordRequest represents the request for changing a user's password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// ChangePassword handles changing a user's password
// @Summary Change user password
// @Description Change the password of an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body ChangePasswordRequest true "Change password request"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /{id}/change-password [post]
func (h *UserHandler) ChangePassword(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var req ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	input := dto.ChangePasswordInput{
		ID:          uint(id),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	if err := h.userUseCase.ChangePassword(c.Request().Context(), input); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password changed successfully"})
}

// ListUsersResponse represents the response for listing users
type ListUsersResponse struct {
	Users      []*UserResponse `json:"users"`
	TotalCount int64           `json:"total_count"`
}

// ListUsers handles listing users
// @Summary List users
// @Description Get a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {object} ListUsersResponse "List of users"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router / [get]
func (h *UserHandler) ListUsers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}

	input := dto.ListUsersInput{
		Page:  page,
		Limit: limit,
	}

	output, err := h.userUseCase.ListUsers(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	users := make([]*UserResponse, len(output.Users))
	for i, user := range output.Users {
		users[i] = toUserResponse(user)
	}

	resp := ListUsersResponse{
		Users:      users,
		TotalCount: output.TotalCount,
	}

	return c.JSON(http.StatusOK, resp)
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(e *echo.Echo, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/v1/users", middlewares...)

	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	g.GET("/:id", h.GetUser)
	g.PUT("/:id", h.UpdateUser)
	g.POST("/:id/change-password", h.ChangePassword)
	g.GET("", h.ListUsers)
}
