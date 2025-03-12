package api

import (
	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	// 准备 HTML 内容
	htmlContent := `<p>Please visit <a href="https://github.com/helloxz/aiload" target="_blank">helloxz/aiload</a> for usage instructions.</p>`

	// 设置 Content-Type 为 text/html
	c.Header("Content-Type", "text/html; charset=utf-8")

	// 返回 HTML 内容
	c.String(200, htmlContent)
}
