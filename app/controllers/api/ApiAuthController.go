package api

import (
	"github.com/revel/revel"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
//	"strconv"
)

// 用户登录后生成一个token, 将这个token保存到session中
// 以后每次的请求必须带这个token, 并从session中获取userId

// 用户登录/注销/找回密码

type ApiAuth struct {
	ApiBaseContrller
}

// 登录
// 成功返回 {Ok: true, Item: token }
// 失败返回 {Ok: false, Code: 1(表示需要验证码), Item: token}
func (c ApiAuth) Login(email, pwd string, captcha string) revel.Result {
	token := c.getToken()
	var msg = ""
	code := S_DEFAULT
	
	// > 5次需要验证码, 直到登录成功
	if sessionService.LoginTimesIsOver(token) && sessionService.GetCaptcha(token) != captcha {
		msg = "captchaError"
		code = S_WRONG_CAPTCHA
	} else {
		userInfo := authService.Login(email, pwd)
		if userInfo.Email != "" {
			sessionService.ClearLoginTimes(token)
			sessionService.SetUserId(token, userInfo.UserId.Hex())
			return c.RenderJson(info.Re{Ok: true, Item: token})
		} else {
			// 登录错误, 则错误次数++
			msg = "wrongUsernameOrPassword"
			code = S_WRONG_USERNAME_PASSWORD
			sessionService.IncrLoginTimes(token)
		}
	}
	
	if sessionService.LoginTimesIsOver(token) {
		code = S_NEED_CAPTCHA
	}
	return c.RenderJson(info.Re{Ok: false, Code: code, Item: token, Msg: c.Message(msg)})
}

// 注销
func (c ApiAuth) Logout() revel.Result {
	token := c.getToken()
	sessionService.Clear(token)
	re := info.NewRe()
	re.Ok = true
	return c.RenderJson(re)
}

// 注册
func (c ApiAuth) Register(email, pwd string) revel.Result {
	re := info.NewRe()
	if !configService.IsOpenRegister() {
		re.Code = S_NOT_OPEN_REGISTER // 未开放注册
		return c.RenderJson(re)
	}
	
	if re.Ok, re.Msg = Vd("email", email); !re.Ok {
		return c.RenderRe(re);
	}
	if re.Ok, re.Msg = Vd("password", pwd); !re.Ok {
		return c.RenderRe(re);
	}
	
	// 注册
	re.Ok, re.Msg = authService.Register(email, pwd)
	
	return c.RenderRe(re)
}

// 找回密码请使用web view