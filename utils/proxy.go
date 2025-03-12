package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
)

var onceClient sync.Once
var client *http.Client

// ProxyConfig 包含代理配置信息
type ProxyConfig struct {
	TargetURL    string            // 目标URL
	Timeout      time.Duration     // 超时时间
	CustomHeader map[string]string // 自定义Header
	RemoveHeader []string          // 需要移除的响应Header
	AddHeader    map[string]string // 需要添加的响应Header
	ModelName    string            // 模型名称
}

// ReverseProxy 实现反向代理功能
func ReverseProxy(c *gin.Context, config ProxyConfig) {
	// 创建HTTP客户端
	// client := &http.Client{
	// 	Timeout: config.Timeout,
	// }
	onceClient.Do(func() {
		client = &http.Client{
			Timeout: config.Timeout,
		}
	})
	// 读取原始请求的Body
	body, err := ioutil.ReadAll(c.Request.Body)
	// 将 Body 转换为字符串
	jsonData := string(body)

	// 使用 gjson 获取 model 字段的值（可选，用于调试）
	// model := gjson.Get(jsonData, "model").String()
	// fmt.Println("Original model:", model)

	// 使用 sjson 修改 model 字段的值
	modifiedJSON, err := sjson.Set(jsonData, "model", config.ModelName)
	// 转为[]byte
	body = []byte(modifiedJSON)
	// fmt.Println(body)

	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to modify model field",
			"data": "",
		})
		return
	}

	// 创建新的请求
	req, err := http.NewRequest(c.Request.Method, config.TargetURL, bytes.NewBuffer(body))
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to create request",
			"data": "",
		})
		return
	}

	// 复制原始请求的Header到新请求
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 添加自定义Header
	for key, value := range config.CustomHeader {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to send request",
			"data": "",
		})
		return
	}
	defer resp.Body.Close()

	// 读取响应Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to read response body",
			"data": "",
		})
		return
	}

	// 移除不需要的响应Header
	for _, header := range config.RemoveHeader {
		resp.Header.Del(header)
	}

	// 添加需要的响应Header
	for key, value := range config.AddHeader {
		resp.Header.Set(key, value)
	}

	// 复制响应Header到Gin的响应
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// 设置响应状态码
	c.Status(resp.StatusCode)

	// 打印header
	// fmt.Println(resp.Header)
	// 打印目标地址
	// fmt.Println(c.Request.Method)
	c.Writer.Write(respBody)
}

// ReverseProxy 实现反向代理功能，支持流式响应
func ReverseStreamProxy(c *gin.Context, config ProxyConfig) {
	// 创建HTTP客户端
	onceClient.Do(func() {
		client = &http.Client{
			Timeout: config.Timeout,
		}
	})

	// 读取原始请求的Body
	body, err := ioutil.ReadAll(c.Request.Body)
	// fmt.Println(body)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to read request body",
			"data": "",
		})
		return
	}

	// 修改请求体中的 model 字段
	jsonData := string(body)
	modifiedJSON, err := sjson.Set(jsonData, "model", config.ModelName)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to modify model field",
			"data": "",
		})
		return
	}
	body = []byte(modifiedJSON)

	// 创建新的请求
	req, err := http.NewRequest(c.Request.Method, config.TargetURL, bytes.NewBuffer(body))
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to create request",
			"data": "",
		})
		return
	}

	// 复制原始请求的Header到新请求
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 添加自定义Header
	for key, value := range config.CustomHeader {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to send request",
			"data": "",
		})
		return
	}
	defer resp.Body.Close()

	// 移除不需要的响应Header
	for _, header := range config.RemoveHeader {
		resp.Header.Del(header)
	}

	// 添加需要的响应Header
	for key, value := range config.AddHeader {
		resp.Header.Set(key, value)
	}

	// 复制响应Header到Gin的响应
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// 设置响应状态码
	c.Status(resp.StatusCode)

	// 流式转发响应体
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		// 如果流式传输中断，记录错误日志
		c.Error(err)
		return
	}
}

// CreateReadCloser 创建一个新的基于字节的 io.ReadCloser
func CreateReadCloser(data []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewBuffer(data))
}
