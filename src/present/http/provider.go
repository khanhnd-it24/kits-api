package httpserver

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"kits/api/src/common/configs"
	"net"
	"net/http"
)

type HttpServer struct {
	srv  *http.Server
	Root *gin.RouterGroup
}

func NewHttpServer(cf *configs.Config) *HttpServer {
	if cf.Mode.IsDev() {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DisableConsoleColor()
	engine := gin.New()
	engine.Use(cors.AllowAll())
	r := engine.RouterGroup.Group(cf.Server.Prefix)
	srv := &http.Server{
		Addr:    cf.Server.Port,
		Handler: engine,
		//ReadTimeout:    10 * time.Second,
		//WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	httpServer := &HttpServer{Root: r, srv: srv}

	return httpServer
}

func (s *HttpServer) Start(ctx context.Context) error {
	addr := s.srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		_ = s.srv.Serve(ln)

	}()

	return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		wErr := fmt.Errorf("failed to stop http server: %w", err)
		return wErr
	}
	return nil
}
