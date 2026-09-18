package main

import "testing"

func TestColumnRoundTripNames(t *testing.T) {
	tests := map[string]string{
		"_id": "id", "UserId": "user_id", "MaxImageNums": "max_image_nums", "Desc": "description",
	}
	for mongo, postgres := range tests {
		if got := mongoKeyToColumn(mongo); got != postgres {
			t.Errorf("mongoKeyToColumn(%q) = %q, want %q", mongo, got, postgres)
		}
	}
	if got := columnToMongoKey("to_comment_id"); got != "ToCommentId" {
		t.Fatalf("columnToMongoKey(to_comment_id) = %q", got)
	}
}
