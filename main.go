package main

import (
	"fmt"
	"gin-seeion/sessions"
	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	//从上下文中获取到session data
	sessionData, ok := c.Get(sessions.ContextID)
	if !ok {
		fmt.Println("获取session data 失败！")
	}
	//为什么一直做接口?就是能在这里动态调用。
	session := sessionData.(sessions.SessionInterface)
	session.Set("username", "张根")
	session.Set("gender", "男")
	session.Set("age", "18")
	session.Save()
	c.JSON(200, gin.H{"data": "登录成功"})
}

func index(c *gin.Context) {
	sessionData, ok := c.Get(sessions.ContextID)
	if !ok {
		fmt.Println("获取session data 失败！")

	}
	session := sessionData.(sessions.SessionInterface)
	c.JSON(200, gin.H{"data": "首页", "姓名": session.Get("username"), "性别": session.Get("gender"), "年龄": session.Get("age")})
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Use(sessions.GinSessionMiddleWare())
	router.GET("/login/", login)
	router.GET("/index/", index)
	router.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", nil)
	})
	err := router.Run(":8002")
	if err != nil {
		fmt.Println("Gin启动失败！")
	}

}
