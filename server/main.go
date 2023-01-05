package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	lg, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	dbPath := "/tmp/mygenji"
	// if os.Getenv("DEBUG") == "true" {
	// 	dbPath = ":memory:"
	// }
	db, err := initGenji(dbPath)
	if err != nil {
		lg.Fatal("failed to init genjidb", zap.String("err", err.Error()))
	}

	h := newHandler(db, lg)

	e := echo.New()
	e.Use(middleware.Logger())
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/posts", h.postCreate)
	e.GET("/posts", h.postList)
	e.DELETE("/posts/:id", h.postDelete)
	e.POST("/posts/:id/comment", h.postComment)

	e.POST("/medias/upload", h.mediaUpload)
	e.GET("/medias/:id", h.mediaDetail)
	e.GET("/medias", h.mediaList)

	e.Logger.Fatal(e.Start(":6000"))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
		c.HTML(http.StatusInternalServerError, fmt.Sprintf("<p>%s</p>", err.Error()))
	}
}
