package storage

// Sattus -
type Status string

const (
	StatusUnknown MsgType = "unknown"
)

// NewTxKind -
func NewStatus(value int64) Status {
	return Status(StatusUnknown)
}
