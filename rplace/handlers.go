package rplace

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"liewell.fun/alioth/core"
)

var (
	colorWhite = "#FFFFFF"                // 初始化颜色
	width      = 32                       // 画布宽度(数量)
	height     = 16                       // 画布高度(数量)
	pixels     = make([][]string, height) // 存储每个像素的颜色
	upgrader   = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = new(sync.Map) // 当前所有的 WebSocket 客户端
	mutex   sync.Mutex
)

// 初始化 pixels 数组
func init() {
	for i := range pixels {
		pixels[i] = make([]string, width)
		for j := range pixels[i] {
			pixels[i][j] = colorWhite
		}
	}
}

// 每个像素的位置信息和颜色
type PositionColor struct {
	X int    `json:"x"`
	Y int    `json:"y"`
	C string `json:"color"`
}

// 处理 WebSocket 连接
func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		core.Logger.Errorf("[HandleWebSocket] failed upgrade: %v", err)
		return
	}
	defer conn.Close()

	// 记录当前客户端链接
	clients.Store(conn, true)
	core.Logger.Info("[HandleWebSocket] client: {} established!", conn.RemoteAddr().String())

	// 建立链接后立刻发送当前所有像素数据
	if err := conn.WriteJSON(pixels); err != nil {
		core.Logger.Errorf("[HandleWebSocket] failed send initial data: ", err)
		return
	}

	// 监听前端发送的像素更新请求
	for {
		var pc PositionColor
		if err := conn.ReadJSON(&pc); err != nil {
			core.Logger.Errorf("[HandleWebSocket] failed read message: ", err)
			clients.Delete(conn)
			return
		}

		// 更新 pixels 数据
		mutex.Lock()
		pixels[pc.Y][pc.X] = pc.C
		mutex.Unlock()

		// 广播
		broadcastUpdate(pc.X, pc.Y, pc.C)
	}
}

func broadcastUpdate(x, y int, color string) {
	clients.Range(func(key, value any) bool {
		conn := key.(*websocket.Conn)
		err := conn.WriteJSON(PositionColor{X: x, Y: y, C: color})
		if err != nil {
			core.Logger.Errorf("[HandleWebSocket] broadcasting error: ", err)
			conn.Close()
			clients.Delete(conn)
		}
		return true
	})
}
