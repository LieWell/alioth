package auth

import (
	"net/http"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"liewell.fun/alioth/core"
	"liewell.fun/alioth/models"
	"liewell.fun/alioth/utils"
)

func Login(c *gin.Context) {

	// 获取参数
	var req UserRegisterLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Logger.Error("[Login] should bind json error: %v", err)
		c.JSON(http.StatusBadRequest, core.BadRequestError(LackOfParamCode, LackOfParamMessage))
		return
	}

	// 查询用户
	user, err := models.FindUserByUsername(req.Username)
	if err != nil {
		core.Logger.Error("[Login] find user error: %v", err)
		c.JSON(http.StatusBadRequest, core.BadRequestError(UserNameOrPasswordIncorrectCode, UserNameOrPasswordIncorrectMessage))
		return
	}

	// 校验密码是否正确
	if checkPassword(req.Password, user.Salt, user.Password) != nil {
		core.Logger.Error("[Login] username or password error!")
		c.JSON(http.StatusBadRequest, core.BadRequestError(UserNameOrPasswordIncorrectCode, UserNameOrPasswordIncorrectMessage))
		return
	}

	// 签发 token
	expireAt := time.Now().Add(time.Second * time.Duration(core.GlobalConfig.JWT.Expire))
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Username: req.Username,
		RegisteredClaims: validator.RegisteredClaims{
			Expiry:   expireAt.Unix(),
			Issuer:   core.GlobalConfig.JWT.Issuer,
			Audience: core.GlobalConfig.JWT.Audience,
		},
	}).SignedString([]byte(core.GlobalConfig.JWT.Secret))
	if err != nil {
		core.Logger.Error("[Login] issue token error: %v", err)
		c.JSON(http.StatusInternalServerError, core.SimpleInternalServerError())
		return
	}

	c.JSON(http.StatusOK, Token{
		Token:    token,
		ExpireAt: expireAt,
	})
}

func Register(c *gin.Context) {

	// 获取参数
	var req UserRegisterLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Logger.Error("[Register] should bind json error: %v", err)
		c.JSON(http.StatusBadRequest, core.BadRequestError(LackOfParamCode, LackOfParamMessage))
		return
	}

	// 校验用户名是否已占用
	user, err := models.FindUserByUsername(req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		core.Logger.Error("[Register] find user error: %v", err)
		c.JSON(http.StatusInternalServerError, core.SimpleInternalServerError())
		return
	}
	if user != nil {
		core.Logger.Error("[Register] dumplicate user: %v", req.Username)
		c.JSON(http.StatusConflict, core.ConflictError(DumplicateUserNameCode, DumplicateUserNameMessage))
		return
	}

	// 随机生成 salt 值
	var salt = utils.RandomString(16)

	// 计算 Hash 值
	password, _ := hashPassword(req.Password, salt)

	// 保存用户
	record := &models.User{
		Username: req.Username,
		NickName: req.Username,
		Email:    req.Email,
		Password: password,
		Salt:     salt,
		Status:   models.UserStatusEnable,
	}
	if _, err := models.SaveUser(record); err != nil {
		core.Logger.Fatalf("[Register] save user error: %v", err)
		c.JSON(http.StatusInternalServerError, core.SimpleInternalServerError())
		return
	}

	c.JSON(http.StatusOK, core.SuccessResponse)
}
