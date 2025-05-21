package order

import (
	"encoding/json"
	"fmt"
)

type OrderStatus int

const (
	StatusCreated OrderStatus = iota
	StatusProcessing
	StatusCompleted
	StatusCancelled
)

func (s OrderStatus) String() string {
	switch s {
	case StatusCreated:
		return "created"
	case StatusProcessing:
		return "processing"
	case StatusCompleted:
		return "completed"
	case StatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// custom json marshalling to prevent sending invalid status
func (s *OrderStatus) UnmarshalJSON(data []byte) error {
	var statusString string
	if err := json.Unmarshal(data, &statusString); err != nil {
		return err
	}

	switch statusString {
	case "created":
		*s = StatusCreated
	case "processing":
		*s = StatusProcessing
	case "completed":
		*s = StatusCompleted
	case "cancelled":
		*s = StatusCancelled
	default:
		return fmt.Errorf("invalid status: %s", statusString)
	}
	return nil
}
