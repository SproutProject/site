package main

import (
    "sort"
    "net/http"
    "encoding/json"
    "code.google.com/p/go-uuid/uuid"
    "github.com/garyburd/redigo/redis"
)

func RoutineLogin (
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    mail := req.PostFormValue("mail")
    passwd := req.PostFormValue("passwd")

    serial,err := AccessLogin(ctx,mail,passwd)
    if err != nil {
	return nil,err
    }

    http.SetCookie(res,&http.Cookie{
	Name: "serial",
	Value: serial,
	Path: "/spt",
	MaxAge: 3600,
	HttpOnly: true,
    })
    return nil,nil
}

func RoutineMg (
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    if ctx.Token.Key == "" {
	return nil,StatusError{STATUS_INVALID}
    }
    return nil,nil
}

type QA struct {
    Id string
    Subject string
    Clas string
    Order int
    Body string
}
type QAList []QA
func (qas QAList) Len() int {
    return len(qas)
}
func (qas QAList) Swap(i int,j int) {
    qas[i],qas[j] = qas[j],qas[i]
}
func (qas QAList) Less(i int,j int) bool {
    return qas[i].Order < qas[j].Order
}
func RoutineQA (
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    ids,err := redis.Strings(ctx.CRs.Do("SMEMBERS","QA_ALL"))
    if err != nil || len(ids) == 0 {
	return []interface{}{},nil
    }

    keys := redis.Args{}
    for i,_ := range(ids) {
	keys = keys.Add("QA@" + ids[i])
    }
    datas,err := redis.Values(ctx.CRs.Do("MGET",keys...))
    if err != nil {
	return []interface{}{},err
    }

    qas := QAList{}
    for i,_ := range(ids) {
	if datas[i] == nil {
	    continue
	}
	qa := QA{}
	json.Unmarshal(datas[i].([]byte),&qa)
	qas = append(qas,qa)
    }
    sort.Sort(qas)

    return qas,nil
}
func RoutineMgQA (
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    if ctx.Token.Key == "" {
	return nil,StatusError{STATUS_INVALID}
    }

    ids,err := redis.Strings(ctx.CRs.Do("SMEMBERS","QA_ALL"))
    if err != nil || len(ids) == 0 {
	return []interface{}{},nil
    }

    keys := redis.Args{}
    for i,_ := range(ids) {
	keys = keys.Add("QA@" + ids[i])
    }
    datas,err := redis.Values(ctx.CRs.Do("MGET",keys...))
    if err != nil {
	return []interface{}{},err
    }

    qas := QAList{}
    for i,_ := range(ids) {
	if datas[i] == nil {
	    continue
	}
	qa := QA{}
	json.Unmarshal(datas[i].([]byte),&qa)
	qas = append(qas,qa)
    }
    sort.Sort(qas)

    return qas,nil
}
func RoutineMgQA_Add (
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    if ctx.Token.Key == "" {
	return nil,StatusError{STATUS_INVALID}
    }

    qa := QA{}
    if err := json.Unmarshal(
	[]byte(req.PostFormValue("data")),
	&qa,
    ); err != nil {
	return nil,err
    }

    if qa.Id == "" {
	if qa.Subject == "" {
	    return nil,StatusError{STATUS_INVALID}
	}
	qa.Id = uuid.New()
    } else {
	qa.Id = uuid.Parse(qa.Id).String()
	if qa.Subject == "" {
	    if _,err := ctx.CRs.Do("DEL","QA@" + qa.Id); err != nil {
		return nil,err
	    }
	    if _,err := ctx.CRs.Do("SREM","QA_ALL",qa.Id); err != nil {
		return nil,err
	    }
	    return nil,nil
	}
    }

    data,err := json.Marshal(qa)
    if err != nil {
	return nil,err
    }
    if _,err := ctx.CRs.Do(
	"SET",
	"QA@" + qa.Id,
	data,
    ); err != nil {
	return nil,err
    }
    if _,err := ctx.CRs.Do("SADD","QA_ALL",qa.Id); err != nil {
	return nil,err
    }

    return nil,nil
}

func RoutinePoll (
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    return PollGetAll(ctx)
}
func RoutineMgPoll (
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    if ctx.Token.Key == "" {
	return nil,StatusError{STATUS_INVALID}
    }

    return PollGetAll(ctx)
}
func RoutineMgPoll_Add (
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    if ctx.Token.Key == "" {
	return nil,StatusError{STATUS_INVALID}
    }

    poll := Poll{}
    if err := json.Unmarshal(
	[]byte(req.PostFormValue("data")),
	&poll,
    ); err != nil {
	return nil,err
    }

    return nil,PollAdd(ctx,&poll)
}

type Request struct {
    Id string
    Name string
    Mail string
    School string
    Phone string
}
