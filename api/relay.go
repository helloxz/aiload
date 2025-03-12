package api

import (
	"aiload/utils"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

// 声明一个结构体
type Model struct {
	BaseURL string `json:"base_url"`
	Model   string `json:"model"`
	ApiKey  string `json:"api_key"`
}

func Relay(c *gin.Context) {
	// 读取 JSON 文件内容
	data, err := os.ReadFile("data/config/config.json")
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to read configuration file!",
			"data": "",
		})
		c.Abort()
		return
	}

	// 获取原始请求体数据
	rawData, err := c.GetRawData()
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to read request body!",
			"data": "",
		})
		c.Abort()
		return
	}

	stream := gjson.GetBytes(rawData, "stream").Bool()

	// 将文件内容转换为字符串
	jsonData := string(data)

	// 重要：将请求体放回，这样后续处理器仍然可以读取它
	c.Request.Body = utils.CreateReadCloser(rawData)

	// 使用 gjson 提取 models 列表
	models := gjson.Get(jsonData, "models")
	if !models.Exists() {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Models list not found!",
			"data": "",
		})
		c.Abort()
		return
	}

	var modelList []Model

	// 遍历 models 列表
	models.ForEach(func(_, model gjson.Result) bool {
		baseURL := model.Get("base_url").String()
		modelName := model.Get("model").String()
		api_key := model.Get("api_key").String()

		// 将每个模型的信息添加到 modelList 中
		modelList = append(modelList, Model{
			BaseURL: baseURL,
			Model:   modelName,
			ApiKey:  api_key,
		})
		return true // 继续遍历下一个元素
	})

	// 随机从 modelList 中选择一个模型
	randIndex := rand.Intn(len(modelList))
	// 检查内存中是否存在index
	index := getModelIndex(c)
	// 如果不存在
	if index == 0 {
		// 将选择的模型的 index 存入缓存
		setModelIndex(c, randIndex)
	} else {
		randIndex = index
	}

	selectedModel := modelList[randIndex]

	// fmt.Println(selectedModel)
	timeout := viper.GetInt("settins.timeout")
	// 如果为0，则设置默认值
	if timeout == 0 {
		timeout = 120
	}
	config := utils.ProxyConfig{
		TargetURL:    selectedModel.BaseURL,
		Timeout:      time.Duration(timeout) * time.Second, // 设置超时时间
		CustomHeader: map[string]string{"Authorization": "Bearer " + selectedModel.ApiKey},
		RemoveHeader: []string{},
		AddHeader:    map[string]string{"X-Model": selectedModel.Model}, // 添加自定义响应头
		ModelName:    selectedModel.Model,
	}

	if stream {
		// 请求转发服务
		utils.ReverseStreamProxy(c, config)
		return
	} else {
		// 请求转发服务
		utils.ReverseProxy(c, config)
	}

}

// 获取模型的index
func getModelIndex(c *gin.Context) int {
	// 获取用户IP
	ip := utils.GetClientIP(c)
	// 生成key
	keyStr := "model_index_" + ip
	// 转为[]byte
	key := []byte(keyStr)
	// 读取缓存
	value, err := utils.Cache.Get(key)
	if err != nil {
		return 0
	}
	if value == nil {
		return 0
	}
	// 转为int
	index, err := strconv.Atoi(string(value))
	if err != nil {
		return 0
	}
	return index
}

// 设置模型的index
func setModelIndex(c *gin.Context, index int) {
	// 获取用户IP
	ip := utils.GetClientIP(c)
	// 生成key
	keyStr := "model_index_" + ip
	// 转为[]byte
	key := []byte(keyStr)
	// 转为[]byte
	value := []byte(strconv.Itoa(index))
	// 设置缓存，有效期5分钟
	utils.Cache.Set(key, value, 60*5)
}
