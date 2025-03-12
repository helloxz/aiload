English | [中文](./README.zh.md)

___

# AILoad

AILoad is a lightweight and efficient tool that distributes requests to multiple AI model interfaces through load balancing, effectively solving the issue of frequency limits for a single AI model. Currently, it supports round-robin and IP-based load balancing strategies, ensuring that requests from the same IP are allocated to the same AI model within a short time frame, thereby enhancing consistency and user experience.

## Main Features

* [x] Load Balancing: Automatically distributes requests randomly to multiple AI model interfaces, avoiding single model frequency limits and improving overall throughput.
* [x] Unified API Entry: Provides a single entry point for invocation.
* [x] Compatible with OpenAI: Compatible with OpenAI API interface format.
* [x] Stream Return: Supports stream transmission.
* [x] IP Affinity: Requests from the same IP are fixedly assigned to one AI model within a certain period, enhancing the consistency of continuous dialogues.
* [ ] Disable IP affinity via configuration
* [ ] Specify particular model invocation
* [ ] Model grouping function


## Use Cases

Some vendors currently provide free AI model interfaces, but they often come with frequency restrictions. By using AILoad, multiple AI interfaces can be integrated together and the issue of frequency limits for a single AI model can be alleviated through load balancing.

## Deployment

Currently, only Docker deployment is supported. Please ensure you have installed the Docker environment.

**Method One: Docker Compose Installation**

The content of `docker-compose.yaml` is as follows:

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

Run the command `docker-compose up -d`.

**Method Two: Docker Command Line Installation**

```bash
docker run -d \
  --name aiload \
  -v /opt/aiload/data:/opt/aiload/data \
  -p 2081:2081 \
  --restart always \
  helloz/aiload
```

### Configuration File

The configuration file is located at `config/config.json` under the mounted directory, for example: `/opt/aiload/data/config.json`, with default content as follows:

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

* `auth_token`: Authorization key (required for calling AILoad), it will be automatically set upon first launch, but you can also modify it manually.
* `models`: Model list, you can add multiple model interfaces here (Note: The model interface must be compatible with the OpenAI API format).
* `settings.timeout`: Backend timeout time, default `120s`, which you can also modify yourself.

> Note: After modifying the configuration file, please verify the correctness of the JSON format first to avoid program runtime anomalies. Also, changes to the configuration require a container restart to take effect!

### Usage

AILoad fully supports the OpenAI API format, and you can call it like this:

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

* `sk-xxx`: Corresponds to the `auth_token` in the configuration file.
* `model`: Fixed as `auto`.

If you need stream transmission, just add the `stream:true` parameter:

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

> Other parameters can also be passed, such as `top_p/temperature` or other OpenAI parameter formats.

## Contact Us

For any questions or feedback, please submit them through GitHub Issues.