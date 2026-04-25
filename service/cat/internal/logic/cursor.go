package logic

import (
	"encoding/base64"
	"encoding/json"
)

// pageCursor stores pagination state for cursor-based pagination.
// last_id is the ID of the last item in the previous page.
// Additional sort fields can be added here as needed.
type pageCursor struct {
	LastID uint64 `json:"last_id"`
}

func encodeCursor(lastID uint64) string {
	c := pageCursor{LastID: lastID}
	b, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(b)
}

func decodeCursor(s string) uint64 {
	if s == "" {
		return 0
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return 0
	}
	var c pageCursor
	if err := json.Unmarshal(b, &c); err != nil {
		return 0
	}
	return c.LastID
}
