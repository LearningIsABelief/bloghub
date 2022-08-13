package response

import "github.com/gin-gonic/gin"

func Response(c *gin.Context, code int, msg string, obj ...interface{}) {
	data := map[string]interface{}{}
	if len(obj) > 1 && len(obj)%2 == 0 {
		for i := 0; i < len(obj); i += 2 {
			data[obj[i].(string)] = obj[i+1]
		}
	}
	c.JSON(code, gin.H{
		"msg":  msg,
		"data": data,
	})
}
