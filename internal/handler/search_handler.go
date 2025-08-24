package handler

import (
	"net/http"
	"strconv"

	"thamaniyah/internal/domain"
	"thamaniyah/internal/service"

	"github.com/gin-gonic/gin"
)

// SearchHandler handles HTTP requests for search operations
type SearchHandler struct {
	searchService service.SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchService service.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// Search godoc
// @Summary Search media content
// @Description Search across all media content with filters
// @Tags search
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param type query string false "Media type (video, podcast)"
// @Param limit query int false "Limit results" default(20)
// @Param offset query int false "Offset results" default(0)
// @Success 200 {object} domain.SearchResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	// Parse query parameters
	var req domain.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "Invalid search parameters",
			Details: err.Error(),
		})
		return
	}

	// Parse limit and offset manually since ShouldBindQuery might not handle them correctly
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			req.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			req.Offset = offset
		}
	}

	// Perform search
	response, err := h.searchService.Search(c.Request.Context(), &req)
	if err != nil {
		if businessErr, ok := err.(*domain.BusinessError); ok {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   businessErr.Code,
				Message: businessErr.Message,
				Details: businessErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Search failed",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Suggest godoc
// @Summary Get search suggestions
// @Description Get search suggestions based on partial query
// @Tags search
// @Accept json
// @Produce json
// @Param query query string true "Partial search query"
// @Param limit query int false "Limit suggestions" default(10)
// @Success 200 {object} domain.SuggestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/search/suggest [get]
func (h *SearchHandler) Suggest(c *gin.Context) {
	// Parse query parameters
	var req domain.SuggestRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "Invalid suggest parameters",
			Details: err.Error(),
		})
		return
	}

	// Parse limit manually
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			req.Limit = limit
		}
	}

	// Get suggestions
	response, err := h.searchService.Suggest(c.Request.Context(), &req)
	if err != nil {
		if businessErr, ok := err.(*domain.BusinessError); ok {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   businessErr.Code,
				Message: businessErr.Message,
				Details: businessErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Suggestions failed",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Reindex godoc
// @Summary Reindex search data
// @Description Rebuild the search index with latest data from CMS service
// @Tags search
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/search/reindex [post]
func (h *SearchHandler) Reindex(c *gin.Context) {
	err := h.searchService.Reindex(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Reindex failed",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Search index rebuilt successfully",
	})
}
