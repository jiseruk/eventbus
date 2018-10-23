package test

import (
	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app/model"
	"net/http"
	"net/http/httptest"
	"strings"
)

func getTopicMock(name string, engine string, resource string) *model.Topic{
	topic := model.Topic{ResourceID:resource, Name:name, Engine:engine}
	topic.CreatedAt = model.Clock.Now()
	topic.UpdatedAt = model.Clock.Now()
	return &topic
}

func executeMockedRequest(router *gin.Engine, method string, uri string, body string) *httptest.ResponseRecorder{
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(method, uri, strings.NewReader(body))
	router.ServeHTTP(rec, req)
	return rec
}
