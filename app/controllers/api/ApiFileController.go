package api

import (
	"github.com/gemsnote/gemsnote/app/controllers"
	"github.com/revel/revel"
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
	path := fileService.GetFile(c.getUserId(), fileId) // 得到路径
	if path == "" {
		return c.RenderText("")
	}
	return controllers.RenderStoredFile(c.Controller, path, "", revel.Inline)
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
