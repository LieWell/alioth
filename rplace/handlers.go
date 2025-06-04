package rplace

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"liewell.fun/alioth/core"
	"liewell.fun/alioth/models"
	"liewell.fun/alioth/utils"
)

var (
	colorWhite = 1  // 初始化颜色
	colorBlack = 0  // 初始化颜色
	width      = 32 // 画布宽度(数量)
	height     = 16 // 画布高度(数量)
	upgrader   = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = new(sync.Map) // 当前所有的 WebSocket 客户端
	mutex   sync.Mutex
)

// 每个像素的位置信息和颜色
type PositionColor struct {
	X int `json:"x"`
	Y int `json:"y"`
	C int `json:"color"`
}

// 处理 WebSocket 连接
func HandleWebSocket(c *gin.Context) {
	// 建立连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		core.Logger.Error("[HandleWebSocket] failed upgrade: ", err)
		return
	}
	defer conn.Close()

	// 查询当前日期数据库中存储的数据
	// 如果已存在数据, 直接返回所有数据
	// 否则初始化一条数据入库, 然后返回
	pixels := initOrReturn()

	// 记录当前客户端链接
	clients.Store(conn, true)
	core.Logger.Infof("[HandleWebSocket] client: [%v] established!", conn.RemoteAddr().String())

	// 建立链接后立刻发送当前所有像素数据
	if err := conn.WriteJSON(pixels); err != nil {
		core.Logger.Error("[HandleWebSocket] failed send initial data: ", err)
		return
	}

	// 监听前端发送的像素更新请求
	for {
		var pc PositionColor
		if err := conn.ReadJSON(&pc); err != nil {
			core.Logger.Error("[HandleWebSocket] failed read message: ", err)
			clients.Delete(conn)
			return
		}

		// 更新 pixels 数据
		// 其值来自前端传入
		mutex.Lock()
		pixels[pc.Y][pc.X] = pc.C
		mutex.Unlock()

		// 广播
		broadcastUpdate(pc.X, pc.Y, pc.C)

		// 持久化数据
		saveOrUpdateRplace(pixels)
	}
}

// 初始化 pixels 数组
func initPixels() [][]int {
	pixels := make([][]int, height)
	for i := range pixels {
		pixels[i] = make([]int, width)
		for j := range pixels[i] {
			pixels[i][j] = colorWhite
		}
	}
	return pixels
}

func initOrReturn() [][]int {
	var pixels [][]int

	// 检查数据是否已存在
	dateNow := utils.StringDate(time.Now())
	rplcaceOld, err := models.FindRplaceByDate(dateNow)

	// 数据已存在,直接返回已经存储的数据
	if err == nil {
		_ = json.Unmarshal(rplcaceOld.Data, &pixels)
	} else {
		// 数据不存在则初始化像素数组
		pixels = initPixels()

		// 如果是数据库记录不存在,则插入一条新纪录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			jsonData, _ := json.Marshal(pixels)
			models.SaveRplace(&models.Rplace{
				Date: dateNow,
				Data: jsonData,
			})
			core.Logger.Infof("[initOrReturn] date[%v] init to db!", dateNow)
		}
	}
	return pixels
}

func broadcastUpdate(x, y int, color int) {
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

func saveOrUpdateRplace(pixels [][]int) {

	// 数据准备
	dateNow := utils.StringDate(time.Now())
	jsonData, _ := json.Marshal(pixels)

	// 检查是否存在
	rplcace, err := models.FindRplaceByDate(dateNow)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		core.Logger.Error("[saveOrUpdateRplace] find rplace error: ", err)
		return
	}

	// 不存在则新建
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		rplcace = models.EmptyRplace
	}

	// 覆盖更新
	rplcace.Date = dateNow
	rplcace.Data = jsonData
	models.SaveRplace(rplcace)
}
