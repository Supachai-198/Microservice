package routes

import (
	"net/http"

	"example.com/order-service/db"
	"example.com/order-service/middlewares"
	"example.com/order-service/models"
	"github.com/gin-gonic/gin"
)

func InitOrderRoute(rg *gin.RouterGroup) {
	routerGroup := rg.Group("/order")

	//localhost:6000/api/v1/order/
	routerGroup.POST("/", middlewares.AuthJWT(), func(c *gin.Context) {
		//get user_id from auth_jwt middleware
		user_id := c.MustGet("user_id").(float64)

		var body models.InputOrder
		c.ShouldBindJSON(&body)

		order := models.Order{
			UserId:    uint(user_id),
			ProductId: body.ProductId,
		}

		result := db.DB.Create(&order)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}

		c.JSON(200, gin.H{
			"message": "บันทึกการสั่งซื้อสำเร็จ",
		})
	})

}
