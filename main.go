package main

import (
	"context"
	"log"

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
	r.LoadHTMLGlob("./templates/*.html")

	err := sqlitedb.Connect()
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

	firebaseApp, err := firebaseservice.NewFirebaseApp("./firebase-sdk.json")
	if err != nil {
		panic(err)
	}

	firebaseAuthClient, err := firebaseservice.NewAuthService(context.Background(), firebaseApp)
	if err != nil {
		panic(err)
	}

	userService := controllers.NewUserController(userRepo)
	postControllrt := controllers.NewPostsController(postRepo)

	authMiddleware := middleware.NewAuthMiddleware(firebaseAuthClient, userRepo)

	r.Use(authMiddleware.ExtractUserFromToken())

	r.GET("/register", userService.CreateUser)
	r.POST("/register", authMiddleware.RegisterNewUser(), userService.CreateUser)
	r.GET("/login", userService.Login)
	r.POST("/login", userService.Login)
	r.GET("/", postControllrt.AllPostsByPage)
	r.GET("/home/:page", postControllrt.AllPostsByPage)
	r.GET("/:page", postControllrt.AllPostsByPage)
	r.GET("/new_post", authMiddleware.Authenticate(), postControllrt.NewPost)
	r.POST("/new_post", authMiddleware.Authenticate(), postControllrt.NewPost)
	r.Any("/post/:id", postControllrt.GetPost)
	r.GET("/edit_post/:id", authMiddleware.Authenticate(), postControllrt.EditPost)
	r.POST("/edit_post/:id", authMiddleware.Authenticate(), postControllrt.EditPost)
	r.GET("/delete_post/:id", authMiddleware.Authenticate(), postControllrt.Delete)

	err = r.Run()
	if err != nil {
		log.Panic(err)
	}
}
