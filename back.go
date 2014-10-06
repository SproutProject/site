package main

import (
    "time"
    "net/http"
    "github.com/zenazn/goji/web"
    "github.com/zenazn/goji/web/middleware"
    "github.com/zenazn/goji/graceful"
    "github.com/garyburd/redigo/redis"
)

func main() {
    crspl := &redis.Pool{
	MaxIdle:4,
	IdleTimeout:600 * time.Second,
	Dial:func() (redis.Conn,error) {
	    return redis.Dial("tcp","127.0.0.1:6379")
	},
    }

    mux := web.New()
    mux.Use(middleware.EnvInit)
    mux.Use(middleware.Logger)
    mux.Use(middleware.NoCache)
    mux.Use(middleware.Recoverer)
    mux.Use(func(c *web.C,hlr http.Handler) http.Handler {
	mhlr := func(res http.ResponseWriter,req *http.Request) {
	    c.Env["crs"] = crspl.Get()
	    hlr.ServeHTTP(res,req)
	}
	return http.HandlerFunc(mhlr)
    })
    mux.Get("/access",RouteAccess)
    graceful.ListenAndServe("127.0.0.1:3000",mux)
}
