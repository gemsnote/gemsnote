package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gemsnote/gemsnote/app/db"
	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
)

type SharingRecipient struct {
	UserId   string
	Username string
	NoteIds  []string
}
type PlacedSharedNote struct {
	info.Note
	IsShared bool
	Perm     int
}
type SharingNavigation struct {
	Owners       []SharingRecipient
	Suggestions  []SharingRecipient
	Recipients   []SharingRecipient
	Placements   map[string]string
	PlacedShares []PlacedSharedNote
}

// Navigation returns only public usernames and notes the requester may see.
// Group and notebook grants use the same effective permission check as reads.
func (s *ShareService) Navigation(userID string) (*SharingNavigation, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, fmt.Errorf("invalid user")
	}
	result := &SharingNavigation{Owners: []SharingRecipient{}, Suggestions: s.shareSuggestions(userID), Recipients: []SharingRecipient{}, Placements: map[string]string{}, PlacedShares: []PlacedSharedNote{}}
	// Keep public sharer names in the snapshot navigation, so offline clients
	// can identify the author without caching private user profiles.
	_, owners := s.GetShareNotebooks(userID)
	for _, owner := range owners {
		result.Owners = append(result.Owners, SharingRecipient{UserId: owner.UserId.Hex(), Username: owner.Username, NoteIds: []string{}})
	}
	placements, err := db.SharedPlacements(userID)
	if err != nil {
		return nil, err
	}
	for _, note := range s.placedNotes(userID, placements) {
		target := placements[note.NoteId.Hex()]
		if !bson.IsObjectIdHex(target) {
			continue
		}
		book := notebookService.GetNotebook(target, userID)
		if book.NotebookId == "" || book.IsDeleted || book.IsTrash {
			continue
		}
		result.Placements[note.NoteId.Hex()] = target
		note.NotebookId = book.NotebookId
		result.PlacedShares = append(result.PlacedShares, PlacedSharedNote{Note: note, IsShared: true})
	}

	var direct []info.ShareNote
	var books []info.ShareNotebook
	q := bson.M{"UserId": bson.ObjectIdHex(userID)}
	db.ListByQ(db.ShareNotes, q, &direct)
	db.ListByQ(db.ShareNotebooks, q, &books)
	// Expand each grant once rather than checking every user against every
	// note. These are the same direct/group/notebook alternatives as getOrQ.
	groupMembers := map[bson.ObjectId][]bson.ObjectId{}
	recipientsFor := func(userID, groupID bson.ObjectId) []bson.ObjectId {
		ids := []bson.ObjectId{}
		if userID.Valid() {
			ids = append(ids, userID)
		}
		if groupID.Valid() {
			members, ok := groupMembers[groupID]
			if !ok {
				var rows []info.GroupUser
				db.ListByQ(db.GroupUsers, bson.M{"GroupId": groupID}, &rows)
				for _, row := range rows {
					members = append(members, row.UserId)
				}
				var group info.Group
				db.Get(db.Groups, groupID.Hex(), &group)
				if group.UserId.Valid() {
					members = append(members, group.UserId)
				}
				groupMembers[groupID] = members
			}
			ids = append(ids, members...)
		}
		return ids
	}
	noteRecipients, bookRecipients := map[bson.ObjectId][]bson.ObjectId{}, map[bson.ObjectId][]bson.ObjectId{}
	noteIDs, bookIDs := []bson.ObjectId{}, []bson.ObjectId{}
	for _, grant := range direct {
		if !grant.NoteId.Valid() {
			continue
		}
		noteIDs = append(noteIDs, grant.NoteId)
		noteRecipients[grant.NoteId] = append(noteRecipients[grant.NoteId], recipientsFor(grant.ToUserId, grant.ToGroupId)...)
	}
	for _, grant := range books {
		if !grant.NotebookId.Valid() {
			continue
		}
		bookIDs = append(bookIDs, grant.NotebookId)
		bookRecipients[grant.NotebookId] = append(bookRecipients[grant.NotebookId], recipientsFor(grant.ToUserId, grant.ToGroupId)...)
	}
	if len(noteIDs) == 0 && len(bookIDs) == 0 {
		return result, nil
	}
	var activeNotes []info.Note
	db.ListByQ(db.Notes, bson.M{"UserId": bson.ObjectIdHex(userID), "IsDeleted": false, "IsTrash": false,
		"$or": []bson.M{{"_id": bson.M{"$in": noteIDs}}, {"NotebookId": bson.M{"$in": bookIDs}}}}, &activeNotes)
	byUser := map[bson.ObjectId]map[string]bool{}
	for _, note := range activeNotes {
		ids := append(append([]bson.ObjectId{}, noteRecipients[note.NoteId]...), bookRecipients[note.NotebookId]...)
		for _, id := range ids {
			if !id.Valid() || id.Hex() == userID {
				continue
			}
			if byUser[id] == nil {
				byUser[id] = map[string]bool{}
			}
			byUser[id][note.NoteId.Hex()] = true
		}
	}
	userIDs := []bson.ObjectId{}
	for id := range byUser {
		userIDs = append(userIDs, id)
	}
	if len(userIDs) > 0 {
		var users []info.User
		db.ListByQWithFields(db.Users, bson.M{"_id": bson.M{"$in": userIDs}, "Flag": bson.M{"$ne": info.FlagDisabled}}, []string{"_id", "Username"}, &users)
		for _, user := range users {
			recipient := SharingRecipient{UserId: user.UserId.Hex(), Username: user.Username, NoteIds: []string{}}
			for id := range byUser[user.UserId] {
				recipient.NoteIds = append(recipient.NoteIds, id)
			}
			sort.Strings(recipient.NoteIds)
			result.Recipients = append(result.Recipients, recipient)
		}
	}
	sort.Slice(result.Recipients, func(i, j int) bool { return result.Recipients[i].Username < result.Recipients[j].Username })
	seen := map[string]bool{}
	for _, u := range result.Suggestions {
		seen[u.UserId] = true
	}
	for _, u := range result.Recipients {
		if !seen[u.UserId] {
			result.Suggestions = append(result.Suggestions, SharingRecipient{UserId: u.UserId, Username: u.Username, NoteIds: []string{}})
		}
	}
	sort.Slice(result.Suggestions, func(i, j int) bool { return result.Suggestions[i].Username < result.Suggestions[j].Username })

	return result, nil
}

