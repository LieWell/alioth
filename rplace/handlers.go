package rplace

import (
	"encoding/json"
	"errors"
	"fmt"
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

	// 获取当前日期（月/日）
	now := time.Now()
	month := int(now.Month())
	day := now.Day()

	// 将日期格式化为两位数（例如：6月 -> "06", 15日 -> "15"）
	dateStr := fmt.Sprintf("%02d%02d", month, day)

	// 5x7 点阵数字定义（7行5列）
	digits := map[rune][7]int{
		'0': {14, 17, 17, 17, 17, 17, 14},
		'1': {4, 12, 4, 4, 4, 4, 14},
		'2': {14, 17, 1, 2, 4, 8, 31},
		'3': {14, 17, 1, 6, 1, 17, 14},
		'4': {17, 17, 17, 31, 1, 1, 1},
		'5': {31, 16, 30, 1, 1, 17, 14},
		'6': {14, 17, 16, 30, 17, 17, 14},
		'7': {31, 1, 2, 4, 4, 4, 4},
		'8': {14, 17, 17, 14, 17, 17, 14},
		'9': {14, 17, 17, 15, 1, 17, 14},
	}

	// 计算起始位置（居中显示）
	startRow := (16 - 7) / 2         // 垂直居中
	startCol := (32 - (5*4 + 3)) / 2 // 水平居中（4个数字，每个5列，3个间距）

	// 在像素阵列中绘制日期数字
	for i, char := range dateStr {
		digitPattern, exists := digits[char]
		if !exists {
			continue
		}

		// 每个数字起始位置（考虑间距）
		numStartCol := startCol + i*(5+1)

		// 绘制数字（7行）
		for row := 0; row < 7; row++ {
			// 确保行位置有效
			pixelRow := startRow + row
			if pixelRow >= 16 {
				continue
			}

			// 绘制一行（5列）
			pattern := digitPattern[row]
			for col := 0; col < 5; col++ {
				// 检查点阵中的位（从最高位开始）
				if pattern&(1<<(4-col)) != 0 {
					pixelCol := numStartCol + col
					if pixelCol < 32 {
						pixels[pixelRow][pixelCol] = 0
					}
				}
			}
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
