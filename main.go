package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/AliAlhajji/Motarjamat/controllers"
	firebaseservice "github.com/AliAlhajji/Motarjamat/firebase_service"
	"github.com/AliAlhajji/Motarjamat/middleware"
	sqlitedb "github.com/AliAlhajji/Motarjamat/sqlite-db"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Static("/assets/js", "./assets/js")
	r.Static("/assets/css", "./assets/css")

	files, err := loadTemplates("templates")
	if err != nil {
		panic(err)
	}
	// r.LoadHTMLGlob("./templates/**/*")
	r.LoadHTMLFiles(files...)

	err = sqlitedb.Connect()
	if err != nil {
		panic(err)
	}

	userRepo, err := sqlitedb.InitUserRepo()
	if err != nil {
		panic(err)
	}

	postRepo, err := sqlitedb.InitPostRepo()
	if err != nil {
		panic(err)
	}

	categoryRepo, err := sqlitedb.InitCategoryRepo()
	if err != nil {
		panic(err)
	}

	settingsRepo, err := sqlitedb.InitSettingsRepo()

	firebaseApp, err := firebaseservice.NewFirebaseApp("./firebase-sdk.json")
	if err != nil {
		panic(err)
	}

	firebaseAuthClient, err := firebaseservice.NewAuthService(context.Background(), firebaseApp)
	if err != nil {
		panic(err)
	}

	userService := controllers.NewUserController(userRepo)
	postControllrt := controllers.NewPostsController(postRepo, categoryRepo)
	categoryController := controllers.NewCategoryController(categoryRepo)
	settingsController := controllers.NewSettingsController(settingsRepo)

	authMiddleware := middleware.NewAuthMiddleware(firebaseAuthClient, userRepo)
	settingsMiddleware := middleware.NewSettingsMiddleware(settingsRepo)

	r.Use(authMiddleware.ExtractUserFromToken(), settingsMiddleware.GetSettings(), postControllrt.AddCategories())

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(authMiddleware.EnsureAdmin())
	adminRoutes.GET("/", settingsController.AdminHome)
	adminRoutes.POST("/update_settings", settingsController.UpdateSettings)
	adminRoutes.GET("/categories", categoryController.ShowAll)
	adminRoutes.GET("/categories/add", categoryController.AddNew)
	adminRoutes.POST("/categories/add", categoryController.AddNew)
	adminRoutes.GET("/category/delete/:category", categoryController.DeleteCategory)
	adminRoutes.POST("/category/update/:category", categoryController.UpdateCategory)
	adminRoutes.GET("/users/:page", userService.GetAll)
	adminRoutes.GET("/users", userService.GetAll)
	adminRoutes.GET("/user/delete/:userID", userService.Delete)

	r.GET("/register", userService.CreateUser)
	r.POST("/register", authMiddleware.RegisterNewUser(), userService.CreateUser)
	r.GET("/login", userService.Login)
	r.POST("/login", userService.Login)
	r.GET("/", postControllrt.AllPostsByPage)
	r.GET("/:page", postControllrt.AllPostsByPage)
	r.GET("/new_post", authMiddleware.Authenticate(), postControllrt.NewPost)
	r.POST("/new_post", authMiddleware.Authenticate(), postControllrt.NewPost)
	r.Any("/post/:postID", postControllrt.GetPost)
	r.GET("/edit_post/:postID", authMiddleware.Authenticate(), postControllrt.EditPost)
	r.POST("/edit_post/:postID", authMiddleware.Authenticate(), postControllrt.EditPost)
	r.GET("/delete_post/:postID", authMiddleware.Authenticate(), postControllrt.Delete)
	r.GET("/category/:category", postControllrt.CategoryPosts)

	err = r.Run()
	if err != nil {
		log.Panic(err)
	}
}

func loadTemplates(root string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			if path != root {
				_, err = loadTemplates(path)
				if err != nil {
					return err
				}
			}
		} else {
			files = append(files, path)
		}
		return err
	})
	return files, err
}
