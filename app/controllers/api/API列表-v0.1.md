# API 列表

By life (life@leanote.com)

除了/api/login, /api/register外其它的都需要另外带参数token=xxxx

## Auth 登录与注册

* /api/login 登录
```
参数: email, pwd, captcha(验证码)
返回: 成功返回Ok = true, Item是token, 否则Ok = false, 若Code == S_NEED_CAPTCHA 表示需要验证码
```

* /api/logout 注销
```
无参数
返回: 成功返回Ok = true, 否则为false
```

* /api/register 注册
```
参数: email, pwd
返回: 成功返回 Ok = true, 否则返回 false,  Msg有相应的提示
```

## User 用户
* /api/user/info 获取用户信息
```
无参数
返回: 成功则返回 Ok = true, Item 为info.User, 否则Ok = false
```

* /api/user/updateUsername 修改用户名
```
参数: username
返回: 成功返回 Ok = true, 否则为false, Msg有相应提示
```

* /api/user/updatePwd 修改密码
```
参数: oldPwd, pwd
返回: 成功返回 Ok = true, 否则为false, Msg有相应提示
```

## Notebook 笔记本
* /api/notebook/getNotebooks 得到所有笔记本 
```
无参数
返回: 成功返回 Ok = true, List为info.SubNotebooks(父与子), 否则返回false
```

* /api/notebook/addNotebook 添加笔记本
```
参数: title, parentNotebookId(父notebookId, 可空)
返回: 成功返回 Ok = true, Item为notebook信息, 否则为false
```

* /api/notebook/deleteNotebook 删除笔记本
```
参数: notebookId
返回: 成功返回 Ok = true, 否则为false
```

* /api/notebook/updateNotebookTitle 修改笔记本标题
```
参数: notebookId, title
返回: 成功返回 Ok = true, 否则为false
```

* /api/notebook/dragNotebooks 拖拽笔记本
```
参数: curNotebookId 当前笔记本Id, parentNotebookId 父笔记本Id, siblings []string 所有的兄笔记本Ids
返回: 成功返回 Ok = true, 否则为false
```

## Note 笔记
* /api/note/getNotes 获得某笔记本下的笔记(无内容)
```
参数: notebookId
返回: 成功返回 Ok = true, List为笔记列表(无内容), 否则返回false
```

* /api/note/getTrashNotes 获得所有Trash笔记(无内容)
```
无参数
返回: 成功返回 Ok = true, List为笔记列表(无内容), 否则返回false
```

* /api/note/getNoteContent 获得笔记内容
```
参数: noteId
返回: 成功返回 Ok = true, Item是info.NoteContent, 否则返回false
```

* /api/note/getNoteAndContent 获得笔记及内容
```
参数: noteId
返回: 成功返回 Ok = true, Item是info.NoteAndContent, 否则返回false
```

* /api/note/getHistories 获得某笔记的历史记录
```
参数: noteId
返回: 成功返回 Ok = true, Item为info.NoteContentHistory, 否则返回false
```

* /api/note/addNote 添加笔记
```
参数: 
		NotebookId string 必传
		Title string 可选
		Desc string 可选
		ImgSrc string 可选
		Tags []string 可选
		Content string 可选
		Abstract string 可选
		IsMarkdown bool 可选
		FromUserId string // 为共享而新建 可选
		IsBlog bool // 是否是blog 可选

返回: 成功返回 Ok = true, Item是note; 否则返回false
```

* /api/note/updateNote 更新笔记
```
参数:
		NoteId string 必传
		Title string 可选
		Desc string 可选
		ImgSrc string 可选
		Tags []string 可选
		Content string 可选
		Abstract string 可选
		
返回: 成功返回 Ok = true; 否则返回false
```

* /api/note/deleteNote 删除笔记
```
参数: noteId, userId 用户Id, isShared 是否是删除共享的笔记
返回: 成功返回 Ok = true; 否则返回false
```

* /api/note/deleteTrash 删除Trash笔记(彻底删除)
```
参数: noteId
返回: 成功返回 Ok = true; 否则返回false
```

* /api/note/moveNote 移动笔记
```
参数: noteId, notebookId 将笔记移动到notebookId下
返回: 成功返回 Ok = true; 否则返回false
```

* /api/note/copyNote 复制笔记
```
参数: noteId, notebookId 将笔记复制到notebookId下
返回: 成功返回 Ok = true; 否则返回false
```

* /api/note/copySharedNote 复制被分享的笔记
```
参数: noteId, notebookId, fromUserId 将笔记复制到notebookId下, 该noteId是fromUserId共享的
返回: 成功返回 Ok = true; 否则返回false
```

## File 文件操作(上传图片, 显示图片)

* /api/file/uploadAvatar 上传头像
```
参数: file 图片
返回: 成功返回Ok = true, Id为图片路径, 如upload/xxx/xxx.jpg; 否则返回false
```

* /api/file/uploadImage 上传图片
```
参数: albumId 相册id(可为空), file 图片
返回: 成功返回Ok = true, Id为fileId; 否则返回false
```

* /api/file/outputImage 显示图片
```
参数: fileId
返回: 成功返回图片(二进制); 否则返回
```

## Captcha 验证码
* /api/captcha/get
```
无参数
返回: 图片二进制
```