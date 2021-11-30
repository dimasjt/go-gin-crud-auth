package api

import (
	"crud-with-auth/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *Api) Auth() {
	noAuth := api.r.Group("/users")
	noAuth.POST("/login", api.UserLoginHandler)
	noAuth.POST("/register", api.UserRegisterHandler)
	noAuth.POST("/refresh-token", api.UserRefreshHandler)
}

func (api Api) UserLoginHandler(c *gin.Context) {
	var userDto db.User
	var user db.User
	c.BindJSON(&userDto)

	if err := api.db.Storage.Where("email = ?", userDto.Email).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Email or password invalid",
		})
		return
	}

	if err := user.ValidatePassword(userDto.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Email or password invalid",
		})
		return
	}

	if token, err := user.Token(); err != nil {
		panic(err)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successfully",
			"token":   token,
		})
	}
}

func (api Api) UserRegisterHandler(c *gin.Context) {
	var user db.User
	c.BindJSON(&user)
	user.GeneratePassword()
	if err := api.db.Storage.Create(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered",
	})
}

func (api Api) UserRefreshHandler(c *gin.Context) {

}
