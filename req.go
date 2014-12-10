package main

import (
    "strconv"
    "net/http"
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "code.google.com/p/go.crypto/sha3"
    "github.com/garyburd/redigo/redis"
)

type Request struct {
    Id string
    Step int
    Clas int
    Answer []int
    Mail string
    Verify string
    Data interface{}
}
type ReqPrePro struct {
    Statement string
    Option []string
    Answer int `json:"-"`
}

var AlgPrePro = []ReqPrePro{
    {
`阿文參加了資訊之芽，下列哪些不是他可能從中學到的技能？
<ol>
    <li>了解電腦的運作邏輯</li>
    <li>針對問題設計高效的解決辦法</li>
    <li>開發手機 APP</li>
    <li>製作充滿 3D 特效的線上遊戲</li>
    <li>了解資訊系不是會寫程式就好</li>
</ol>`,
	[]string{
	    "1,2",
	    "3,4",
	    "1,5",
	    "1,2,3,5",
	    "2,3,4",
	},
	2,
    },
    {
`資訊之芽除了上課時間以外，歷屆學員平均每週大約需要額外花多少時間？`,
	[]string{
	    "0-4 hrs",
	    "4-8 hrs",
	    "8-12 hrs",
	    "12-16 hrs",
	    "16-20 hrs",
	},
	2,
    },
}

func ReqCreate(ctx *Context,clas int) (Request,[]ReqPrePro,error) {
    prolist,_ := ReqGenPrePro(clas)
    answer := make([]int,0,10)
    for _,pro := range prolist {
	answer = append(answer,pro.Answer)
    }
    req := Request{
	"",
	0,
	clas,
	answer,
	"",
	"",
	[]interface{}{},
    }

    rnd := make([]byte,512)
    if _,err := rand.Read(rnd); err != nil {
	return req,prolist,StatusError{STATUS_INVALID}
    }
    req.Id = hex.EncodeToString(rnd)

    rnd = make([]byte,3)
    if _,err := rand.Read(rnd); err != nil {
	return req,prolist,StatusError{STATUS_INVALID}
    }
    req.Verify = hex.EncodeToString(rnd)

    if ReqStore(ctx,&req) != nil {
	return req,prolist,StatusError{STATUS_INVALID}
    }

    return req,prolist,nil
}
func ReqStore(ctx *Context,req *Request) error {
    rnd,err := hex.DecodeString(req.Id)
    if err != nil {
	return err
    }
    md := sha3.New512()
    md.Write([]byte(rnd))
    data,err := json.Marshal(req)
    if err != nil {
	return err
    }
    _,err = ctx.CRs.Do(
	"SETEX",
	"REQUEST@" + hex.EncodeToString(md.Sum(nil)),
	3600000,
	data,
    )
    return err
}
func ReqLoad(ctx *Context,req *http.Request) (Request,error) {
    request := Request{}

    cookie,err := req.Cookie("req")
    if err != nil {
	return request,err
    }

    rnd,err := hex.DecodeString(cookie.Value)
    if err != nil {
	return request,nil
    }
    md := sha3.New512()
    md.Write([]byte(rnd))
    data,err := redis.Bytes(ctx.CRs.Do(
	"GET",
	"REQUEST@" + hex.EncodeToString(md.Sum(nil)),
    ))
    if err := json.Unmarshal(data,&request); err != nil {
	return request,nil
    }
    return request,nil
}
func ReqDel(ctx *Context,id string) error {
    rnd,err := hex.DecodeString(id)
    if err != nil {
	return err
    }
    md := sha3.New512()
    md.Write([]byte(rnd))
    _,err = ctx.CRs.Do(
	"DEL",
	"REQUEST@" + hex.EncodeToString(md.Sum(nil)),
    )
    return err
}
func ReqGenPrePro(clas int) ([]ReqPrePro,error) {
    return AlgPrePro,nil
}
func ReqHashPrePro(selec []int) string {
    ans := "xfzl3(E)qEU,WO,AUWE09,uOPSUAS80A98D0_"
    for _,option := range selec {
	ans += "_" + strconv.Itoa(option)
    }
    md := sha3.New512()
    md.Write([]byte(ans))
    return hex.EncodeToString(md.Sum(nil))
}
