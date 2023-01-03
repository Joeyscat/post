package main

import (
	"net/http"

	"github.com/genjidb/genji"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type handler struct {
	db *genji.DB
	lg *zap.Logger
}

func newHandler(db *genji.DB, lg *zap.Logger) handler {
	return handler{db: db, lg: lg}
}

type ErrorResponse struct {
	Msg string `json:"msg"`
}

func errorResp(msg string) *ErrorResponse {
	return &ErrorResponse{msg}
}

func internalServerError(c echo.Context, msg string) error {
	return c.JSON(http.StatusInternalServerError, errorResp(msg))
}

func badRequestError(c echo.Context, msg string) error {
	return c.JSON(http.StatusBadRequest, errorResp(msg))
}
