package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/genjidb/genji"
	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/types"
	"github.com/labstack/echo/v4"
)

func (h *handler) mediaDetail(c echo.Context) error {
	id := c.Param("id")

	db := h.db.WithContext(c.Request().Context())
	doc, err := db.QueryDocument("SELECT * FROM media WHERE id = ?", id)
	if err != nil {
		if genji.IsNotFoundError(err) {
			return c.JSON(http.StatusNotFound, nil)
		}

		return internalServerError(c, fmt.Sprintf("query media error: %s", err.Error()))
	}

	var media Media
	err = document.StructScan(doc, &media)
	if err != nil {
		return internalServerError(c, fmt.Sprintf("scan fields error: %s", err.Error()))
	}

	return c.JSON(http.StatusOK, media)
}

func (h *handler) mediaList(c echo.Context) error {
	type query struct {
		UserId *string `query:"userId"`
		Page   *uint   `query:"page"`
	}
	var q query
	err := c.Bind(&q)
	if err != nil {
		return badRequestError(c, fmt.Sprintf("bad request: %s", err.Error()))
	}
	if q.UserId == nil {
		return badRequestError(c, "userId is required")
	}

	db := h.db.WithContext(c.Request().Context())
	stream, err := db.Query("SELECT * FROM media WHERE userId = ?", q.UserId)
	if err != nil {
		return internalServerError(c, fmt.Sprintf("query media error: %s", err.Error()))
	}
	defer stream.Close()

	medias := []*Media{}
	err = stream.Iterate(func(d types.Document) error {
		var media Media
		err = document.StructScan(d, &media)
		if err != nil {
			return internalServerError(c, fmt.Sprintf("scan fields error: %s", err.Error()))
		}
		medias = append(medias, &media)

		return nil
	})

	return c.JSON(http.StatusOK, medias)
}

func (h *handler) mediaUpload(c echo.Context) error {
	user, err := getUserFromReq(c)
	if err != nil {
		return badRequestError(c, err.Error())
	}
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	if len(files) == 0 {
		return badRequestError(c, "file not found in form")
	}

	db := h.db.WithContext(c.Request().Context())
	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			return internalServerError(c, fmt.Sprintf("open file error: %s", err.Error()))
		}
		defer src.Close()

		ext := filepath.Ext(file.Filename)
		sid, err := shortId()
		if err != nil {
			return internalServerError(c, fmt.Sprintf("generate shortid error: %s", err.Error()))
		}
		newFilename := sid + ext
		// Destination
		dst, err := os.Create(fmt.Sprintf("/tmp/medias/%s", newFilename))
		if err != nil {
			return internalServerError(c, fmt.Sprintf("create new file error: %s", err.Error()))
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return internalServerError(c, fmt.Sprintf("copy file error: %s", err.Error()))
		}

		posted := false
		media := &Media{
			Id:     MediaId(sid),
			UserId: user.Id,
			Type:   mediaType(newFilename),
			URL:    fmt.Sprintf("http://localhost:8000/m/%s", newFilename),
			Posted: &posted,
			Time:   time.Now().Unix(),
		}

		err = db.Exec("INSERT INTO media VALUES ?", media)
		if err != nil {
			return internalServerError(c, fmt.Sprintf("INSERT media error: %s", err.Error()))
		}
	}

	return nil
}

func mediaType(path string) MediaType {
	ext := filepath.Ext(path)

	mt := MediaType_Unknown
	switch ext {
	case "jpg":
		mt = MediaType_Picture
	}

	return mt
}

func getUserFromReq(c echo.Context) (*User, error) {

	user := User{
		Id: DEV_USER_ID,
	}

	return &user, nil
}
