package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/genjidb/genji"
	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/types"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *handler) postList(c echo.Context) error {
	type query struct {
		UserId *string `query:"userId"`
		Page   *uint   `query:"page"`
	}
	var q query
	err := c.Bind(&q)
	if err != nil {
		return badRequestError(c, fmt.Sprintf("bad request: %s", err.Error()))
	}

	db := h.db.WithContext(c.Request().Context())

	sql := "SELECT * FROM post"
	if q.UserId != nil {
		sql = fmt.Sprintf("%s WHERE userId = '%s'", sql, *q.UserId)
	}

	stream, err := db.Query(sql, q.UserId)
	if err != nil {
		return internalServerError(c, fmt.Sprintf("query post error: %s", err.Error()))
	}
	defer stream.Close()

	posts := []*Post{}
	err = stream.Iterate(func(d types.Document) error {
		var post Post
		err = document.StructScan(d, &post)
		if err != nil {
			return internalServerError(c, fmt.Sprintf("scan fields error: %s", err.Error()))
		}
		posts = append(posts, &post)

		return nil
	})

	return c.JSON(http.StatusOK, posts)
}

func (h *handler) postCreate(c echo.Context) error {
	user, err := getUserFromReq(c)
	if err != nil {
		return badRequestError(c, err.Error())
	}

	type params struct {
		Content *string  `json:"content"`
		Pics    []string `json:"pics"`
		Video   *string  `json:"video"`
	}
	var p params
	err = c.Bind(&p)
	if err != nil {
		return badRequestError(c, fmt.Sprintf("bad request: %s", err.Error()))
	}

	db := h.db.WithContext(c.Request().Context())

	sid, err := shortId()
	if err != nil {
		return internalServerError(c, fmt.Sprintf("generate shortid error: %s", err.Error()))
	}
	post := Post{
		Id:     PostId(sid),
		UserId: user.Id,
		Medias: []*MediaId{},
		Time:   time.Now().Unix(),
	}
	h.lg.Debug("creating post", zap.Any("post", post))

	if p.Content != nil {
		post.Content = *p.Content
	} else {
		if len(p.Pics) == 0 && p.Video == nil {
			return badRequestError(c, "empty post")
		}
	}
	if len(p.Pics) != 0 {
		if p.Video != nil {
			return badRequestError(c, "not supported post")
		}
		for _, v := range p.Pics {
			// TODO 确认所有Picture有效
			_, err = db.QueryDocument("SELECT id FROM media WHERE id = ?", v)
			if err != nil {
				return badRequestError(c, fmt.Sprintf("query media error: %s", err.Error()))
			}

			mid := MediaId(v)
			post.Medias = append(post.Medias, &mid)
		}
	}
	if p.Video != nil {
		// TODO 确认Video有效
		_, err = db.QueryDocument("SELECT id FROM media WHERE id = ?", p.Video)
		if err != nil {
			return badRequestError(c, fmt.Sprintf("query media error: %s", err.Error()))
		}

		mid := MediaId(*p.Video)
		post.Medias = append(post.Medias, &mid)
	}

	// post插入，media更新需要在一个事务中完成
	tx, err := db.Begin(true)
	if err != nil {
		return internalServerError(c, fmt.Sprintf("starts new transaction error: %s", err.Error()))
	}
	defer tx.Rollback()

	err = tx.Exec("INSERT INTO post VALUES ?", post)
	if err != nil {
		return internalServerError(c, fmt.Sprintf("insert post error: %s", err.Error()))
	}
	err = tx.Exec("UPDATE media SET posted = true WHERE id IN (?)", post.Medias)
	if err != nil {
		return internalServerError(c, fmt.Sprintf("update media error: %s", err.Error()))
	}

	err = tx.Commit()
	if err != nil {
		return internalServerError(c, fmt.Sprintf("commit error: %s", err.Error()))
	}

	return c.JSON(http.StatusCreated, nil)
}

func (h *handler) postComment(c echo.Context) error {
	user, err := getUserFromReq(c)
	if err != nil {
		return badRequestError(c, err.Error())
	}

	id := c.Param("id")
	type params struct {
		Content *string `json:"content"`
	}
	var p params
	err = c.Bind(&p)
	if err != nil {
		return badRequestError(c, fmt.Sprintf("bad request: %s", err.Error()))
	}

	newComment := Comment{
		UserId:  user.Id,
		Content: *p.Content,
		Time:    time.Now().Unix(),
	}
	h.lg.Debug("new comment for post", zap.Any("comment", newComment), zap.String("postId", id))

	db := h.db.WithContext(c.Request().Context())

	// post查询，comment更新需要在一个事务中完成
	tx, err := db.Begin(true)
	if err != nil {
		return internalServerError(c, fmt.Sprintf("starts new transaction error: %s", err.Error()))
	}
	defer tx.Rollback()

	doc, err := tx.QueryDocument("SELECT * FROM post WHERE id = ?", id)
	if err != nil {
		if genji.IsNotFoundError(err) {
			return c.JSON(http.StatusNotFound, nil)
		}
		return internalServerError(c, fmt.Sprintf("query post error: %s", err.Error()))
	}

	var post Post
	err = document.StructScan(doc, &post)
	if err != nil {
		return internalServerError(c, fmt.Sprintf("scan fields error: %s", err.Error()))
	}

	post.Comments = append(post.Comments, &newComment)
	err = tx.Exec("UPDATE post SET comments = ? WHERE id = ?", post.Comments, post.Id)
	if err != nil {
		return internalServerError(c, fmt.Sprintf("append comment error: %s", err.Error()))
	}

	err = tx.Commit()
	if err != nil {
		return internalServerError(c, fmt.Sprintf("commit error: %s", err.Error()))
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *handler) postDelete(c echo.Context) error {
	user, err := getUserFromReq(c)
	if err != nil {
		return badRequestError(c, err.Error())
	}

	id := c.Param("id")
	h.lg.Debug("deleting post", zap.Any("postId", id))

	db := h.db.WithContext(c.Request().Context())
	_, err = db.QueryDocument("SELECT id FROM post WHERE id = ? AND userId = ?", id, user.Id)
	if err != nil {
		return badRequestError(c, err.Error())
	}

	err = db.Exec("DELETE FROM post WHERE id = ? AND userId = ?", id, user.Id)
	if err != nil {
		return internalServerError(c, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}
