FROM golang:1.18.2
RUN go env -w GOPROXY="https://proxy.golang.com.cn,direct"
