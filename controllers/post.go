package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AliAlhajji/Motarjamat/models"
	"github.com/AliAlhajji/Motarjamat/utils"
	"github.com/gin-gonic/gin"
)

type postServer interface {
	CreatePost(post *models.Post) (int64, error)
	GetPost(postID int64) (*models.Post, error)
	GetAllPostsPaged(page int) ([]*models.Post, error)
	GetByCategory(categoryID int64, page int) ([]*models.Post, error)
	EditPost(post *models.Post) error
	DeletePost(postID int64) error
}

type postsController struct {
	postServer     postServer
	categoryServer categoryServer
}

func NewPostsController(postServer postServer, catecategoryServer categoryServer) *postsController {
	return &postsController{
		postServer:     postServer,
		categoryServer: catecategoryServer,
	}
}

func (c *postsController) AddCategories() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cats, _ := c.categoryServer.GetUsedCategories()
		ctx.Set("cats", cats)
		ctx.Next()
	}
}

func (c *postsController) NewPost(ctx *gin.Context) {
	if ctx.Request.Method == http.MethodGet {
		cats, err := c.categoryServer.GetAll()
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"err": err})

		}

		ctx.HTML(http.StatusOK, "new_post.html", gin.H{
			"title":      "New Post",
			"user":       utils.GetContextUser(ctx),
			"categories": cats,
		})
		return
	}

	user := utils.GetContextUser(ctx)

	var post models.Post

	err := ctx.Bind(&post)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "new_post.html", gin.H{"err": err, "user": utils.GetContextUser(ctx)})
		return
	}

	if post.Title == "" || post.Body == "" || post.Link == "" {
		ctx.HTML(http.StatusInternalServerError, "new_post.html", gin.H{"err": "All fields are required", "user": utils.GetContextUser(ctx)})
		return
	}

	post.UserID = user.UUID

	postID, err := c.postServer.CreatePost(&post)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"err": err})

	}

	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "new_post.html", gin.H{
			"err":   err,
			"title": "New Post",
			"user":  utils.GetContextUser(ctx),
		})
		return
	}

	postLink := fmt.Sprintf("/post/%d", postID)
	ctx.Redirect(http.StatusTemporaryRedirect, postLink)
}

func (c *postsController) GetPost(ctx *gin.Context) {
	postID := ctx.Param("postID")

	id, err := strconv.Atoi(postID)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
	}

	post, err := c.postServer.GetPost(int64(id))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"err": err})
	}

	cats := ctx.MustGet("cats")

	ctx.HTML(http.StatusOK, "post.html", gin.H{"post": post, "title": post.Title, "user": utils.GetContextUser(ctx), "categories": cats})
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

	user := utils.GetContextUser(ctx)

	cats := ctx.MustGet("cats")

	ctx.HTML(http.StatusOK, "home.html", gin.H{"title": "مترجمات",
		"settings":   ctx.MustGet("siteSettings"),
		"posts":      posts,
		"user":       user,
		"categories": cats,
	})

}

func (c *postsController) CategoryPosts(ctx *gin.Context) {
	postID := ctx.Param("page")

	page, err := strconv.Atoi(postID)
	if err != nil {
		page = 1
	}

	cat := ctx.Param("category")

	categoryID, err := strconv.Atoi(cat)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"err": err})
		return

	}

	category, err := c.categoryServer.GetCategory(int64(categoryID))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"err": err})
		return
	}

	posts, err := c.postServer.GetByCategory(int64(categoryID), page)

	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"err": err})
		return
	}

	user := utils.GetContextUser(ctx)

	cats := ctx.MustGet("cats")

	ctx.HTML(http.StatusOK, "category.html", gin.H{
		"title":      category.Title,
		"posts":      posts,
		"user":       user,
		"categories": cats,
		"category":   category,
	})

}

func (c *postsController) EditPost(ctx *gin.Context) {
	paramPostID := ctx.Param("postID")

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

		ctx.HTML(http.StatusOK, "edit_post.html", gin.H{"title": "Edit Post", "post": post, "user": utils.GetContextUser(ctx)})
		return
	}

	user := utils.GetContextUser(ctx)

	var post models.Post

	err = ctx.Bind(&post)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "edit_post.html", gin.H{"err": err, "user": utils.GetContextUser(ctx)})
		return
	}

	post.ID = int64(postID)

	if post.Title == "" || post.Body == "" || post.Link == "" {
		ctx.HTML(http.StatusInternalServerError, "edit_post.html", gin.H{"err": "All fields are required", "user": utils.GetContextUser(ctx)})
		return
	}

	post.UserID = user.UUID

	err = c.postServer.EditPost(&post)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "edit_post.html", gin.H{"err": err, "title": "Edit Post", "user": utils.GetContextUser(ctx)})
		return
	}

	postLink := fmt.Sprintf("/post/%d", post.ID)
	ctx.Redirect(http.StatusTemporaryRedirect, postLink)
}

func (c *postsController) Delete(ctx *gin.Context) {
	paramPostID := ctx.Param("postID")

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

	user := utils.GetContextUser(ctx)

	if post.UserID != user.UUID {
		ctx.HTML(http.StatusUnauthorized, "error.html", gin.H{"err": "You cannot delete this post"})

		return
	}

	err = c.postServer.DeletePost(int64(postID))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "edit_post.html", gin.H{"err": err, "title": "Edit Post", "user": utils.GetContextUser(ctx)})
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/home")
}
