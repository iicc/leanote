# API 列表

除了/login其它的都需要另外带参数token=xxxx

## Auth

* /api/login 登录
```
参数 email, pwd, captcha(验证码)
返回 成功返回Ok = true, Item是token, 否则Ok = false, 若Code == S_NEED_CAPTCHA 表示需要验证码
```
* /api/logout 注销
```
无参数
返回	成功返回Ok = true, 否则为false
```

* /api/register 注册
```
参数 email, pwd
返回 成功返回 Ok = true, 否则返回 false,  Msg有相应的提示
```

## User