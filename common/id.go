package common

import "github.com/google/uuid"

// NewID generates a new UUID
func NewID() uuid.UUID {
    return uuid.New()
}