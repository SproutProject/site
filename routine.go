package main

import (
    "fmt"
    "sort"
    "strconv"
    "io/ioutil"
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


func RoutineReqGetPre(
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    type Form struct {
	Clas int
	Recaptcha string
    }
    type RecaptchaResp struct {
	Success bool `json:"success"`
    }

    form := Form{}
    if err := json.Unmarshal(
	[]byte(req.PostFormValue("data")),
	&form,
    ); err != nil {
	return nil,StatusError{STATUS_INVALID}
    }

    url := fmt.Sprintf(
	"https://www.google.com/recaptcha/api/siteverify?secret=%s&response=%s",
	Recaptcha_Secret,
	form.Recaptcha,
    )
    resp,err := http.Get(url)
    if err != nil {
	return nil,StatusError{STATUS_INVALID}
    }
    defer resp.Body.Close()
    data,err := ioutil.ReadAll(resp.Body)
    if err != nil {
	return nil,StatusError{STATUS_INVALID}
    }
    recaptcha_resp := RecaptchaResp{}
    if err := json.Unmarshal(
	data,
	&recaptcha_resp,
    ); err != nil {
	return nil,StatusError{STATUS_INVALID}
    }
    if recaptcha_resp.Success == false {
	return nil,StatusError{STATUS_INVALID}
    }

    clas := 1
    if form.Clas == 0 {
	clas = 0
    }
    request,prepro,_ := ReqCreate(ctx,clas)

    http.SetCookie(res,&http.Cookie{
	Name: "req",
	Value: request.Id,
	Path: "/spt",
	MaxAge: 3600000,
	HttpOnly: true,
    })
    return prepro,nil
}
func RoutineReqCheckPre(
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    request,err := ReqLoad(ctx,req)
    if err != nil {
	return nil,StatusError{STATUS_INVALID}
    }
    ReqDel(ctx,request.Id)
    if request.Step != 0 {
	return nil,StatusError{STATUS_INVALID}
    }

    answer := []int{}
    if err := json.Unmarshal(
	[]byte(req.PostFormValue("data")),
	&answer,
    ); err != nil {
	return nil,StatusError{STATUS_INVALID}
    }
    if len(answer) != len(request.Answer) {
	return nil,StatusError{STATUS_INVALID}
    }
    for i,x := range request.Answer {
	if answer[i] != x {
	    return nil,StatusError{STATUS_INVALID}
	}
    }

    request.Step = 1
    ReqStore(ctx,&request)
    return nil,nil
}
func RoutineReqCheckMail(
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    request,err := ReqLoad(ctx,req)
    if err != nil {
	return nil,StatusError{STATUS_INVALID}
    }
    if request.Step != 1 {
	return nil,StatusError{STATUS_INVALID}
    }

    mail := req.PostFormValue("data")
    if request.Clas == 0 {
	resp,err := http.Get("http://reg.cms.sprout.csie.org/checker?" + mail)
	if err != nil {
	    return nil,StatusError{STATUS_INVALID}
	}
	defer resp.Body.Close()
	data,err := ioutil.ReadAll(resp.Body)
	if err != nil {
	    return nil,StatusError{STATUS_INVALID}
	}
	score,err := strconv.ParseInt(string(data),10,32)
	if err != nil {
	    return nil,StatusError{STATUS_INVALID}
	}
	if score < 250 {
	    return nil,StatusError{STATUS_INVALID}
	}
    }

    MailVerify(mail,request.Verify)
    request.Mail = mail

    request.Step = 2
    ReqStore(ctx,&request)
    return nil,nil
}
func RoutineReqVerify(
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    request,err := ReqLoad(ctx,req)
    if err != nil {
	return nil,StatusError{STATUS_INVALID}
    }
    if request.Step != 2 {
	return nil,StatusError{STATUS_INVALID}
    }

    code := req.PostFormValue("data")
    if request.Verify != code {
	return nil,StatusError{STATUS_INVALID}
    }

    request.Step = 3
    ReqStore(ctx,&request)
    return nil,nil
}
func RoutineReqData(
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    request,err := ReqLoad(ctx,req)
    if err != nil {
	return nil,StatusError{STATUS_INVALID}
    }
    if request.Step != 3 {
	return nil,StatusError{STATUS_INVALID}
    }

    request.Data = req.PostFormValue("data")

    request.Step = 4
    ReqStore(ctx,&request)
    if ReqDone(ctx,request.Id) != nil {
	return nil,StatusError{STATUS_INVALID}
    }
    return nil,nil
}

func RoutineMgReq (
    ctx *Context,
    res http.ResponseWriter,
    req *http.Request,
) (interface{},error) {
    if ctx.Token.Key == "" {
	return nil,StatusError{STATUS_INVALID}
    }
    ids,err := redis.Strings(ctx.CRs.Do("SMEMBERS","REQUEST_DONE"))
    if err != nil || len(ids) == 0 {
	return []interface{}{},nil
    }

    keys := redis.Args{}
    for i,_ := range(ids) {
	keys = keys.Add("REQUEST@" + ids[i])
    }
    datas,err := redis.Values(ctx.CRs.Do("MGET",keys...))
    if err != nil {
	return []interface{}{},err
    }

    requests := []Request{}
    for i,_ := range(ids) {
	if datas[i] == nil {
	    continue
	}
	request := Request{}
	json.Unmarshal(datas[i].([]byte),&request)
	requests = append(requests,request)
    }

    return requests,err
}
