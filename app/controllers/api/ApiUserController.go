package api

import (
	"github.com/revel/revel"
//	"encoding/json"
//	"gopkg.in/mgo.v2/bson"
	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/info"
//	"github.com/leanote/leanote/app/types"
//	"io/ioutil"
//	"fmt"
//	"math"
//	"os"

//	"path"
//	"strconv"
)

type ApiUser struct {
	ApiBaseContrller
}

// 获取用户信息
func (c ApiUser) Info() revel.Result {
	re := info.NewRe()
	
	userInfo := c.getUserInfo()
	if userInfo.UserId == "" {
		return c.RenderJson(re)
	}
	re.Ok = true
	re.Item = userInfo
	return c.RenderJson(re)
}

// 修改用户名
func (c ApiUser) UpdateUsername(username string) revel.Result {
	re := info.NewRe();
	if(c.GetUsername() == "demo") {
		re.Msg = "cannotUpdateDemo"
		return c.RenderRe(re);
	}
	
	if re.Ok, re.Msg = Vd("username", username); !re.Ok {
		return c.RenderRe(re);
	}
	
	re.Ok, re.Msg = userService.UpdateUsername(c.GetUserId(), username)
	if(re.Ok) {
		c.UpdateSession("Username", username)
	}
	return c.RenderRe(re);
}

// 修改密码
func (c ApiUser) UpdatePwd(oldPwd, pwd string) revel.Result {
	re := info.NewRe();
	if(c.GetUsername() == "demo") {
		re.Msg = "cannotUpdateDemo"
		return c.RenderRe(re);
	}
	if re.Ok, re.Msg = Vd("password", oldPwd); !re.Ok {
		return c.RenderRe(re);
	}
	if re.Ok, re.Msg = Vd("password", pwd); !re.Ok {
		return c.RenderRe(re);
	}
	re.Ok, re.Msg = userService.UpdatePwd(c.GetUserId(), oldPwd, pwd)
	return c.RenderRe(re);
}