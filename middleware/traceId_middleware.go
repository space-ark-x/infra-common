package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/space-ark-x/infra-common/utils"
)

func NewTraceIdMiddleware() iris.Handler {
	return func(ctx iris.Context) {
		// 从Authorization头获取token
		traceId := ctx.GetHeader("X-Trace-Id")
		if traceId == "" {
			traceId = utils.GenUUID()
		}
		ctx.Values().Set("traceId", traceId)
		ctx.Next()
		ctx.Header("X-Trace-Id", traceId)
	}
}
