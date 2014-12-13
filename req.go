package main

import (
    "time"
    "net/http"
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "code.google.com/p/go.crypto/sha3"
    "github.com/garyburd/redigo/redis"
)

type Request struct {
    Id string
    Time string
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
	1,
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
	3,
    },
    {
`菊菊迫不及待要報名 2015 的資訊之芽了，請問哪項是他報名時需要知道的事情？`,
	[]string{
	    "在填寫報名表前，須先完成一份關於報名規則的小測驗，如果沒有全對就無法繼續報名的程序",
	    "有北區語法班、北區算法班、竹區語法班和竹區算法班四班可以選擇",
	    "有北區和竹區可以選擇，報名者可以根據自己的地理位置選擇到台灣大學或者交通大學上課",
	    "僅限高中生、高職生與五專等高中一到三年級同等學歷的學生報名",
	    "如果報名的班別為算法班，則除填寫報名表與備審資料外，仍須完成主辦單位提供的預試題才算完成報名手續",
	},
	4,
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
	    "阿哲：資訊之芽採取等第制度，有可能某兩人原始成績不同，最後拿到的證書卻是一樣的",
	    "阿英：資訊之芽這種東西就是先報名再說，反正覺得無聊就中途退出不去上課也不會有任何損失",
	    "阿義：資訊之芽是個殘酷的單位，如果全班作業成績都只有總分的 60% 不到，講師團隊將不排除死當全班",
	},
	2,
    },
    {
`資訊之芽算法班今年競爭激烈，最後恰只剩下一個正取名額，由下列五名報名者競爭。請問根據篩選規則最後應該獎落誰家？`,
	[]string{
	    "入芽考獲得滿分，備審資料也獲得極高評價，且年年都拿書卷獎的認真大學生小杉",
	    "入芽考靠著一只天外機器貓獲得滿分，但備審資料亂寫一通被打 0 分的高中生大雄",
	    "入芽考靠著實力勉強通過門檻，但備審資料用心填寫獲得一致肯定的正妹高中生靜香",
	    "入芽考沒有通過門檻，但備審資料震撼人心勇奪滿分的悲情派角色代表高中生小夫",
	    "入芽考獲得 0 分，備審資料慘不忍睹，但身為政府高官的老爸有來關說的高中生胖虎",
	},
	2,
    },
    {
`承上題，假如<b>除上題之正取者外</b>主辦單位要決定一個人作為備取人選，則根據篩選規則此人應當是誰？`,
	[]string{
	    "入芽考獲得滿分，備審資料也獲得極高評價，且年年都拿書卷獎的認真大學生小杉",
	    "入芽考靠著一只天外機器貓獲得滿分，但備審資料亂寫一通被打 0 分的高中生大雄",
	    "入芽考靠著實力勉強通過門檻，但備審資料用心填寫獲得一致肯定的正妹高中生靜香",
	    "入芽考沒有通過門檻，但備審資料震撼人心勇奪滿分的悲情派角色代表高中生小夫",
	    "入芽考獲得 0 分，備審資料慘不忍睹，但身為政府高官的老爸有來關說的高中生胖虎",
	},
	1,
    },
    {
`承上題，若備取人數決定擴增為兩人，則備取第二名根據篩選規則應當是誰？`,
	[]string{
	    "入芽考獲得滿分，備審資料也獲得極高評價，且年年都拿書卷獎的認真大學生小杉",
	    "入芽考靠著一只天外機器貓獲得滿分，但備審資料亂寫一通被打 0 分的高中生大雄",
	    "入芽考靠著實力勉強通過門檻，但備審資料用心填寫獲得一致肯定的正妹高中生靜香",
	    "入芽考沒有通過門檻，但備審資料震撼人心勇奪滿分的悲情派角色代表高中生小夫",
	    "入芽考獲得 0 分，備審資料慘不忍睹，但身為政府高官的老爸有來關說的高中生胖虎",
	},
	0,
    },
    {
`下列關於資訊之芽教學方式的敘述，何者正確？`,
	[]string{
	    "語法班偶爾請假無傷大雅，反正採翻轉課堂制，只要回家多看幾次教學影片即可",
	    "算法班偶爾請假無傷大雅，反正採翻轉課堂制，沒有到場只是沒辦法問講師問題",
	    "資訊之芽的作業很簡單，如果寫不出來一定是因為我不夠努力或是腦帶有洞，乾脆放棄這條路好了",
	    "資訊之芽和一般學校的課程一樣，著重知識的傳遞而非透過討論獲得答案",
	    "資訊之芽請假制度相當彈性，只要不超過一定次數都不會影響結業與否",
	},
	4,
    },
    {
`下列關於保證金制度的說明，何者<b>錯誤</b>？`,
	[]string{
	    "2015 年的保證金為 500 台幣",
	    "保證金必須在課程的第一週繳交，否則視同放棄錄取資格",
	    "如果第一階段結束後決定繼續參加第二階段，則不須再次繳交保證金，保證金將會在第二階段結束時退還",
	    "如果出席率未達參加證書門檻，則不論成績如何皆無法取得證書與保證金",
	    "如果出席率已達門檻但成績未達一般結業門檻，則僅退還 70% 的保證金",
	},
	4,
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
	time.Now().UTC().Format("2006-01-02 15:04:05 +0000"),
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
    /*for i := 0; i  < len(AlgPrePro); i += 1 {
	j := (weakrand.Int() % (len(AlgPrePro) - i)) + i
	AlgPrePro[i],AlgPrePro[j] = AlgPrePro[j],AlgPrePro[i]
    }*/
    return AlgPrePro,nil
}
