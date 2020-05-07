package work

import "github.com/google/uuid"

// Work presents for a work need to be done for the client
type Work struct {
	ID   uuid.UUID   // Unique UUID for each data
	Data interface{} // Data need to be streamed to client
}
