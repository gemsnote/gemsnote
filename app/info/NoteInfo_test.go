package info

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// 契约测试: 移动端客户端依赖以下 JSON 字段名, 变更需同步客户端
func TestHistoryMetaJSONContract(t *testing.T) {
	meta := HistoryMeta{Index: 3, UpdatedUserId: "550c0bee2ec82a2eb5000000"}
	data, err := json.Marshal(meta)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"Index":3`, `"UpdatedUserId":"550c0bee2ec82a2eb5000000"`, `"UpdatedTime"`} {
		if !strings.Contains(string(data), key) {
			t.Errorf("HistoryMeta JSON missing %s: %s", key, data)
		}
	}
}

func TestEachHistoryJSONContract(t *testing.T) {
	each := EachHistory{Content: "# hello"}
	data, err := json.Marshal(each)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"Content":"# hello"`, `"UpdatedTime"`, `"UpdatedUserId"`} {
		if !strings.Contains(string(data), key) {
			t.Errorf("EachHistory JSON missing %s: %s", key, data)
		}
	}
}

func TestHistoryMetaZeroTimeRoundTrip(t *testing.T) {
	meta := HistoryMeta{Index: 0, UpdatedTime: time.Unix(1480012191, 0).UTC()}
	data, err := json.Marshal(meta)
	if err != nil {
		t.Fatal(err)
	}
	var back HistoryMeta
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatal(err)
	}
	if back.Index != 0 || !back.UpdatedTime.Equal(meta.UpdatedTime) {
		t.Errorf("round trip mismatch: %+v", back)
	}
}
