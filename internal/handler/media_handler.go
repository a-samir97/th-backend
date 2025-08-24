package handler

import (
	"net/http"
	"strconv"

	"thamaniyah/internal/domain"
	"thamaniyah/internal/service"

	"github.com/gin-gonic/gin"
)

// MediaHandler handles HTTP requests for media operations
type MediaHandler struct {
	mediaService service.MediaService
}

// NewMediaHandler creates a new media handler
func NewMediaHandler(mediaService service.MediaService) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
	}
}

// CreateUploadURL godoc
// @Summary Generate upload URL
// @Description Generate a presigned URL for media file upload
// @Tags media
// @Accept json
// @Produce json
// @Param request body domain.UploadRequest true "Upload request"
// @Success 200 {object} domain.UploadURL
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/media/upload-url [post]
func (h *MediaHandler) CreateUploadURL(c *gin.Context) {
	var req domain.UploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	uploadURL, err := h.mediaService.CreateUploadURL(c.Request.Context(), &req)
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
			Message: "Failed to create upload URL",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, uploadURL)
}

// ConfirmUpload godoc
// @Summary Confirm file upload
// @Description Confirm that a file has been uploaded successfully
// @Tags media
// @Produce json
// @Param id path string true "Media ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/media/{id}/confirm [post]
func (h *MediaHandler) ConfirmUpload(c *gin.Context) {
	mediaID := c.Param("id")
	if mediaID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "Media ID is required",
		})
		return
	}

	err := h.mediaService.ConfirmUpload(c.Request.Context(), mediaID)
	if err != nil {
		if err == domain.ErrMediaNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "MEDIA_NOT_FOUND",
				Message: "Media not found",
			})
			return
		}
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
			Message: "Failed to confirm upload",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Upload confirmed successfully",
	})
}

// GetMedia godoc
// @Summary Get media by ID
// @Description Retrieve media details by ID
// @Tags media
// @Produce json
// @Param id path string true "Media ID"
// @Success 200 {object} domain.Media
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/media/{id} [get]
func (h *MediaHandler) GetMedia(c *gin.Context) {
	mediaID := c.Param("id")
	if mediaID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "Media ID is required",
		})
		return
	}

	media, err := h.mediaService.GetMedia(c.Request.Context(), mediaID)
	if err != nil {
		if err == domain.ErrMediaNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "MEDIA_NOT_FOUND",
				Message: "Media not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to get media",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, media)
}

// GetAllMedia godoc
// @Summary List all media
// @Description Retrieve all media with pagination
// @Tags media
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} MediaListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/media [get]
func (h *MediaHandler) GetAllMedia(c *gin.Context) {
	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	mediaList, total, err := h.mediaService.GetAllMedia(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to get media list",
			Details: err.Error(),
		})
		return
	}

	response := MediaListResponse{
		Items:  mediaList,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateMedia godoc
// @Summary Update media metadata
// @Description Update media metadata
// @Tags media
// @Accept json
// @Produce json
// @Param id path string true "Media ID"
// @Param request body domain.UpdateMediaRequest true "Update request"
// @Success 200 {object} domain.Media
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/media/{id} [put]
func (h *MediaHandler) UpdateMedia(c *gin.Context) {
	mediaID := c.Param("id")
	if mediaID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "Media ID is required",
		})
		return
	}

	var req domain.UpdateMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	media, err := h.mediaService.UpdateMedia(c.Request.Context(), mediaID, &req)
	if err != nil {
		if err == domain.ErrMediaNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "MEDIA_NOT_FOUND",
				Message: "Media not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to update media",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, media)
}

// DeleteMedia godoc
// @Summary Delete media
// @Description Soft delete a media record
// @Tags media
// @Produce json
// @Param id path string true "Media ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/media/{id} [delete]
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	mediaID := c.Param("id")
	if mediaID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "Media ID is required",
		})
		return
	}

	err := h.mediaService.DeleteMedia(c.Request.Context(), mediaID)
	if err != nil {
		if err == domain.ErrMediaNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "MEDIA_NOT_FOUND",
				Message: "Media not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Failed to delete media",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Media deleted successfully",
	})
}

// Response types

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
}

// MediaListResponse represents a paginated media list response
type MediaListResponse struct {
	Items  []*domain.Media `json:"items"`
	Total  int64           `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}
