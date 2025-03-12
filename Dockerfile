FROM alpine:3.21
RUN mkdir -p /opt/aiload/data
WORKDIR /opt/aiload
COPY aiload config.json /opt/aiload/
# 暴露文件夹和端口
EXPOSE 2081
# 启动程序
CMD ["./aiload", "start"]