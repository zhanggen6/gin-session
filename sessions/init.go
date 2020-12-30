package sessions

import (
	"github.com/gin-gonic/gin"
)

const (
	sessionID = "sessionID" //在cookie中key名称
	ContextID = "session"   //在gin contentext中key的名称
	domain    = "/"         //设置cookie的作用域名
	maxAge    = 3600        //设置cookie的超时时间


)

var SessionStore StoreInterface

func init() {
	SessionStore = NewMemStore()

}

//实现1个gin框架的中间件
func GinSessionMiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {
		//1.根据约定好的sessionID从请求的cookie中获取唯一标识
		uniqueString, err := c.Cookie(sessionID)
		//2.cookie中获取不到sessionID--->给用户创建一个存储session的地方
		if err != nil {
			uniqueString = SessionStore.SetSession()
		}
		sessionData := SessionStore.GetSession(uniqueString)
		//2.session仓库里（服务端重启了！）获取不到sessionID--->给用户创建一个存储session的地方
		if sessionData == nil {
			//更新服务端session信息
			uniqueString = SessionStore.SetSession()
			sessionData = SessionStore.GetSession(uniqueString)
		}
		//3.初始化记录session的session id
		sessionData.Init(uniqueString)
		c.Set(ContextID, sessionData)
		//4.回写cookie一定要再next（）也就是视图函数返回之前
		c.SetCookie(sessionID, uniqueString, maxAge, "/", domain, false, false)
		c.Next()

	}

}
