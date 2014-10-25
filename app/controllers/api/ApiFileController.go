package api

import (
	"github.com/revel/revel"
//	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	. "github.com/leanote/leanote/app/lea"
//	"github.com/leanote/leanote/app/lea/netutil"
	"github.com/leanote/leanote/app/info"
	"io/ioutil"
	"os"
//	"strconv"
	"strings"
)

// 文件操作, 图片, 头像上传, 输出

type ApiFile struct {
	ApiBaseContrller
}

// 头像设置
func (c ApiFile) UploadAvatar() revel.Result {
	re := c.uploadImage("logo", "");

	if re.Ok {
		re.Ok = userService.UpdateAvatar(c.getUserId(), re.Id)
		if re.Ok {
			c.UpdateSession("Logo", re.Id);
		}
	}
	
	return c.RenderJson(re)
}

// 上传黑乎乎
func (c ApiFile) UploadImage(albumId string) revel.Result {
	re := c.uploadImage("", albumId);
	return c.RenderJson(re)
}

// 上传图片, 公用方法
// upload image common func
func (c ApiFile) uploadImage(from, albumId string) (re info.Re) {
	var fileUrlPath = ""
	var fileId = ""
	var resultCode = 0 // 1表示正常
	var resultMsg = "内部错误" // 错误信息
	var Ok = false
	
	defer func() {
		re.Id = fileId // 只是id, 没有其它信息
		re.Code = resultCode
		re.Msg = resultMsg
		re.Ok = Ok
	}()
	
	file, handel, err := c.Request.FormFile("file")
	if err != nil {
		return re
	}
	defer file.Close()
	// 生成上传路径
	if(from == "logo") {
		fileUrlPath = "public/upload/" + c.getUserId() + "/images/logo"
	} else {
		fileUrlPath = "files/" + c.getUserId() + "/images"
	}
	dir := revel.BasePath + "/" +  fileUrlPath
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return re
	}
	// 生成新的文件名
	filename := handel.Filename
	
	var ext string;
	if from == "pasteImage" {
		ext = ".png"; // TODO 可能不是png类型
	} else {
		_, ext = SplitFilename(filename)
		if(ext != ".gif" && ext != ".jpg" && ext != ".png" && ext != ".bmp" && ext != ".jpeg") {
			resultMsg = "不是图片"
			return re
		}
	}

	filename = NewGuid() + ext
	data, err := ioutil.ReadAll(file)
	if err != nil {
		LogJ(err)
		return re
	}
	
	// > 2M?
	if(len(data) > 5 * 1024 * 1024) {
		resultCode = 0
		resultMsg = "图片大于2M"
		return re
	}
	
	toPath := dir + "/" + filename;
	err = ioutil.WriteFile(toPath, data, 0777)
	if err != nil {
		LogJ(err)
		return re
	}
	// 改变成gif图片
	_, toPathGif := TransToGif(toPath, 0, true)
	filename = GetFilename(toPathGif)
	filesize := GetFilesize(toPathGif)
	fileUrlPath += "/" + filename
	resultCode = 1
	resultMsg = "上传成功!"
	
	// File
	fileInfo := info.File{Name: filename,
		Title: handel.Filename,
		Path: fileUrlPath,
		Size: filesize}
		
	id := bson.NewObjectId();
	fileInfo.FileId = id
	fileId = id.Hex()
	if(from == "logo") {
		fileId = "public/upload/" + c.getUserId() + "/images/logo/" + filename
	}
	
	Ok = fileService.AddImage(fileInfo, albumId, c.getUserId())
	
	fileInfo.Path = ""; // 不要返回
	re.Item = fileInfo
	
	return re
}

// get all images by userId with page
func (c ApiFile) GetImages(albumId, key string, page int) revel.Result {
	imagesPage := fileService.ListImagesWithPage(c.getUserId(), albumId, key, page, 12)
	re := info.NewRe()
	re.Ok = true
	re.Item = imagesPage
	return c.RenderJson(re)
}

func (c ApiFile) UpdateImageTitle(fileId, title string) revel.Result {
	re := info.NewRe()
	re.Ok = fileService.UpdateImageTitle(c.getUserId(), fileId, title)
	return c.RenderJson(re)
}

func (c ApiFile) DeleteImage(fileId string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = fileService.DeleteImage(c.getUserId(), fileId)
	return c.RenderJson(re)
}

//-----------

// 输出image
// 权限判断
func (c ApiFile) OutputImage(fileId string) revel.Result {
	path := fileService.GetFile(c.getUserId(), fileId); // 得到路径
	if path == "" {
		return c.RenderText("")
	}
	fn := revel.BasePath + "/" +  strings.TrimLeft(path, "/")
    file, _ := os.Open(fn)
    return c.RenderFile(file, revel.Inline) // revel.Attachment
}

// 协作时复制图片到owner
func (c ApiFile) CopyImage(userId, fileId, toUserId string) revel.Result {
	re := info.NewRe()
	
	re.Ok, re.Id = fileService.CopyImage(userId, fileId, toUserId)
	
	return c.RenderJson(re)
}