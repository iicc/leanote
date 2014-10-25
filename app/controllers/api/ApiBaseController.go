package api

import (
//	"github.com/revel/revel"
//	"gopkg.in/mgo.v2/bson"
//	"encoding/json"
//	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/controllers"
	"github.com/leanote/leanote/app/info"
//	"io/ioutil"
//	"fmt"
//	"math"
//	"strconv"
//	"strings"
)

// 公用Controller, 其它Controller继承它
type ApiBaseContrller struct {
	controllers.BaseController // 不能用*BaseController
}


// 得到token, 这个token是在AuthInterceptor设到Session中的
func (c ApiBaseContrller) getToken() string {
	return c.Session["_token"]
}

// userId
// _userId是在AuthInterceptor设置的
func (c ApiBaseContrller) getUserId() string {
	return c.Session["_userId"]
}

// 得到用户信息
func (c ApiBaseContrller) getUserInfo() info.User {
	userId := c.Session["_userId"]
	if userId == "" {
		return info.User{}
	}
	return userService.GetUserInfo(userId);
}