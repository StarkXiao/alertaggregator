package fingerprint

import (
	"alertaggregator/internal/model"
	"alertaggregator/internal/normalize"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

func Build(e model.Event) string {
	keys := make([]string, 0, len(e.Labels))
	for k := range e.Labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString(e.Service + "|" + e.Environment + "|" + e.Level + "|" + normalize.Message(e.Message))
	for _, k := range keys {
		switch strings.ToLower(k) {
		case "trace_id", "request_id", "span_id":
			continue
		}
		b.WriteString("|" + k + "=" + e.Labels[k])
	}
	h := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(h[:])
}
