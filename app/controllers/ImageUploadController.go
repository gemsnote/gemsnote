package controllers

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/gemsnote/gemsnote/app/db"
	"github.com/gemsnote/gemsnote/app/info"
	"github.com/gemsnote/gemsnote/app/service"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

func imageUploadResult(file info.File) info.Re {
	return info.Re{Ok: true, Id: file.FileId.Hex(), Item: map[string]interface{}{
		"FileId": file.FileId.Hex(), "NoteId": file.UploadNoteId.Hex(),
		"State": file.UploadState, "Digest": file.UploadDigest,
		"CanCleanup": file.UploadState == "pending" || file.UploadState == "deleting",
	}}
}

func (c File) PasteImageV2(noteId, uploadId string) revel.Result {
	// Existing Web clients may upload before their new note is saved. Such an
	// upload has no receipt and is deliberately outside the cleanup contract.
	if noteId == "" && uploadId == "" {
		return c.PasteImage("")
	}
	if !bson.IsObjectIdHex(noteId) || !shareService.HasUpdateNotePerm(noteId, c.GetUserId()) {
		return c.RenderJSON(info.Re{Msg: "noPermission"})
	}
	var data []byte
	c.Params.Bind(&data, "file")
	if len(data) == 0 {
		return c.RenderJSON(info.Re{Msg: "fileNotFound"})
	}
	digest := sha256.Sum256(data)
	id, err := service.NewImageUploadID()
	if err != nil {
		return c.RenderJSON(info.Re{Msg: "randomSourceUnavailable"})
	}
	file := info.File{FileId: bson.ObjectIdHex(id), UploadState: "pending", UploadNoteId: bson.ObjectIdHex(noteId), UploadDigest: hex.EncodeToString(digest[:])}
	if uploadId != "" {
		key, err := service.ImageUploadKey(c.GetUserId(), uploadId)
		if err != nil {
			return c.RenderJSON(info.Re{Msg: err.Error()})
		}
		file.UploadKey = key
		existing, err := db.ImageUploadByKey(c.GetUserId(), key)
		if err == nil {
			return c.existingImageUpload(file, existing)
		}
		if err.Error() != "uploadNotFound" {
			return c.RenderJSON(info.Re{Msg: err.Error()})
		}
	}
	re := c.uploadImage("pasteImage", "", file)
	if !re.Ok && uploadId != "" {
		// Another process may have committed the same upload key first.
		if existing, err := db.ImageUploadByKey(c.GetUserId(), file.UploadKey); err == nil {
			return c.existingImageUpload(file, existing)
		}
	}
	if re.Ok && uploadId != "" {
		existing, err := db.ImageUploadByKey(c.GetUserId(), file.UploadKey)
		if err != nil {
			return c.RenderJSON(info.Re{Msg: err.Error()})
		}
		return c.existingImageUpload(file, existing)
	}
	return c.RenderJSON(re)
}

func (c File) existingImageUpload(request, existing info.File) revel.Result {
	if existing.UploadNoteId != request.UploadNoteId || existing.UploadDigest != request.UploadDigest {
		return c.RenderJSON(info.Re{Msg: "uploadKeyConflict"})
	}
	if existing.UploadState == "deleted" || existing.UploadState == "deleting" {
		return c.RenderJSON(info.Re{Msg: "uploadAlreadyCleaned"})
	}
	return c.RenderJSON(imageUploadResult(existing))
}

func (c File) ImageUploadStatus(fileId, uploadId string) revel.Result {
	if uploadId != "" {
		key, err := service.ImageUploadKey(c.GetUserId(), uploadId)
		if err != nil {
			return c.RenderJSON(info.Re{Msg: err.Error()})
		}
		file, err := db.ImageUploadByKey(c.GetUserId(), key)
		if err != nil {
			return c.RenderJSON(info.Re{Msg: err.Error()})
		}
		if fileId != "" && fileId != file.FileId.Hex() {
			return c.RenderJSON(info.Re{Msg: "invalidRequest"})
		}
		return c.RenderJSON(imageUploadResult(file))
	}
	file, err := db.ImageUpload(c.GetUserId(), fileId)
	if err != nil {
		return c.RenderJSON(info.Re{Msg: err.Error()})
	}
	return c.RenderJSON(imageUploadResult(file))
}

func (c File) CleanupImageUpload(fileId string) revel.Result {
	if err := service.CleanupNoteImage(revel.BasePath, c.GetUserId(), fileId); err != nil {
		return c.RenderJSON(info.Re{Msg: err.Error()})
	}
	return c.RenderJSON(info.Re{Ok: true})
}
