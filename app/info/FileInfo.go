package info

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type File struct {
	FileId         bson.ObjectId `bson:"_id,omitempty"` //
	UserId         bson.ObjectId `bson:"UserId"`
	AlbumId        bson.ObjectId `bson:"AlbumId"`
	Name           string        `Name`  // file name
	Title          string        `Title` // file name or user defined for search
	Size           int64         `Size`  // file size (byte)
	Type           string        `Type`  // file type, "" = image, "doc" = word
	Path           string        `Path`  // the file path
	IsDefaultAlbum bool          `IsDefaultAlbum`
	CreatedTime    time.Time     `CreatedTime`

	FromFileId bson.ObjectId `bson:"FromFileId,omitempty"` // copy from fileId, for collaboration

	// Only API2 1.1 note uploads participate in orphan cleanup. Empty state
	// means an existing/legacy resource, which must never be cleaned this way.
	UploadState  string        `bson:"UploadState,omitempty" json:"-"`
	UploadNoteId bson.ObjectId `bson:"UploadNoteId,omitempty" json:"-"`
	UploadDigest string        `bson:"UploadDigest,omitempty" json:"-"`
	UploadKey    string        `bson:"UploadKey,omitempty" json:"-"`
}
