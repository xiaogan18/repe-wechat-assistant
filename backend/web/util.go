package web

import (
	"errors"
	"github.com/xiaogan18/repe-wechat-assistant/backend/web/resp"
	"github.com/gin-gonic/gin"
	"strconv"
)

var errNotFound error = errors.New("query key not found")

func queryInt64(ctx *gin.Context, key string) (int64, error) {
	q := ctx.Query(key)
	if len(key) > 0 && len(q) > 0 {
		v, err := strconv.ParseInt(q, 10, 64)
		if err != nil {
			return 0, err
		}
		return v, nil
	}
	return 0, errNotFound
}
func queryString(ctx *gin.Context, key string) (string, error) {
	q := ctx.Query(key)
	if len(key) > 0 && len(q) > 0 {
		return q, nil
	}
	return q, errNotFound
}
func queryPage(ctx *gin.Context) resp.PageSearch {
	v := resp.PageSearch{
		Index: 0,
		Size:  20,
	}
	if index, err := queryInt64(ctx, "page_index"); err == nil {
		v.Index = int(index)
	}
	if size, err := queryInt64(ctx, "page_size"); err == nil {
		v.Size = int(size)
	}
	return v
}
