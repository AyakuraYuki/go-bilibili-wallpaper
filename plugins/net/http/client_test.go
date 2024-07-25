package nhttp

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	cjson "github.com/AyakuraYuki/bilibili-wallpaper/plugins/json"
)

var ports = fmt.Sprintf(":%d", rand.Intn(10000)+30000)

func TestMain(m *testing.M) {
	engine := gin.New()
	engine.Use(gin.Recovery(), gin.Logger())
	engine.GET("/test-get", func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusOK, gin.H{"hello": "world"}) })
	engine.HEAD("/test-head", func(c *gin.Context) { c.AbortWithStatus(http.StatusNoContent) })
	engine.POST("/test-post", func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusOK, gin.H{"hello": "world"}) })
	engine.NoRoute(func(c *gin.Context) { c.AbortWithStatus(http.StatusNotFound) })
	server := &http.Server{Addr: ports, Handler: engine}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	defer func() { _ = server.Shutdown(context.Background()) }()

	m.Run()
}

func debugShowResponseResults(t *testing.T, data []byte, rspHeader http.Header, httpCode int) {
	t.Logf("http code: %v\n", httpCode)
	t.Logf("response data: %v\n", string(data))
	bs, _ := cjson.JSON.Marshal(rspHeader)
	t.Logf("headers: %v\n", string(bs))
}

func TestClient(t *testing.T) {
	t.Run("GetRaw", func(t *testing.T) {
		requestUrl := fmt.Sprintf("http://127.0.0.1%s/test-get", ports)
		data, header, code, err := GetRaw(nil, requestUrl, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		debugShowResponseResults(t, data, header, code)
	})
	t.Run("PostRaw", func(t *testing.T) {
		requestUrl := fmt.Sprintf("http://127.0.0.1%s/test-post", ports)
		data, header, code, err := PostRaw(nil, requestUrl, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		debugShowResponseResults(t, data, header, code)
	})
	t.Run("Head", func(t *testing.T) {
		requestUrl := fmt.Sprintf("http://127.0.0.1%s/test-head", ports)
		header, code, err := Head(nil, requestUrl, nil, 10000, 3)
		if err != nil {
			t.Fatal(err)
		}
		debugShowResponseResults(t, nil, header, code)
	})
}
