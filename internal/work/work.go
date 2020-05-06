package work

import "github.com/google/uuid"

// Work presents for a data need to be streamed to client
type Work struct {
	ID   uuid.UUID   // unique UUID for each data
	Data interface{} // data need to be streamed to client
}
