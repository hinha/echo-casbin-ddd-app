package handler

import (
	"net/http"
	"strconv"

	"github.com/hinha/echo-casbin-ddd-app/internal/application/dto"
	"github.com/hinha/echo-casbin-ddd-app/internal/application/interfaces"
	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
	"github.com/labstack/echo/v4"
)

// @title API Client API
// @version 1.0
// @description API for managing API clients
// @BasePath /api/clients

// APIClientHandler handles HTTP requests for API clients
type APIClientHandler struct {
	apiClientUseCase interfaces.APIClientUseCase
}

// NewAPIClientHandler creates a new APIClientHandler
func NewAPIClientHandler(apiClientUseCase interfaces.APIClientUseCase) *APIClientHandler {
	return &APIClientHandler{
		apiClientUseCase: apiClientUseCase,
	}
}

// APIClientResponse represents an API client in the response
type APIClientResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	APIKey      string `json:"api_key"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// toAPIClientResponse converts an API client entity to an API client response
func toAPIClientResponse(client *entity.APIClient) *APIClientResponse {
	return &APIClientResponse{
		ID:          client.ID,
		Name:        client.Name,
		Description: client.Description,
		APIKey:      client.APIKey,
		Active:      client.Active,
		CreatedAt:   client.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   client.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CreateRequest represents the request for creating an API client
type CreateAPIClientRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// Create handles creating an API client
// @Summary Create a new API client
// @Description Create a new API client with the provided details
// @Tags api-clients
// @Accept json
// @Produce json
// @Param request body CreateAPIClientRequest true "API Client creation request"
// @Success 201 {object} APIClientResponse "Created API client"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router / [post]
func (h *APIClientHandler) Create(c echo.Context) error {
	var req CreateAPIClientRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	input := dto.CreateAPIClientInput{
		Name:        req.Name,
		Description: req.Description,
	}

	output, err := h.apiClientUseCase.Create(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, toAPIClientResponse(output.APIClient))
}

// GetByID handles getting an API client by ID
// @Summary Get an API client by ID
// @Description Retrieve an API client by its ID
// @Tags api-clients
// @Accept json
// @Produce json
// @Param id path int true "API Client ID"
// @Success 200 {object} APIClientResponse "API client details"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "API client not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /{id} [get]
func (h *APIClientHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid API client ID"})
	}

	input := dto.GetAPIClientByIDInput{
		ID: uint(id),
	}

	client, err := h.apiClientUseCase.GetByID(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if client == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "API client not found"})
	}

	return c.JSON(http.StatusOK, toAPIClientResponse(client))
}

// UpdateRequest represents the request for updating an API client
type UpdateAPIClientRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// Update handles updating an API client
// @Summary Update an API client
// @Description Update an existing API client with the provided details
// @Tags api-clients
// @Accept json
// @Produce json
// @Param id path int true "API Client ID"
// @Param request body UpdateAPIClientRequest true "API Client update request"
// @Success 200 {object} APIClientResponse "Updated API client"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /{id} [put]
func (h *APIClientHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid API client ID"})
	}

	var req UpdateAPIClientRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	input := dto.UpdateAPIClientInput{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
	}

	client, err := h.apiClientUseCase.Update(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toAPIClientResponse(client))
}

// RegenerateAPIKey handles regenerating an API key for an API client
// @Summary Regenerate API key
// @Description Regenerate the API key for an existing API client
// @Tags api-clients
// @Accept json
// @Produce json
// @Param id path int true "API Client ID"
// @Success 200 {object} APIClientResponse "API client with new API key"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /{id}/regenerate-key [post]
func (h *APIClientHandler) RegenerateAPIKey(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid API client ID"})
	}

	input := dto.RegenerateAPIKeyInput{
		ID: uint(id),
	}

	client, err := h.apiClientUseCase.RegenerateAPIKey(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toAPIClientResponse(client))
}

// SetActiveRequest represents the request for setting an API client's active status
type SetAPIClientActiveRequest struct {
	Active bool `json:"active"`
}

// SetActive handles setting an API client's active status
// @Summary Set API client active status
// @Description Set the active status of an existing API client
// @Tags api-clients
// @Accept json
// @Produce json
// @Param id path int true "API Client ID"
// @Param request body SetAPIClientActiveRequest true "Set active status request"
// @Success 200 {object} APIClientResponse "Updated API client"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /{id}/set-active [post]
func (h *APIClientHandler) SetActive(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid API client ID"})
	}

	var req SetAPIClientActiveRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	input := dto.SetAPIClientActiveInput{
		ID:     uint(id),
		Active: req.Active,
	}

	client, err := h.apiClientUseCase.SetActive(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toAPIClientResponse(client))
}

// Delete handles deleting an API client
// @Summary Delete an API client
// @Description Delete an existing API client
// @Tags api-clients
// @Accept json
// @Produce json
// @Param id path int true "API Client ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /{id} [delete]
func (h *APIClientHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid API client ID"})
	}

	input := dto.DeleteAPIClientInput{
		ID: uint(id),
	}

	if err := h.apiClientUseCase.Delete(c.Request().Context(), input); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "API client deleted successfully"})
}

// ListAPIClientsResponse represents the response for listing API clients
type ListAPIClientsResponse struct {
	APIClients []*APIClientResponse `json:"api_clients"`
	TotalCount int64                `json:"total_count"`
}

// List handles listing API clients
// @Summary List API clients
// @Description Get a paginated list of API clients
// @Tags api-clients
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {object} ListAPIClientsResponse "List of API clients"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router / [get]
func (h *APIClientHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}

	input := dto.ListAPIClientsInput{
		Page:  page,
		Limit: limit,
	}

	output, err := h.apiClientUseCase.List(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	clients := make([]*APIClientResponse, len(output.APIClients))
	for i, client := range output.APIClients {
		clients[i] = toAPIClientResponse(client)
	}

	resp := ListAPIClientsResponse{
		APIClients: clients,
		TotalCount: output.TotalCount,
	}

	return c.JSON(http.StatusOK, resp)
}

// RegisterRoutes registers the API client routes
func (h *APIClientHandler) RegisterRoutes(e *echo.Echo, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/clients", middlewares...)

	g.POST("", h.Create)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.POST("/:id/regenerate-key", h.RegenerateAPIKey)
	g.POST("/:id/set-active", h.SetActive)
	g.DELETE("/:id", h.Delete)
	g.GET("", h.List)
}
