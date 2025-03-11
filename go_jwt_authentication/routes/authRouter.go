package routes

import(
	"github.com/gin-gonic/gin"
	"go_jwt/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("/users/signup", controllers.Signup())
	incomingRoutes.POST("users/login", controllers.Login())
}
