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
	session := sessionData.(*sessions.SessionCell)
	session.Set("username", "张根")
	session.Set("gender", "男")
	session.Set("age", "18")
	c.JSON(200, gin.H{"data": "登录成功"})
}

func index(c *gin.Context) {
	sessionData, ok := c.Get(sessions.ContextID)
	if !ok {
		fmt.Println("获取session data 失败！")

	}
	session := sessionData.(*sessions.SessionCell)
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
