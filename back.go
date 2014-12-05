package main

import (
    "fmt"
    "time"
    "net/http"
    "encoding/json"
    "github.com/garyburd/redigo/redis"
)

type Context struct {
    Token AccessToken
    CRs redis.Conn
}
type StatusError struct {
    code string
}
func (status StatusError) Error() string {
    return status.code
}
const (
    STATUS_SUCCESS = "SUCCES"
    STATUS_INVALID = "EINVAL"
)

type APIHandler struct {
    CRsPool *redis.Pool
    Routine func (
	ctx *Context,
	res http.ResponseWriter,
	req *http.Request,
    ) (interface{},error)
}
func (hlr APIHandler) ServeHTTP(res http.ResponseWriter,req *http.Request) {
    var err error

    if req.Method != "POST" {
	res.WriteHeader(404)
	return
    }

    ctx := &Context{}
    if ctx.CRs,err = hlr.CRsPool.Dial(); err != nil {
	res.WriteHeader(500)
	return
    }
    serial := ""
    if cookie,err := req.Cookie("serial"); err == nil {
	serial = cookie.Value
    }
    ctx.Token,err = AccessGetToken(ctx,serial)
    if err != nil {
	res.WriteHeader(500)
	return
    }

    ret := map[string]interface{}{}
    data,err := hlr.Routine(ctx,res,req)
    if err == nil {
	ret["status"] = STATUS_SUCCESS
	ret["data"] = data
    } else {
	if status,ok := err.(StatusError); ok {
	    ret["status"] = status.code
	    ret["data"] = nil
	} else {
	    fmt.Println(err)
	    res.WriteHeader(500)
	    return
	}
    }

    if json,err := json.Marshal(ret); err != nil {
	res.WriteHeader(500)
    } else {
	res.Header().Set("Content-Type","application/json")
	res.Write(json)
    }
}

func main() {
    crspl := &redis.Pool{
	MaxIdle:4,
	IdleTimeout:600 * time.Second,
	Dial:func() (redis.Conn,error) {
	    return redis.Dial("tcp","127.0.0.1:6379")
	},
    }

    http.Handle("/qa",APIHandler{crspl,RoutineQA})
    http.Handle("/login",APIHandler{crspl,RoutineLogin})

    http.Handle("/mg",APIHandler{crspl,RoutineMg})
    http.Handle("/mg/qa",APIHandler{crspl,RoutineMgQA})
    http.Handle("/mg/qa_add",APIHandler{crspl,RoutineMgQA_Add})

    http.ListenAndServe("127.0.0.1:3000",nil)
}
