package api

import (
	"github.com/revel/revel"
	"encoding/json"
	"github.com/leanote/leanote/app/info"
	"gopkg.in/mgo.v2/bson"
//	. "github.com/leanote/leanote/app/lea"
//	"io/ioutil"
)

// 笔记本API

type ApiNotebook struct {
	ApiBaseContrller
}

// 得到用户的所有笔记本
// info.SubNotebooks
func (c ApiNotebook) GetNotebooks() revel.Result {
	re := info.NewRe()
	notebooks := notebookService.GetNotebooks(c.getUserId())
	if notebooks.Len() > 0 {
		re.Ok = true
		re.List = notebooks
	}
	return c.RenderJson(re)
}

// 删除
func (c ApiNotebook) DeleteNotebook(notebookId string) revel.Result {
	re, msg := notebookService.DeleteNotebook(c.getUserId(), notebookId)
	return c.RenderJson(info.Re{Ok: re, Msg: msg})
}

// 添加notebook
func (c ApiNotebook) AddNotebook(notebookId, title, parentNotebookId string) revel.Result {
	notebook := info.Notebook{NotebookId: bson.ObjectIdHex(notebookId), 
		Title: title,
		Seq: -1,
		UserId: c.GetObjectUserId()}
	if(parentNotebookId != "") {
		notebook.ParentNotebookId = bson.ObjectIdHex(parentNotebookId)
	}
	re := info.NewRe()
	re.Ok = notebookService.AddNotebook(notebook)
	
	if(re.Ok) {
		re.Item = notebook
	}
	return c.RenderJson(re)
}
// 修改标题
func (c ApiNotebook) UpdateNotebookTitle(notebookId, title string) revel.Result {
	re := info.NewRe()
	re.Ok = notebookService.UpdateNotebookTitle(notebookId, c.getUserId(), title)
	return c.RenderJson(re)
}

// 调整notebooks, 可能是排序, 可能是移动到其它笔记本下
type DragNotebooksInfo struct {
	CurNotebookId string
	ParentNotebookId string
	Siblings []string
}
// 传过来的data是JSON.stringfy数据
func (c ApiNotebook) DragNotebooks(data string) revel.Result {
	re := info.NewRe()
	
	info := DragNotebooksInfo{}
	json.Unmarshal([]byte(data), &info)
	re.Ok = notebookService.DragNotebooks(c.getUserId(), info.CurNotebookId, info.ParentNotebookId, info.Siblings)
	return c.RenderJson(re)
}