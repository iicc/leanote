package api

import (
	"github.com/revel/revel"
//	"encoding/json"
//	"gopkg.in/mgo.v2/bson"
//	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/lea/captcha"
//	"github.com/leanote/leanote/app/types"
//	"io/ioutil"
//	"fmt"
//	"math"
//	"os"
//	"path"
//	"strconv"
//	"net/http"
)

// 验证码服务
type ApiCaptcha struct {
	ApiBaseContrller
}

func (c ApiCaptcha) Get() revel.Result {
	c.Response.ContentType = "image/png"
	image, str := captcha.Fetch()
	image.WriteTo(c.Response.Out)

	sessionService.SetCaptcha(c.getToken(), str)
	
	return c.Render()
}