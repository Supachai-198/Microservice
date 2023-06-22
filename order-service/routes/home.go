package routes

import "github.com/gin-gonic/gin"

func InitHomeRoute(rg *gin.RouterGroup) {
	routerGroup := rg.Group("/")

	//localhost:6000/api/v1/
	routerGroup.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Order Service v1.0.0",
		})
	})

}
