package controllers

import (
	"net/http"
	"strconv"

	"github.com/AliAlhajji/Motarjamat/models"
	"github.com/gin-gonic/gin"
)

type categoryServer interface {
	Add(title string) error
	GetAll() ([]*models.Category, error)
	GetUsedCategories() ([]*models.Category, error)
	GetCategory(id int64) (*models.Category, error)
	DeleteCategory(id int64) error
	UpdateCategory(id int64, updatedCategory *models.Category) error
}

type cateogryController struct {
	categoryServer categoryServer
}

func NewCategoryController(categoryServer categoryServer) *cateogryController {
	return &cateogryController{
		categoryServer: categoryServer,
	}
}

func (c *cateogryController) ShowAll(ctx *gin.Context) {
	data := gin.H{}
	data["title"] = "All Categories"

	cats, err := c.categoryServer.GetAll()
	if err != nil {
		data["err"] = err
		ctx.HTML(http.StatusInternalServerError, "error", data)
		return
	}

	data["categories"] = cats
	data["settings"] = ctx.MustGet("siteSettings")

	ctx.HTML(http.StatusOK, "all_categories.html", data)
}

func (c *cateogryController) AddNew(ctx *gin.Context) {
	data := gin.H{}

	data["title"] = "Add Category"

	if ctx.Request.Method == http.MethodGet {
		ctx.HTML(http.StatusOK, "add_category.html", data)
		return
	}

	var category models.Category

	err := ctx.Bind(&category)
	if err != nil {
		data["err"] = err
		ctx.HTML(http.StatusBadRequest, "error.html", data)
		return
	}

	if category.Title == "" {
		data["err"] = "Title field is required"
		ctx.HTML(http.StatusBadRequest, "add_category.html", data)
		return
	}
	err = c.categoryServer.Add(category.Title)
	if err != nil {
		data["err"] = err
		ctx.HTML(http.StatusBadRequest, "add_category.html", data)
		return
	}

	ctx.HTML(http.StatusTemporaryRedirect, "all_categories.html", gin.H{"title": "All Categories", "msg": "Added", "settings": ctx.MustGet("siteSettings")})
}

func (c *cateogryController) DeleteCategory(ctx *gin.Context) {
	paramCategoryID := ctx.Param("category")

	categoryID, err := strconv.Atoi(paramCategoryID)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
		return
	}

	err = c.categoryServer.DeleteCategory(int64(categoryID))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "all_categories.html", gin.H{"err": err, "title": "All Categories"})
		return
	}

	ctx.HTML(http.StatusOK, "all_categories.html", gin.H{"title": "All Categories", "msg": "تم الحذف بنجاح"})
}
func (c *cateogryController) UpdateCategory(ctx *gin.Context) {
	data := gin.H{}

	paramCategoryID := ctx.Param("category")

	var category models.Category

	categoryID, err := strconv.Atoi(paramCategoryID)
	if err != nil {
		data["err"] = err
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
		return
	}

	err = ctx.Bind(&category)
	if err != nil {
		data["err"] = err
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
		return
	}

	category.ID = int64(categoryID)

	err = c.categoryServer.UpdateCategory(category.ID, &category)
	if err != nil {
		data["err"] = err
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
		return
	}

	data["title"] = "All Categories"
	data["msg"] = "Updated"
	ctx.HTML(http.StatusOK, "all_categories.html", data)
}
