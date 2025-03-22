package main

import (
	"github.com/gin-gonic/gin"
	"webook/internal/web"
)

func main() {
	server := gin.Default()

	u := &web.UserHandler{}
	u.RegisterRoutesv1(server.Group("/users"))
	u.RegisterRoutes(server)
	
	server.Run(":8080")
}
