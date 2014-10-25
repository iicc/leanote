package api

import (
	"github.com/revel/revel"
//	"encoding/json"
	"gopkg.in/mgo.v2/bson"
//	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/info"
//	"github.com/leanote/leanote/app/types"
//	"io/ioutil"
//	"fmt"
//	"bytes"
//	"os"
)

// 笔记API

type ApiNote struct {
	ApiBaseContrller
}

// 笔记首页, 判断是否已登录
// 已登录, 得到用户基本信息(notebook, shareNotebook), 跳转到index.html中
// 否则, 转向登录页面
func (c ApiNote) Index() revel.Result {
	c.SetLocale()
	
	userInfo := c.GetUserInfo()
	
	userId := userInfo.UserId.Hex()
	
	// 没有登录
	if userId == "" {
		return c.Redirect("/login")
	}
	
	c.RenderArgs["openRegister"] = configService.IsOpenRegister()
	
	// 已登录了, 那么得到所有信息
	notebooks := notebookService.GetNotebooks(userId)
	shareNotebooks, sharedUserInfos := shareService.GetShareNotebooks(userId)
	
	// 还需要按时间排序(DESC)得到notes
	notes := []info.Note{}
	noteContent := info.NoteContent{}
	if len(notebooks) > 0 {
//		_, notes = noteService.ListNotes(c.getUserId(), "", false, c.GetPage(), pageSize, defaultSortField, false, false);
		// 变成最新
		_, notes = noteService.ListNotes(c.getUserId(), "", false, c.GetPage(), 50, defaultSortField, false, false);
		if len(notes) > 0 {
			noteContent = noteService.GetNoteContent(notes[0].NoteId.Hex(), userId)
		}
	}
	// 当然, 还需要得到第一个notes的content
	//...
	
	c.RenderArgs["isAdmin"] = leanoteUserId == userInfo.Username
	c.RenderArgs["userInfo"] = userInfo
	c.RenderArgs["userInfoJson"] = c.Json(userInfo)
	c.RenderArgs["notebooks"] = c.Json(notebooks)
	c.RenderArgs["shareNotebooks"] = c.Json(shareNotebooks)
	c.RenderArgs["sharedUserInfos"] = c.Json(sharedUserInfos)
	
	c.RenderArgs["notes"] = c.Json(notes)
	c.RenderArgs["noteContentJson"] = c.Json(noteContent)
	c.RenderArgs["noteContent"] = noteContent.Content
	
	c.RenderArgs["tagsJson"] = c.Json(tagService.GetTags(c.getUserId()))
	
	if isDev, _ := revel.Config.Bool("mode.dev"); isDev {
		return c.RenderTemplate("note/note-dev.html")
	} else {
		return c.RenderTemplate("note/note.html")
	}
}

// 得到笔记本下的笔记
func (c ApiNote) GetNotes(notebookId string) revel.Result {
	re := info.NewRe()
	_, notes := noteService.ListNotes(c.getUserId(), notebookId, false, c.GetPage(), pageSize, defaultSortField, false, false);
	
	if len(notes) > 0 {
		re.Ok = true
		re.List = notes
	}
	
	return c.RenderJson(re)
}

// 得到trash
func (c ApiNote) GetTrashNotes() revel.Result {
	re := info.NewRe()
	_, notes := noteService.ListNotes(c.getUserId(), "", true, c.GetPage(), pageSize, defaultSortField, false, false);
	
	if len(notes) > 0 {
		re.Ok = true
		re.List = notes
	}
	
	return c.RenderJson(re)
}

// 得到note和内容
func (c ApiNote) GetNoteAndContent(noteId string) revel.Result {
	re := info.NewRe()
	re.Item = noteService.GetNoteAndContent(noteId, c.getUserId())
	return c.RenderJson(re)
}

// 得到内容
func (c ApiNote) GetNoteContent(noteId string) revel.Result {
	re := info.NewRe()
	noteContent := noteService.GetNoteContent(noteId, c.getUserId())
	re.Item = noteContent
	return c.RenderJson(re)
}

