package middleware

import (
	"fmt"
	"net/http"

	"go_jwt/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc{
	return func(c *gin.Context){
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error":fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
			}
			
			class, err := helpers.ValidateToken(clientToken)
			if err != "" {
				c.JSON(http.StatusInternalServerError, gin.H{"error":err})
				c.Abort()
				return
			}
			c.Set("email", class.Email)
			c.Set("first_name", class.First_name)
			c.Set("last_name", class.Last_name)
			c.Set("uid", class.Uid)
			c.Set("user_type", class.User_type)
			c.Next()
		

	}
}