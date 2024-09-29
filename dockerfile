# 使用官方 Golang 映像作為基礎映像
FROM golang:1.22.1-alpine

# 在容器中建立應用程序的工作目錄
WORKDIR /app

# 將本地的代碼複製到容器中
COPY . .

# 使用 go mod 下載依賴
RUN go mod download

ARG ENV_MODE

ENV APP_ENV = $ENV_MODE

# 編譯應用程序
RUN go build main.go

# 暴露應用程序運行的端口號
EXPOSE 8080

# 運行應用程序

CMD ["./main"]