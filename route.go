package main

import (
    "net/http"
    "encoding/json"
    "github.com/zenazn/goji/web"
    "github.com/garyburd/redigo/redis"
)

type AccessToken struct {
    ID int
    Name string
    Token string
}

func RouteAccess(c web.C,res http.ResponseWriter,req *http.Request) {
    crs := c.Env["crs"].(redis.Conn)
    if _,err := crs.Do("GET","TOKEN@"); err != nil {
	res.WriteHeader(1000)
	return
    }

    token := AccessToken{
	ID:0,
	Name:"NoAmI",
	Token:"DEADBEEF",
    }
    if out,err := json.Marshal(token); err == nil {
	res.Write(out)
    } else {
	res.WriteHeader(500)
    }
}
