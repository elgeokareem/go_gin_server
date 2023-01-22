package middlewares

import (
	"github.com/gin-gonic/gin"
)

func GetJwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token := c.Request.Header["Authorization"]
		// tokenParsed, err := jwt.Parse(token[0], func(token *jwt.Token) (interface{}, error) {
		// 	return []byte("your-256-bit-secret"), nil
		// })

		// fmt.Println(tokenParsed)

		// if err != nil {
		// 	fmt.Println(err)
		// 	c.JSON(http.StatusInternalServerError, "Error")
		// }

		// hash := sha256.New()

		// c.JSON(http.StatusOK, "todo bien")
	}
}
