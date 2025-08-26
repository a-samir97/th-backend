package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError_Error(t *testing.T) {
	// Given
	err := ValidationError{
		Field:   "title",
		Message: "is required",
	}

	// When
	result := err.Error()

	// Then
	expected := "validation error on field 'title': is required"
	assert.Equal(t, expected, result)
}

func TestValidationErrors_Error(t *testing.T) {
	tests := []struct {
		name     string
		errors   ValidationErrors
		expected string
	}{
		{
			name:     "empty errors",
			errors:   ValidationErrors{},
			expected: "validation errors",
		},
		{
			name: "single error",
			errors: ValidationErrors{
				{Field: "title", Message: "is required"},
			},
			expected: "validation errors: validation error on field 'title': is required",
		},
		{
			name: "multiple errors",
			errors: ValidationErrors{
				{Field: "title", Message: "is required"},
				{Field: "file_size", Message: "must be positive"},
			},
			expected: "validation errors: validation error on field 'title': is required, validation error on field 'file_size': must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errors.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidationErrors_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		errors   ValidationErrors
		expected bool
	}{
		{
			name:     "empty errors",
			errors:   ValidationErrors{},
			expected: false,
		},
		{
			name: "with errors",
			errors: ValidationErrors{
				{Field: "title", Message: "is required"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errors.HasErrors()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidationErrors_Add(t *testing.T) {
	// Given
	var errors ValidationErrors

	// When
	errors.Add("title", "is required")
	errors.Add("file_size", "must be positive")

	// Then
	assert.True(t, errors.HasErrors())
	assert.Len(t, errors, 2)
	assert.Equal(t, "title", errors[0].Field)
	assert.Equal(t, "is required", errors[0].Message)
	assert.Equal(t, "file_size", errors[1].Field)
	assert.Equal(t, "must be positive", errors[1].Message)
}

func TestBusinessError_Error(t *testing.T) {
	tests := []struct {
		name     string
		error    BusinessError
		expected string
	}{
		{
			name: "error without details",
			error: BusinessError{
				Code:    "INVALID_REQUEST",
				Message: "Request validation failed",
			},
			expected: "INVALID_REQUEST: Request validation failed",
		},
		{
			name: "error with details",
			error: BusinessError{
				Code:    "INVALID_REQUEST",
				Message: "Request validation failed",
				Details: "Title field is missing",
			},
			expected: "INVALID_REQUEST: Request validation failed (Title field is missing)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.error.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewBusinessError(t *testing.T) {
	// When
	err := NewBusinessError("TEST_ERROR", "This is a test error")

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, "TEST_ERROR", err.Code)
	assert.Equal(t, "This is a test error", err.Message)
	assert.Empty(t, err.Details)
}

func TestNewBusinessErrorWithDetails(t *testing.T) {
	// When
	err := NewBusinessErrorWithDetails("TEST_ERROR", "This is a test error", "Additional details")

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, "TEST_ERROR", err.Code)
	assert.Equal(t, "This is a test error", err.Message)
	assert.Equal(t, "Additional details", err.Details)
}

func TestDomainErrors(t *testing.T) {
	// Test that domain errors are defined
	assert.NotNil(t, ErrMediaNotFound)
	assert.NotNil(t, ErrInvalidMediaStatus)
	assert.NotNil(t, ErrInvalidRequest)
	assert.NotNil(t, ErrUnauthorized)
	assert.NotNil(t, ErrForbidden)
	assert.NotNil(t, ErrInternalError)
	assert.NotNil(t, ErrServiceUnavailable)

	// Test error messages
	assert.Equal(t, "media not found", ErrMediaNotFound.Error())
	assert.Equal(t, "invalid media status", ErrInvalidMediaStatus.Error())
	assert.Equal(t, "invalid request", ErrInvalidRequest.Error())
	assert.Equal(t, "unauthorized", ErrUnauthorized.Error())
	assert.Equal(t, "forbidden", ErrForbidden.Error())
	assert.Equal(t, "internal server error", ErrInternalError.Error())
	assert.Equal(t, "service unavailable", ErrServiceUnavailable.Error())
}

func TestBusinessError_ImplementsError(t *testing.T) {
	// Test that BusinessError implements the error interface
	var err error = &BusinessError{
		Code:    "TEST_ERROR",
		Message: "Test message",
	}

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "TEST_ERROR")
}

func TestValidationError_ImplementsError(t *testing.T) {
	// Test that ValidationError implements the error interface
	var err error = ValidationError{
		Field:   "test_field",
		Message: "test message",
	}

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "test_field")
}

func TestValidationErrors_ImplementsError(t *testing.T) {
	// Test that ValidationErrors implements the error interface
	var err error = ValidationErrors{
		{Field: "field1", Message: "message1"},
	}

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "field1")
}

func TestErrorIsComparable(t *testing.T) {
	// Test that domain errors can be compared using errors.Is
	err1 := ErrMediaNotFound
	err2 := ErrMediaNotFound

	assert.True(t, errors.Is(err1, err2))
	assert.False(t, errors.Is(err1, ErrInvalidRequest))
}
