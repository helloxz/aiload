# AILoad

AILoad 是一个轻量级且高效的工具，通过负载均衡的方式将请求分发到多个 AI 模型接口，从而有效解决单一 AI 模型频率上限问题。目前支持轮询和基于 IP 的负载均衡策略，能够在短时间内确保同一 IP 的请求被分配到同一个 AI 模型，提升一致性和用户体验。

## 主要功能

* 负载均衡 ：自动将请求随机分发到多个 AI 模型接口，避免单个模型频率上限，提升整体吞吐量。
* 统一API入口 ：提供单一入口调用
* 兼容OpenAI：兼容OpenAI API接口格式
* 流式返回：支持流式传输
* IP 亲和性：同一个IP在一定时间内固定分配到一个AI模型，提升连续对话一致性


## 使用场景

目前市面部分厂商提供了免费的AI模型接口，但往往存在频率限制，通过AILoad可以将多个AI接口整合到一起，并通过负载均衡的方式来缓解单一AI模型频率上限的问题。

## 部署

目前仅支持Docker部署，请确保您已经安装Docker环境。

**方式一：Docker Compose安装**

`docker-compose.yaml`内容如下：

```yaml
version: '3'
services:
    aiload:
        container_name: aiload
        volumes:
            - '/opt/aiload/data:/opt/aiload/data'
        restart: always
        ports:
          - '2081:2081'
        image: 'helloz/aiload'
```

输入命令`docker-compose up -d`运行。

**方式二：Docker命令行安装**

```bash
docker run -d \
  --name aiload \
  -v /opt/aiload/data:/opt/aiload/data \
  -p 2081:2081 \
  --restart always \
  helloz/aiload
```

### 配置文件

配置文件位于挂载目录下的`config/config.json`，比如：`/opt/aiload/data/config.json`，默认内容如下：

```json
{
    "auth_token": "",
    "models": [
      {
        "api_key": "sk-xxx",
        "base_url": "https://api.openai.com/v1/chat/completions",
        "model": "gpt-4o"
      }
    ],
    "server": {
      "model": "release",
      "port": 2081
    },
    "settings":{
      "timeout":120
    }
}
```

* `auth_token`：授权密钥（调用AILoad需要），首次启动会自动设置，您也可以手动修改
* `models`：模型列表，您可以在这里填写多个模型接口（注意：模型接口需要兼容OpenAI API格式）
* `settings.timeout`：后端超时时间，默认`120s`，您也可以自行修改

> 注意：修改配置文件后请先校验json格式是否正确，以免导致程序运行异常。另外每次修改配置后需要重启容器才会生效！

### 使用

AILoad完全兼容OpenAI API格式，您可以像下面这样调用：

```bash
curl http://IP:2081/v1/chat/completions \
    -H "Authorization: Bearer sk-xxx" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "auto",
    "messages": [
      {
        "role": "assistant",
        "content": "You are a helpful assistant."
      },
      {
        "role": "user",
        "content": "who are you?"
      }
    ]
  }'
```

* `sk-xxx`：对应配置文件中的`auth_token`
* `model`：固定为`auto`

如果需要流式传输，只需要添加`stream:true`参数即可：

```json
{
    "model": "auto",
    "stream": true,
    "messages": [
      {
        "role": "assistant",
        "content": "You are a helpful assistant."
      },
      {
        "role": "user",
        "content": "who are you?"
      }
    ]
}
```

> 也支持传递其他参数，比如`top_p/temperature`或其他OpenAI参数格式。

## 联系我们

如有任何问题或反馈，请通过 GitHub Issues提交。