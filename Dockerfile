# docker build --network=host -t alioth:nightly .
FROM golang:1.22 as builder
WORKDIR /app
ADD . ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o alioth

FROM alpine:3.20.1
WORKDIR /app
COPY --from=builder /app/alioth ./
EXPOSE 80
EXPOSE 443
CMD ["/app/alioth","-c","/app/config/config.yaml"]