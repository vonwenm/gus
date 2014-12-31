package storage

import (
	"encoding/json"
)

type HeaderMap map[string]string

func (h *HeaderMap) ToString() string {
	if str, err := json.Marshal(h); err == nil {
		return string(str)
	}
	return ""
}

func NewHeaderMap(source string) HeaderMap {
	var hm HeaderMap
	if source != "" {
		json.Unmarshal([]byte(source), &hm)
	}
	return hm
}