// Only load explicitly placed notes; never use a saved placement as a grant.
func (s *ShareService) placedNotes(userID string, placements map[string]string) []info.Note {
	ids := []bson.ObjectId{}
	for id := range placements {
		if bson.IsObjectIdHex(id) {
			ids = append(ids, bson.ObjectIdHex(id))
		}
	}
	if len(ids) == 0 {
		return []info.Note{}
	}
	var notes []info.Note
	db.ListByQ(db.Notes, bson.M{"_id": bson.M{"$in": ids}, "IsDeleted": false, "IsTrash": false}, &notes)
	result := []info.Note{}
	for _, note := range notes {
		if note.UserId.Hex() != userID && s.HasReadPerm(note.UserId.Hex(), userID, note.NoteId.Hex()) {
			result = append(result, note)
		}
	}
	return result
}

func (s *ShareService) PlaceReceivedNote(userID, noteID, notebookID string) error {
	if !bson.IsObjectIdHex(noteID) || !bson.IsObjectIdHex(notebookID) {
		return fmt.Errorf("invalidId")
	}
	note := noteService.GetNoteById(noteID)
	if note.NoteId == "" || note.UserId.Hex() == userID || note.IsDeleted || note.IsTrash || !s.HasReadPerm(note.UserId.Hex(), userID, noteID) {
		return fmt.Errorf("noAuth")
	}
	book := notebookService.GetNotebook(notebookID, userID)
	if book.NotebookId == "" || book.IsDeleted || book.IsTrash {
		return fmt.Errorf("invalidNotebook")
	}
	return db.SetSharedPlacement(userID, note.NoteId.Hex(), book.NotebookId.Hex())
}

// Suggestions are scoped to this owner's existing sharing relationships. They
// never enumerate the user directory or disclose recipient email addresses.
func (s *ShareService) shareSuggestions(ownerID string) []SharingRecipient {
	result := []SharingRecipient{}
	var history []info.HasShareNote
	db.ListByQ(db.HasShareNotes, bson.M{"UserId": bson.ObjectIdHex(ownerID)}, &history)
	seen := map[bson.ObjectId]bool{}
	ids := []bson.ObjectId{}
	for _, row := range history {
		if row.ToUserId.Valid() && row.ToUserId.Hex() != ownerID && !seen[row.ToUserId] {
			seen[row.ToUserId] = true
			ids = append(ids, row.ToUserId)
		}
	}
	if len(ids) == 0 {
		return result
	}
	var users []info.User
	db.ListByQWithFields(db.Users, bson.M{"_id": bson.M{"$in": ids}, "Flag": bson.M{"$ne": info.FlagDisabled}}, []string{"_id", "Username"}, &users)
	for _, user := range users {
		if user.Username != "" {
			result = append(result, SharingRecipient{UserId: user.UserId.Hex(), Username: user.Username, NoteIds: []string{}})
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Username < result[j].Username })
	return result
}

// API2 accepts a username or email in the existing emails array. Legacy actions
// keep their original email-only behavior and response shape.
func (s *ShareService) AddShareNoteByIdentity(resourceID, ownerID, identity string) (bool, string, string) {
	if !bson.IsObjectIdHex(resourceID) || !bson.IsObjectIdHex(ownerID) {
		return false, "noAuth", ""
	}
	note := noteService.GetNote(resourceID, ownerID)
	if note.NoteId == "" || note.IsDeleted || note.IsTrash {
		return false, "noAuth", ""
	}
	identity = strings.TrimSpace(identity)
	if identity == "" {
		return false, "无该用户", ""
	}
	user := userService.GetUserInfoByName(identity)
	if !user.UserId.Valid() || user.Flag == info.FlagDisabled {
		return false, "无该用户", ""
	}
	if user.UserId.Hex() == ownerID {
		return false, "不能分享给自己", ""
	}
	return s.AddShareNoteToUserId(resourceID, 0, ownerID, user.UserId.Hex())
}
