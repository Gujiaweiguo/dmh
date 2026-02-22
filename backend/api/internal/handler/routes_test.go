package handler

import (
	"testing"

	"dmh/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func TestRegisterHandlersSymbol(t *testing.T) {
	f := RegisterHandlers
	_ = f
}

func TestRegisterHandlers(t *testing.T) {
	server := rest.MustNewServer(rest.RestConf{Host: "127.0.0.1", Port: 0})
	defer server.Stop()

	ctx := &svc.ServiceContext{}
	ctx.Config.Auth.AccessSecret = "unit-test-secret"

	RegisterHandlers(server, ctx)
}
