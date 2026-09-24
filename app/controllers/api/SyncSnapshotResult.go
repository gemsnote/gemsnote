package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/revel/revel"
)

// Only the API2 bulk-content snapshot uses this result. Keep legacy API
// responses unchanged and do not depend on deployment-specific gzip settings.
type syncSnapshotResult struct{ data interface{} }

func (r syncSnapshotResult) Apply(req *revel.Request, resp *revel.Response) {
	body, compressed, err := encodeSyncSnapshot(req.Header.Get("Accept-Encoding"), r.data)
	if err != nil {
		revel.PlaintextErrorResult{Error: err}.Apply(req, resp)
		return
	}
	resp.Out.Header().Set("Vary", "Accept-Encoding")
	resp.Out.Header().Set("Cache-Control", "private, no-store")
	if compressed {
		// Revel's optional CompressFilter skips an already encoded response.
		resp.Out.Header().Set("Content-Encoding", "gzip")
	}
	resp.WriteHeader(http.StatusOK, "application/json; charset=utf-8")
	if _, err := resp.GetWriter().Write(body); err != nil {
		revel.AppLog.Warn("Write sync snapshot failed", "error", err)
	}
}

func encodeSyncSnapshot(accept string, data interface{}) ([]byte, bool, error) {
	body, err := json.Marshal(data)
	if err != nil || !acceptsSnapshotGzip(accept) {
		return body, false, err
	}
	var buffer bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buffer, gzip.BestSpeed)
	if err != nil {
		return nil, false, err
	}
	if _, err := writer.Write(body); err != nil {
		return nil, false, err
	}
	if err := writer.Close(); err != nil {
		return nil, false, err
	}
	return buffer.Bytes(), true, nil
}

func acceptsSnapshotGzip(accept string) bool {
	wildcard := false
	for _, coding := range strings.Split(accept, ",") {
		parts := strings.Split(coding, ";")
		name := strings.ToLower(strings.TrimSpace(parts[0]))
		if name != "gzip" && name != "*" {
			continue
		}
		quality := 1.0
		for _, parameter := range parts[1:] {
			key, value, ok := strings.Cut(strings.TrimSpace(parameter), "=")
			if ok && strings.EqualFold(strings.TrimSpace(key), "q") {
				var err error
				quality, err = strconv.ParseFloat(strings.TrimSpace(value), 64)
				if err != nil || quality < 0 || quality > 1 {
					quality = 0
				}
			}
		}
		if name == "gzip" {
			return quality > 0 // Explicit gzip;q=0 overrides a wildcard.
		}
		wildcard = quality > 0
	}
	return wildcard
}
