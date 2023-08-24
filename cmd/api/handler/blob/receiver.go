package blob

import "context"

//go:generate mockgen -source=$GOFILE -destination=mock.go -package=blob -typed
type Receiver interface {
	Blobs(ctx context.Context, height uint64, hash ...string) ([]Blob, error)
}

type Blob struct {
	Namespace    string `json:"namespace"`
	Data         string `json:"data"`
	ShareVersion int    `json:"share_version"`
	Commitment   string `json:"commitment"`
}
