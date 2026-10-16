package service

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/gemsnote/gemsnote/app/db"
	"github.com/gemsnote/gemsnote/app/info"
	//	. "github.com/gemsnote/gemsnote/app/lea"
	"gopkg.in/mgo.v2/bson"
	//	"time"
)

// 历史记录
type NoteContentHistoryService struct {
}

// 每个历史记录最大值
var maxSize = 20

// 新建一个note, 不需要添加历史记录
// 添加历史
func (this *NoteContentHistoryService) AddHistory(noteId, userId string, eachHistory info.EachHistory) bool {
	if eachHistory.HistoryId == "" {
		eachHistory.HistoryId = bson.NewObjectId()
	}

	// 先查是否存在历史记录, 没有则添加之
	history := info.NoteContentHistory{}
	db.GetByIdAndUserId(db.NoteContentHistories, noteId, userId, &history)
	if history.NoteId == "" {
		return this.newHistory(noteId, userId, eachHistory)
	} else {
		history = this.ensureHistoryIDs(noteId, userId, history)
		// Keep a note-scoped archive directory in the newest stored snapshot.
		// These internal records are never returned by history APIs. Retain
		// originals until permanent note deletion, even when visible versions
		// are trimmed, so restore/pruning cannot race a physical unlink.
		seen := map[string]bool{}
		for _, a := range eachHistory.AttachmentRecords {
			seen[a.AttachId.Hex()] = true
		}
		for _, previous := range history.Histories {
			for _, a := range previous.AttachmentRecords {
				if a.NoteId.Hex() == noteId && !seen[a.AttachId.Hex()] {
					eachHistory.AttachmentRecords = append(eachHistory.AttachmentRecords, a)
					seen[a.AttachId.Hex()] = true
				}
			}
		}

		// 判断是否超出 maxSize, 如果超出则pop最后一个, 再push之, 不用那么麻烦, 直接update吧, 虽然影响性能
		// TODO
		l := len(history.Histories)
		if l >= maxSize {
			// history.Histories = history.Histories[l-maxSize:] // BUG, 致使都是以前的
			history.Histories = history.Histories[:maxSize-1]
		}
		newHistory := []info.EachHistory{eachHistory}
		newHistory = append(newHistory, history.Histories...) // 在开头加了, 最近的在最前
		history.Histories = newHistory

		// 更新之
		return db.UpdateByIdAndUserId(db.NoteContentHistories, noteId, userId, history)
	}
}

// 新建历史
func (this *NoteContentHistoryService) newHistory(noteId, userId string, eachHistory info.EachHistory) bool {
	history := info.NoteContentHistory{NoteId: bson.ObjectIdHex(noteId),
		UserId:    bson.ObjectIdHex(userId),
		Histories: []info.EachHistory{eachHistory},
	}

	// 保存之
	return db.Insert(db.NoteContentHistories, history)
}

// 列表展示
func (this *NoteContentHistoryService) ListHistories(noteId, userId string) []info.EachHistory {
	histories := info.NoteContentHistory{}
	db.GetByIdAndUserId(db.NoteContentHistories, noteId, userId, &histories)
	return this.ensureHistoryIDs(noteId, userId, histories).Histories
}

// GetHistoryByID returns a version without relying on its position in the
// newest-first array. Empty/invalid IDs are intentionally not resolved.
func (this *NoteContentHistoryService) GetHistoryByID(noteId, userId, historyID string) (info.EachHistory, bool) {
	if !bson.IsObjectIdHex(historyID) {
		return info.EachHistory{}, false
	}
	for _, history := range this.ListHistories(noteId, userId) {
		if history.HistoryId.Hex() == strings.ToLower(historyID) {
			return history, true
		}
	}
	return info.EachHistory{}, false
}

func (this *NoteContentHistoryService) ensureHistoryIDs(noteId, userId string, history info.NoteContentHistory) info.NoteContentHistory {
	for i := range history.Histories {
		if history.Histories[i].HistoryId == "" {
			// Deterministic derivation prevents concurrent first reads from
			// assigning different IDs to the same legacy history entry.
			seed := noteId + "\x00" + userId + "\x00" + history.Histories[i].UpdatedTime.UTC().Format("2006-01-02T15:04:05.999999999Z07:00") + "\x00" + history.Histories[i].UpdatedUserId.Hex() + "\x00" + history.Histories[i].Content
			sum := sha256.Sum256([]byte(seed))
			history.Histories[i].HistoryId = bson.ObjectIdHex(hex.EncodeToString(sum[:12]))
		}
	}
	return history
}
