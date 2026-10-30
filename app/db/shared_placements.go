package db

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// A placement belongs to the reader, never to the source note. It grants no
// access: callers must recheck the share and the destination notebook on read.
func SharedPlacements(userID string) (map[string]string, error) {
	result := map[string]string{}
	if !bson.IsObjectIdHex(userID) {
		return nil, fmt.Errorf("invalid user")
	}
	switch database := globalDatabase.(type) {
	case *PostgresDatabase:
		rows, err := database.db.Query(`SELECT note_id, notebook_id FROM shared_note_placements WHERE user_id=$1`, userID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var noteID, notebookID string
			if err := rows.Scan(&noteID, &notebookID); err != nil {
				return nil, err
			}
			result[noteID] = notebookID
		}
		return result, rows.Err()
	case *MongoDatabase:
		var rows []struct {
			NoteID     string `bson:"NoteId"`
			NotebookID string `bson:"NotebookId"`
		}
		if err := database.db.C("shared_note_placements").Find(bson.M{"UserId": userID}).All(&rows); err != nil {
			return nil, err
		}
		for _, row := range rows {
			result[row.NoteID] = row.NotebookID
		}
		return result, nil
	default:
		return nil, fmt.Errorf("database unavailable")
	}
}

func SetSharedPlacement(userID, noteID, notebookID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(noteID) || !bson.IsObjectIdHex(notebookID) {
		return fmt.Errorf("invalid placement")
	}
	switch database := globalDatabase.(type) {
	case *PostgresDatabase:
		_, err := database.db.Exec(`INSERT INTO shared_note_placements(user_id,note_id,notebook_id) VALUES($1,$2,$3)
			ON CONFLICT(user_id,note_id) DO UPDATE SET notebook_id=excluded.notebook_id`, userID, noteID, notebookID)
		return err
	case *MongoDatabase:
		_, err := database.db.C("shared_note_placements").UpsertId(userID+":"+noteID, bson.M{"UserId": userID, "NoteId": noteID, "NotebookId": notebookID})
		return err
	default:
		return fmt.Errorf("database unavailable")
	}
}
