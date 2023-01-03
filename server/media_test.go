package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/genjidb/genji"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_handler_mediaDetail(t *testing.T) {

	db, err := initGenji("/tmp/Test_handler_mediaDetail")
	assert.Nil(t, err)
	lg, err := zap.NewDevelopment()
	assert.Nil(t, err)
	defer os.RemoveAll("/tmp/Test_handler_mediaDetail")

	// db.Exec()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/medias/:id")
	c.SetParamNames("id")
	c.SetParamValues("12")

	type fields struct {
		db *genji.DB
		lg *zap.Logger
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "OK",
			fields:  fields{db, lg},
			args:    args{c},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler{
				db: tt.fields.db,
				lg: tt.fields.lg,
			}
			if err := h.mediaDetail(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("handler.mediaDetail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
