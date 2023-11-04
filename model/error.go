// Package model contains the models used by httphandle.
package model

import (
	"context"

	"github.com/google/uuid"

	"github.com/MicahParks/httphandle/middleware/ctxkey"
)

// Error is the model for an error.
type Error struct {
	Code            int             `json:"code"`
	Message         string          `json:"message"`
	RequestMetadata RequestMetadata `json:"requestMetadata"`
}

// NewError creates a new error.
func NewError(ctx context.Context, code int, message string) Error {
	return Error{
		Code:    code,
		Message: message,
		RequestMetadata: RequestMetadata{
			UUID: ctx.Value(ctxkey.ReqUUID).(uuid.UUID),
		},
	}
}

// RequestMetadata is the model for request metadata.
type RequestMetadata struct {
	UUID uuid.UUID `json:"uuid"`
}
