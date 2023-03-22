package api

import (
	"fmt"
	"github.com/tony/mot-server/cmd/server/model"
	"strconv"

	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	// "github.com/tony/mot-server/cmd/server/model"
	"github.com/tony/mot-server/cmd/server/service"
	"net/http"
	"strings"
)

type LoginAPI struct {
	LoginSrv service.LoginSrv
}

type LoginParam struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

const (
	userkey = "user"
)

func (a LoginAPI) Login(c *gin.Context) {

	// 获取 body 中的所有数据
	var loginParam LoginParam

	err := c.BindJSON(&loginParam)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	username := loginParam.UserName
	password := loginParam.Password

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	fmt.Println("LoginAPI api --- ", username)

	// 验证用户名和密码是否正确
	verify, user := a.LoginSrv.VerifyByUsername(username, password)

	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码错误"})
		return
	}

	fmt.Println("useruseruser Status：", user)

	if verify == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 如果状态未0，说明改用户刚提交注册申请，还在请求注册中
	if user.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该用户注册申请正在审批中"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user", "data": user})
}

func (a LoginAPI) Register(c *gin.Context) {
	var user model.User

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	username := user.Username
	password := user.Password
	role := user.Role
	user.Status = 0
	if username == "admin" && password == "admin" {
		user.Status = 1
	}

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" || strings.Trim(role, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	err = service.AddUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create user error", "data": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": user.Username})
}

func (a LoginAPI) GetUserList(c *gin.Context) {
	var users []model.User
	users, err := service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "success"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": users})
}

func (a LoginAPI) GetRoleByUsername(c *gin.Context) {
	var users []model.User
	users, err := service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "success"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": users})
}

func (a LoginAPI) PermitUser(c *gin.Context) {
	userId := c.Param("userId")
	id, err := strconv.Atoi(userId)
	fmt.Println("id：", id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad param"})
		return
	}
	err = service.PermitUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "success"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "permit success"})
}

func (a LoginAPI) DeleteUser(c *gin.Context) {
	userId := c.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad param"})
		return
	}
	err = service.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "success"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "delete success"})
}
