package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserResponse 是用户查询接口返回的数据结构。
type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// NewRouter 创建本章示例使用的 Gin 路由引擎。
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// 基础健康检查路由，用于确认服务已经能够处理请求。
	router.GET("/ping", pingHandler)

	// 通过静态路由和动态路由演示静态路径的匹配优先级。
	router.GET("/user/create", createUserHandler)
	router.GET("/user/:name", userHandler)

	// 通过查询参数演示字符串读取、默认值和整数校验。
	router.GET("/search", searchHandler)

	// 通配路由演示如何读取多个后续路径段。
	router.GET("/files/*path", filesHandler)

	// API 路由使用分组表达版本和公共路径前缀。
	apiV1 := router.Group("/api/v1")
	apiV1.GET("/users", listUsersHandler)
	apiV1.GET("/users/:id", getUserHandler)

	return router
}

// pingHandler 返回服务健康检查结果。
func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// createUserHandler 返回静态路由命中的结果。
func createUserHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "static"})
}

// userHandler 返回动态路径参数。
func userHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"name": c.Param("name")})
}

// searchHandler 读取查询参数并校验分页参数。
func searchHandler(c *gin.Context) {
	pageText := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageText)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be a positive integer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"keyword": c.Query("q"),
		"name":    c.Query("name"),
		"page":    page,
	})
}

// filesHandler 返回通配路由捕获到的剩余路径。
func filesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"path": c.Param("path")})
}

// listUsersHandler 返回内存中的示例用户列表。
func listUsersHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"users": []UserResponse{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
		},
	})
}

// getUserHandler 解析用户 ID 并返回一个示例用户。
func getUserHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{ID: id, Name: "User"})
}

func main() {
	if err := NewRouter().Run(":8080"); err != nil {
		// 启动失败必须显式报告，不能静默退出。
		panic(err)
	}
}
