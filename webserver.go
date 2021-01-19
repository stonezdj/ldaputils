package main

import (
	configModels "github.com/goharbor/ldaputils/dao/models"
	"github.com/goharbor/ldaputils/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func WebServer() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&configModels.LdapConfig{})
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.File("/", "public/index.html")
	e.File("/test", "public/test.html")
	e.File("/vue.js", "public/vue.js")
	e.File("/sample", "public/sample.html")
	e.GET("/configs", handlers.GetConfigs(db))
	e.PUT("/configs", handlers.PutConfig(db))
	e.DELETE("/configs/:id", handlers.DeleteConfig(db))
	e.POST("/testconfig/:id", handlers.TestingConfig(db))
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
