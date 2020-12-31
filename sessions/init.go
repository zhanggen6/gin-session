package sessions

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	sessionID     = "sessionID" //在cookie中key名称
	ContextID     = "session"   //在gin contentext中key的名称
	domain        = "127.0.0.1" //设置cookie的域
	path          = "/"         //设置cookie的域名
	MaxAge        = 3600        //设置cookie的超时时间
	storageEngine = "redis"     //local|redis
)

//存储配置
type Option struct {
	ip  string
	pwd string
}

var sessionStorage StoreInterface

//初始化数据引擎：这里可以增加配置文件和反射
func init() {
	switch storageEngine {
	case "redis":
		sessionStorage = NewRedisStore()
	case "local":
		sessionStorage = NewMemStore()

	}
	sessionStorage.InitStore(Option{ip: "YourRedisHost:Port", pwd: "YourPassword"})

}

//实现1个gin框架的中间件
func GinSessionMiddleWare() gin.HandlerFunc {
	var sessionData SessionInterface
	return func(c *gin.Context) {
		//1.根据约定好的sessionID从请求的cookie中获取唯一标识
		uniqueString, err := c.Cookie(sessionID)
		//2.cookie中获取不到sessionID--->给用户创建一个存储session的地方
		if err != nil || len(uniqueString) < 1 {
			fmt.Println("用户有可能删了sessionID，或者首次登录！")
			sessionData = sessionStorage.CreateSession()
			uniqueString = sessionData.GetKey()
		}
		//每次都从redis获取数据：可
		sessionData = sessionStorage.GetSession(uniqueString)
		//2.session仓库里（服务端重启了！）获取不到sessionID--->给用户创建一个存储session的地方
		if sessionData == nil {
			//更新服务端session信息
			sessionData = sessionStorage.CreateSession()
			uniqueString = sessionData.GetKey()
		}
		//3.初始化记录session的session id
		c.Set(ContextID, sessionData)
		//4.回写cookie一定要再next（）也就是视图函数返回之前
		c.SetCookie(sessionID, uniqueString, MaxAge, path, domain, false, false)
		c.Next()

	}

}
