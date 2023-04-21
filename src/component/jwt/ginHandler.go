package jwts

import (
	"github.com/gin-gonic/gin"
)

func JwtHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		if !Serve(context) {
			return
		}
		context.Next()
	}
}
