package order

type Status int

const (
	StatusCreated Status = iota
	StatusProcessing
	StatusCompleted
	StatusCancelled
)

func (s Status) String() string {
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
