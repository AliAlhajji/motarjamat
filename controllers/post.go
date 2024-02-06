package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AliAlhajji/Motarjamat/models"
	"github.com/gin-gonic/gin"
)

type postServer interface {
	CreatePost(post *models.Post) (int64, error)
	GetPost(postID int64) (*models.Post, error)
	GetAllPostsPaged(page int) ([]*models.Post, error)
	EditPost(post *models.Post) error
	DeletePost(postID int64) error
}

type postsController struct {
	postServer postServer
}

func NewPostsController(postServer postServer) *postsController {
	return &postsController{
		postServer: postServer,
	}
}

func (c *postsController) NewPost(ctx *gin.Context) {
	if ctx.Request.Method == http.MethodGet {
		ctx.HTML(http.StatusOK, "new_post.html", gin.H{"title": "New Post", "user": getContextUser(ctx)})
		return
	}

	user := getContextUser(ctx)

	var post models.Post

	err := ctx.Bind(&post)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "new_post.html", gin.H{"err": err, "user": getContextUser(ctx)})
		return
	}

	if post.Title == "" || post.Body == "" || post.Link == "" {
		ctx.HTML(http.StatusInternalServerError, "new_post.html", gin.H{"err": "All fields are required", "user": getContextUser(ctx)})
		return
	}

	post.UserID = user.UUID

	postID, err := c.postServer.CreatePost(&post)

	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "new_post.html", gin.H{"err": err, "title": "New Post", "user": getContextUser(ctx)})
		return
	}

	postLink := fmt.Sprintf("/post/%d", postID)
	ctx.Redirect(http.StatusTemporaryRedirect, postLink)
}

func (c *postsController) GetPost(ctx *gin.Context) {
	postID := ctx.Param("id")

	id, err := strconv.Atoi(postID)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
	}

	post, err := c.postServer.GetPost(int64(id))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"err": err})
	}

	ctx.HTML(http.StatusOK, "post.html", gin.H{"post": post, "title": post.Title, "user": getContextUser(ctx)})
}

func (c *postsController) AllPostsByPage(ctx *gin.Context) {
	postID := ctx.Param("page")

	page, err := strconv.Atoi(postID)
	if err != nil {
		page = 1
	}

	posts, err := c.postServer.GetAllPostsPaged(page)

	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"err": err})
	}

	user := getContextUser(ctx)

	ctx.HTML(http.StatusOK, "home.html", gin.H{"title": "مترجمات", "posts": posts, "user": user})

}

func (c *postsController) EditPost(ctx *gin.Context) {
	paramPostID := ctx.Param("id")

	postID, err := strconv.Atoi(paramPostID)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
	}

	if ctx.Request.Method == http.MethodGet {
		post, err := c.postServer.GetPost(int64(postID))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
			ctx.Abort()
			return
		}

		ctx.HTML(http.StatusOK, "edit_post.html", gin.H{"title": "Edit Post", "post": post, "user": getContextUser(ctx)})
		return
	}

	user := getContextUser(ctx)

	var post models.Post

	err = ctx.Bind(&post)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "edit_post.html", gin.H{"err": err, "user": getContextUser(ctx)})
		return
	}

	post.ID = int64(postID)

	if post.Title == "" || post.Body == "" || post.Link == "" {
		ctx.HTML(http.StatusInternalServerError, "edit_post.html", gin.H{"err": "All fields are required", "user": getContextUser(ctx)})
		return
	}

	post.UserID = user.UUID

	err = c.postServer.EditPost(&post)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "edit_post.html", gin.H{"err": err, "title": "Edit Post", "user": getContextUser(ctx)})
		return
	}

	postLink := fmt.Sprintf("/post/%d", post.ID)
	ctx.Redirect(http.StatusTemporaryRedirect, postLink)
}

func (c *postsController) Delete(ctx *gin.Context) {
	paramPostID := ctx.Param("id")

	postID, err := strconv.Atoi(paramPostID)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
	}

	post, err := c.postServer.GetPost(int64(postID))
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
		ctx.Abort()
		return
	}

	user := getContextUser(ctx)

	if post.UserID != user.UUID {
		ctx.HTML(http.StatusUnauthorized, "error.html", gin.H{"err": "You cannot delete this post"})

		return
	}

	err = c.postServer.DeletePost(int64(postID))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "edit_post.html", gin.H{"err": err, "title": "Edit Post", "user": getContextUser(ctx)})
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/home")
}
