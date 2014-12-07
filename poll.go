package main

import (
    "sort"
    "encoding/json"
    "code.google.com/p/go-uuid/uuid"
    "github.com/garyburd/redigo/redis"
)

type Poll struct {
    Id string
    Order int
    Subject string
    Body string
}
type PollList []Poll
func (polls PollList) Len() int {
    return len(polls)
}
func (polls PollList) Swap(i int,j int) {
    polls[i],polls[j] = polls[j],polls[i]
}
func (polls PollList) Less(i int,j int) bool {
    return polls[i].Order < polls[j].Order
}

func PollGetAll(ctx *Context) (interface{},error) {
    ids,err := redis.Strings(ctx.CRs.Do("SMEMBERS","POLL_ALL"))
    if err != nil || len(ids) == 0 {
	return []interface{}{},nil
    }

    keys := redis.Args{}
    for i,_ := range(ids) {
	keys = keys.Add("POLL@" + ids[i])
    }
    datas,err := redis.Values(ctx.CRs.Do("MGET",keys...))
    if err != nil {
	return []interface{}{},err
    }

    polls := PollList{}
    for i,_ := range(ids) {
	if datas[i] == nil {
	    continue
	}
	poll := Poll{}
	json.Unmarshal(datas[i].([]byte),&poll)
	polls = append(polls,poll)
    }
    sort.Sort(polls)

    return polls,nil
}
func PollAdd(ctx *Context,poll *Poll) error {
    if poll.Id == "" {
	if poll.Subject == "" {
	    return StatusError{STATUS_INVALID}
	}
	poll.Id = uuid.New()
    } else {
	poll.Id = uuid.Parse(poll.Id).String()
	if poll.Subject == "" {
	    if _,err := ctx.CRs.Do("DEL","POLL@" + poll.Id); err != nil {
		return err
	    }
	    if _,err := ctx.CRs.Do("SREM","POLL_ALL",poll.Id); err != nil {
		return err
	    }
	    return nil
	}
    }

    data,err := json.Marshal(poll)
    if err != nil {
	return err
    }
    if _,err := ctx.CRs.Do(
	"SET",
	"POLL@" + poll.Id,
	data,
    ); err != nil {
	return err
    }
    if _,err := ctx.CRs.Do("SADD","POLL_ALL",poll.Id); err != nil {
	return err
    }
    return nil
}
