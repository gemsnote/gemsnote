package api

import (
	"github.com/gemsnote/gemsnote/app/controllers"
	"github.com/gemsnote/gemsnote/app/info"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	//	"encoding/json"
	//	. "github.com/gemsnote/gemsnote/app/lea"
	//	"gopkg.in/mgo.v2/bson"
	//	"github.com/gemsnote/gemsnote/app/lea/netutil"
	//	"github.com/gemsnote/gemsnote/app/info"
	//	"io/ioutil"
	//	"strconv"
)

// 文件操作, 图片, 头像上传, 输出

type ApiFile struct {
	ApiBaseContrller
}

/*
// 协作时复制图片到owner
func (c ApiFile) CopyImage(userId, fileId, toUserId string) revel.Result {
	re := info.NewRe()

	re.Ok, re.Id = fileService.CopyImage(userId, fileId, toUserId)

	return c.RenderJSON(re)
}

// get all images by userId with page
func (c ApiFile) GetImages(albumId, key string, page int) revel.Result {
	imagesPage := fileService.ListImagesWithPage(c.getUserId(), albumId, key, page, 12)
	re := info.NewRe()
	re.Ok = true
	re.Item = imagesPage
	return c.RenderJSON(re)
}

func (c ApiFile) UpdateImageTitle(fileId, title string) revel.Result {
	re := info.NewRe()
	re.Ok = fileService.UpdateImageTitle(c.getUserId(), fileId, title)
	return c.RenderJSON(re)
}

func (c ApiFile) DeleteImage(fileId string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = fileService.DeleteImage(c.getUserId(), fileId)
	return c.RenderJSON(re)
}

*/

//-----------

// 输出image
// [OK]
func (c ApiFile) GetImage(fileId string) revel.Result {
	userId := c.getUserId()
	if userId == "" {
		// Browser <img> requests cannot send an Authorization header, so the
		// URL may carry a short-lived, file-scoped capability instead.
		userId = controllers.ImageCapabilityUser(c.Controller, fileId)
	}
	path := fileService.GetFile(userId, fileId) // 得到路径
	if path == "" {
		return c.RenderText("")
	}
	return controllers.RenderStoredFile(c.Controller, path, "", revel.Inline)
}

// ImageToken issues a short-lived, file-scoped URL so browser <img> elements
// can load an authorised image without an Authorization header.
func (c ApiFile) ImageToken(fileId string) revel.Result {
	re := info.NewRe()
	userId := c.getUserId()
	if userId == "" || !bson.IsObjectIdHex(fileId) || fileService.GetFile(userId, fileId) == "" {
		re.Msg = "noPermission"
		return c.RenderJSON(re)
	}
	url := controllers.SignedImagePath(fileId, userId)
	if url == "" {
		re.Msg = "noPermission"
		return c.RenderJSON(re)
	}
	re.Ok = true
	re.Item = map[string]interface{}{"Url": url, "ExpiresIn": int(controllers.ImageTokenTTL().Seconds())}
	return c.RenderJSON(re)
}

// 下载附件
// [OK]
func (c ApiFile) GetAttach(fileId string) revel.Result {
	attach := attachService.GetAttach(fileId, c.getUserId()) // 得到路径
	path := attach.Path
	if path == "" {
		return c.RenderText("No Such File")
	}
	return controllers.RenderStoredFile(c.Controller, path, attach.Title, revel.Attachment)
}

// 下载所有附件
// [OK]
func (c ApiFile) GetAllAttachs(noteId string) revel.Result {
	note := noteService.GetNoteById(noteId)
	if note.NoteId == "" {
		return c.RenderText("")
	}
	// 得到文件列表
	attachs := attachService.ListAttachs(noteId, c.getUserId())
	if attachs == nil || len(attachs) == 0 {
		return c.RenderText("")
	}

	return controllers.RenderAttachmentArchive(c.Controller, note.Title, attachs)
}
