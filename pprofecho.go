package pprofecho

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	rtp "runtime/pprof"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type PprofGroupConfig struct {
	Prefix  string
	Skipper middleware.Skipper
}

func pprofMiddleware(skipper middleware.Skipper) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return echo.NewHTTPError(404)
			}
			return next(c)
		}
	}
}

func NewPprofGroup(e *echo.Echo, config PprofGroupConfig) {
	if len(strings.TrimSpace(config.Prefix)) == 0 {
		config.Prefix = "/debug/pprof"
	}
	g := e.Group(config.Prefix)
	g.Use(pprofMiddleware(config.Skipper))

	g.GET("/", newWrapperHandler(pprof.Index))
	g.GET("/cmdline", newWrapperHandler(pprof.Cmdline))
	g.GET("/profile", newWrapperHandler(pprof.Profile))
	g.GET("/symbol", newWrapperHandler(pprof.Symbol))
	g.POST("/symbol", newWrapperHandler(pprof.Symbol))
	g.GET("/trace", newWrapperHandler(pprof.Trace))

	for _, v := range rtp.Profiles() {
		ppName := v.Name()
		g.GET(fmt.Sprintf("/%s", ppName), newWrapperHandler(pprof.Handler(ppName).ServeHTTP))
	}
}

func newWrapperHandler(handlerFunc func(w http.ResponseWriter, r *http.Request)) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		handlerFunc(ctx.Response().Writer, ctx.Request())
		return nil
	}
}
