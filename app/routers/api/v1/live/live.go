package live

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/app/controller/clive"
	"sports_service/global/app/errdef"
	"sports_service/global/app/log"
	"sports_service/models/mcontest"
)

func PushStreamCallback(c *gin.Context) {
	reply := errdef.New(c)
	param := &mcontest.StreamCallbackInfo{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("live_trace: invalid param, params:%+v", param)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	if param.EventType != 1 {
		log.Log.Errorf("live_trace: invalid evenType:%d", param.EventType)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	param.Status = 1
	svc := clive.New(c)
	code := svc.PushOrDisconnectStreamCallback(param)
	reply.Response(http.StatusOK, code)
}

func DisconnectStreamCallback(c *gin.Context) {
	reply := errdef.New(c)
	param := &mcontest.StreamCallbackInfo{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("live_trace: invalid param, params:%+v", param)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	if param.EventType != 0 {
		log.Log.Errorf("live_trace: invalid evenType:%d", param.EventType)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	param.Status = 2
	svc := clive.New(c)
	code := svc.PushOrDisconnectStreamCallback(param)
	reply.Response(http.StatusOK, code)
}

func TranscribeStreamCallback(c *gin.Context) {
	reply := errdef.New(c)
	param := &mcontest.StreamCallbackInfo{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("live_trace: invalid param, params:%+v", param)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	if param.EventType != 100 {
		log.Log.Errorf("live_trace: invalid evenType:%d", param.EventType)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := clive.New(c)
	code := svc.TranscribeStreamCallback(param)
	reply.Response(http.StatusOK, code)
}
