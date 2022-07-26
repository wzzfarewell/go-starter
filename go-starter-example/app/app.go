package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wzzfarewell/go-starter-example/delivery/http/middleware"
	"net/http"
	"time"
)

func Run(path string) error {
	if err := initConfigs(path); err != nil {
		return errors.Wrap(err, "init config failed")
	}
	if err := initRepositories(); err != nil {
		return errors.Wrap(err, "init repository failed")
	}
	if err := initServices(); err != nil {
		return errors.Wrap(err, "init service failed")
	}
	return serveHTTP()
}

func serveHTTP() error {
	e := gin.New()
	_ = e.SetTrustedProxies(nil)
	e.Use(gin.Recovery())
	e.NoRoute(middleware.NoRouteHandler())
	e.Use(middleware.ErrorHandler())
	cfg := configs.HTTPConfig
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      e,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
	}
	return s.ListenAndServe()
}
