package main

import (
    "net/http"
    weakrand "math/rand"
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
    {
`下列關於資訊之芽證書頒發的規定，何者為真？`,
	[]string{
	    "只要從頭到尾都沒有缺席就至少可以拿到結業證書",
	    "阿文就讀資訊之芽語法班，以正式學員的身份結束兩階段的課程並且成績達到結業的門檻，他從頭到尾應會拿到 2 張證書",
	    "阿文就讀資訊之芽語法班，非常上進努力，每次作業都有在時間以前全數完成並且拿到滿分，雖然在兩次認證考時都不幸狀況不佳獲得 0 分，還是有機會獲得優秀結業證書",
	    "阿文就讀資訊之芽算法班，非常上進努力，每次作業都有在時間以前全數完成並且拿到滿分，雖然在兩次認證考時都不幸狀況不佳獲得 0 分，還是有機會獲得優秀結業證書",
	    "阿文的雙胞胎弟弟阿又也參加了資訊之芽，雖然阿又不同於阿文，喜歡打混摸魚作業一次也沒繳，但只要每次都有到場還是可以繼續第二階段的課程",
	},
	2,
    },
    {
`時光來到 2016 年，阿金、阿文、阿哲、阿英、阿義都宣稱自己參加過 2015 的資訊之芽，根據他們的發言，請判斷誰才是真正參加過資訊之芽的人？`,
	[]string{
	    "阿金：資訊之芽就是個資訊競賽補習班，專門教人怎麼在資訊競賽中嶄露頭角",
	    "阿文：我在資訊之芽待了一個階段，讓頂尖教授上了八堂資訊專業課程獲益良多",
	    "阿哲：我參加資訊之芽到一半原本為了參選公職想要中途退出，卻發現只有在每階段結束時才有機會無痛領回保證金，只好把退出的計畫延後到第一階段結束",
	    "阿英：資訊之芽這種東西就是先報名再說，反正覺得無聊就中途退出不去上課也不會有任何損失",
	    "阿義：成績很糟糕也沒有關係，只要和老師混好關係成績曲線到最後就會轉彎",
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
func ReqDone(ctx *Context,id string) error {
    rnd,err := hex.DecodeString(id)
    if err != nil {
	return err
    }
    md := sha3.New512()
    md.Write([]byte(rnd))
    _,err = ctx.CRs.Do(
	"SADD",
	"REQUEST_DONE",hex.EncodeToString(md.Sum(nil)),
    )
    return err
}
func ReqGenPrePro(clas int) ([]ReqPrePro,error) {
    for i := 0; i  < len(AlgPrePro); i += 1 {
	j := (weakrand.Int() % (len(AlgPrePro) - i)) + i
	AlgPrePro[i],AlgPrePro[j] = AlgPrePro[j],AlgPrePro[i]
    }
    return AlgPrePro,nil
}
