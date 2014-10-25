# API设计

By life (life@leanote.com)

## api url

所有api的url前面带/api/, 如:

`/api/user/info?userId=xxxx&token=xxxx`

## 文件目录结构
* 所有API的Controller都在app/api文件夹下
* 文件命名: Api功能Controller.go, 如ApiUserController.go
* 结构体命名为 Api功能, 如ApiUser
* API公用Controller: ApiBaseController
* init.go 注入service和定义拦截器

## 流程
用户登录后返回一个token, 以后所有的请求都携带该token. 
在init.go中的拦截器会得到token并调用sessionService判断是否登录了 

## 返回值结构
全部返回Json, 统一结构如下(注意, 首字大写)
```
{
	Ok: bool, 返回是否正确
	Msg: string, 信息, 比如用户名或密码不正确
	Code: int, 代码, -1表示未登录
	Item: interface{}, 单条记录
	List: []interface{}, 多条记录
}
```

Code 值列表:

```
S_DEFAULT = iota // 0 默认, 无意义
S_NOT_LOGIN // 1 未登录 
S_WRONG_USERNAME_PASSWORD  // 2
S_WRONG_CAPTCHA // 3
S_NEED_CAPTCHA // 4
S_NOT_OPEN_REGISTER // 4

```

如果返回用户信息, 则Item是UserInfo, 如果返回用户的笔记列表, 则List是笔记列表

## 问题

### 问题1
登录后所有的请求都带上token, 而之前的笔记中图片的链接也是一个请求, 这个怎么带token? 该token是否可以放在http header中?
而且这个请求不是API的一部分, 是WEB的controller

### 问题2
怎么同步? 同步机制, 怎么判断需要同步什么? 不可能从服务器上全部获取, 主要是笔记的同步!!

机制1: 每次同步记录下时间t, 下次同步的时候只同步时间>t的笔记. 有些笔记删除了怎么办?