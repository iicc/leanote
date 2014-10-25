package api

import (
	"github.com/revel/revel"
//	"encoding/json"
//	"gopkg.in/mgo.v2/bson"
//	. "github.com/leanote/leanote/app/lea"
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