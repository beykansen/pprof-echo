# pprof-echo
Echo wrapper for pprof.
## How To Use?
```golang
	e := echo.New()
	pprofecho.NewPprofGroup(e, pprofecho.PprofGroupConfig{
		Prefix:  "",
		Skipper: func(c echo.Context) bool {
			return false
		},
	})
```
