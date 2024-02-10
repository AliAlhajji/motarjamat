package controllers

import (
	"net/http"
	"strconv"

	"github.com/AliAlhajji/Motarjamat/middleware"
	"github.com/AliAlhajji/Motarjamat/models"
	"github.com/gin-gonic/gin"
)

const cookieToken string = "token"

type UserServer interface {
	CreateUser(uuid string, username string, email string, name string) (id int64, err error)
	VerifyCredentials(usernameOrEmail string, password string) (uuid string, err error)
	GetUserByUsername(username string) (user *models.User, err error)
	GetUserByID(id int64) (user *models.User, err error)
	GetUserByUUID(uuid string) (user *models.User, err error)
	GetAllByPage(page int) ([]*models.User, error)
	Delete(id int64) error
}

type userController struct {
	repo UserServer
}

func NewUserController(repo UserServer) *userController {
	return &userController{
		repo: repo,
	}
}

func (s *userController) GetAll(c *gin.Context) {
	data := gin.H{}
	data["title"] = "All Users"
	data["categories"] = c.MustGet("cats")

	postID := c.Param("page")

	page, err := strconv.Atoi(postID)
	if err != nil {
		page = 1
	}

	users, err := s.repo.GetAllByPage(page)
	if err != nil {
		data["err"] = err
		c.HTML(http.StatusInternalServerError, "admin_all_users.html", data)
	}

	data["users"] = users

	c.HTML(http.StatusOK, "admin_all_users.html", data)
}

func (s *userController) CreateUser(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "register.html", nil)
		return
	}

	var userForm models.User

	err := c.Bind(&userForm)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
		return
	}

	uuid := c.GetString("uuid")
	if uuid == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"err": "Firebase uuid missing"})
		return
	}

	_, err = s.repo.CreateUser(uuid, userForm.Username, userForm.Email, userForm.Name)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"err": "user information missing"})
		return
	}

	c.HTML(http.StatusCreated, "register.html", gin.H{"user": userForm})

}

func (s *userController) Login(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}

	_, ok := c.Get(middleware.ContextUser)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"err": "could not find user"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/home")
}

func (s *userController) Logout(c *gin.Context) {
	c.SetCookie(cookieToken, "", -1, "", "", false, true)

	c.Redirect(http.StatusFound, "/home")
}

func (s *userController) GetUserByUsername(c *gin.Context) {
	var user models.User

	err := c.BindUri(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"errorMessage": err,
		})
	}
	u, err := s.repo.GetUserByUsername(user.Username)
	if err != nil {
		c.JSON(400, gin.H{
			"errorMessage": err,
		})
	}

	c.HTML(200, "user.html", u)
}

func (c *userController) Delete(ctx *gin.Context) {
	paramsUserID := ctx.Param("userID")

	userID, err := strconv.Atoi(paramsUserID)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
	}

	err = c.repo.Delete(int64(userID))
	if err != nil {
		users, _ := c.repo.GetAllByPage(1)

		ctx.HTML(http.StatusInternalServerError, "admin_all_users.html", gin.H{"err": err, "title": "All Users", "users": users})
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/admin/users")
}