// 更新note或content
// 肯定会传userId(谁的), NoteId
// 会传Title, Content, Tags, 一种或几种
type NoteOrContent struct {
	NotebookId string
	NoteId string
	UserId string
	Title string
	Desc string
	ImgSrc string
	Tags []string
	Content string
	Abstract string
	IsNew bool
	IsMarkdown bool
	FromUserId string // 为共享而新建
	IsBlog bool // 是否是blog, 更新note不需要修改, 添加note时才有可能用到, 此时需要判断notebook是否设为Blog
}
// 这里不能用json, 要用post
func (c ApiNote) UpdateNoteOrContent(noteOrContent NoteOrContent) revel.Result {
	re := info.NewRe()
	
	// 新添加note
	if noteOrContent.IsNew {
		userId := c.GetObjectUserId();
		myUserId := userId
		// 为共享新建?
		if noteOrContent.FromUserId != "" {
			userId = bson.ObjectIdHex(noteOrContent.FromUserId)
		}
		
		note := info.Note{UserId: userId, 
			NoteId: bson.ObjectIdHex(noteOrContent.NoteId), 
			NotebookId: bson.ObjectIdHex(noteOrContent.NotebookId), 
			Title: noteOrContent.Title, 
			Tags: noteOrContent.Tags,
			Desc: noteOrContent.Desc,
			ImgSrc: noteOrContent.ImgSrc,
			IsBlog: noteOrContent.IsBlog,
			IsMarkdown: noteOrContent.IsMarkdown,
		};
		noteContent := info.NoteContent{NoteId: note.NoteId, 
			UserId: userId, 
			IsBlog: note.IsBlog,
			Content: noteOrContent.Content, 
			Abstract: noteOrContent.Abstract};
		
		note = noteService.AddNoteAndContent(note, noteContent, myUserId)
		re.Ok = true
		re.Item = note
		return c.RenderJson(re)
	}
	
	noteUpdate := bson.M{}
	needUpdateNote := false
	
	// Desc前台传来
	if c.Has("Desc") {
		needUpdateNote = true
		noteUpdate["Desc"] = noteOrContent.Desc;
	}
	if c.Has("ImgSrc") {
		needUpdateNote = true
		noteUpdate["ImgSrc"] = noteOrContent.ImgSrc;
	}
	if c.Has("Title") {
		needUpdateNote = true
		noteUpdate["Title"] = noteOrContent.Title;
	}
	
	if c.Has("Tags[]") {
		needUpdateNote = true
		noteUpdate["Tags"] = noteOrContent.Tags;
	}
	
	if needUpdateNote { 
		noteService.UpdateNote(noteOrContent.UserId, c.getUserId(), 
			noteOrContent.NoteId, noteUpdate)
	}
	
	//-------------
	
	if c.Has("Content") {
		noteService.UpdateNoteContent(noteOrContent.UserId, c.getUserId(), 
			noteOrContent.NoteId, noteOrContent.Content, noteOrContent.Abstract)
	}
	
	re.Ok = true
	return c.RenderJson(re)
}

// 删除note/ 删除别人共享给我的笔记
// userId 是note.UserId
func (c ApiNote) DeleteNote(noteId, userId string, isShared bool) revel.Result {
	if(!isShared) {
		return c.RenderJson(trashService.DeleteNote(noteId, c.getUserId()));
	}
	
	re := info.NewRe()
	re.Ok = trashService.DeleteSharedNote(noteId, userId, c.getUserId())
	return c.RenderJson(re);
}
// 删除trash
func (c ApiNote) DeleteTrash(noteId string) revel.Result {
	re := info.NewRe()
	re.Ok = trashService.DeleteTrash(noteId, c.getUserId())
	return c.RenderJson(re);
}
// 移动note
func (c ApiNote) MoveNote(noteId, notebookId string) revel.Result {
	re := info.NewRe()
	note := noteService.MoveNote(noteId, notebookId, c.getUserId())
	if note.NoteId != "" {
		re.Ok = true
		re.Item = note
	}
	return c.RenderJson(re);
}
// 复制note
func (c ApiNote) CopyNote(noteId, notebookId string) revel.Result {
	re := info.NewRe()
	note := noteService.CopyNote(noteId, notebookId, c.getUserId())
	if note.NoteId != "" {
		re.Ok = true
		re.Item = note
	}
	return c.RenderJson(re);
}
// 复制别人共享的笔记给我
func (c ApiNote) CopySharedNote(noteId, notebookId, fromUserId string) revel.Result {
	re := info.NewRe()
	note := noteService.CopySharedNote(noteId, notebookId, fromUserId, c.getUserId())
	if note.NoteId != "" {
		re.Ok = true
		re.Item = note
	}
	return c.RenderJson(re);
}

//------------
// search
// 通过title搜索
func (c ApiNote) SearchNote(key string) revel.Result {
	_, notes := noteService.SearchNote(key, c.getUserId(), c.GetPage(), pageSize, "UpdatedTime", false, false)
	re := info.NewRe()
	if len(notes) > 0 {
		re.Ok = true
		re.List = notes
	}
	return c.RenderJson(re)
}
// 通过tags搜索
func (c ApiNote) SearchNoteByTags(tags []string) revel.Result {
	_, notes := noteService.SearchNoteByTags(tags, c.getUserId(), c.GetPage(), pageSize, "UpdatedTime", false)
	re := info.NewRe()
	if len(notes) > 0 {
		re.Ok = true
		re.List = notes
	}
	return c.RenderJson(re)
}

// 得到历史列表
func (c ApiNote) GetHistories(noteId string) revel.Result {
	re := info.NewRe()
	histories := noteContentHistoryService.ListHistories(noteId, c.getUserId())
	if len(histories) > 0 {
		re.Ok = true
		re.List = histories
	}
	return c.RenderJson(re)
}