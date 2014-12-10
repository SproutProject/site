package main

import (
    "time"
    "hash"
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "code.google.com/p/go.crypto/bcrypt"
    "code.google.com/p/go.crypto/sha3"
    "github.com/garyburd/redigo/redis"
)

const (
    Success = "SUCCES"
    EInval = "EINVAL"
)

type AccessToken struct {
    Key string
    Timestamp time.Time
}

func AccessGetToken(ctx *Context,serial string) (AccessToken,error) {
    token := AccessToken{
	Key: "",
    }

    if serial == "" {
	return token,nil
    }

    rnd,err := hex.DecodeString(serial)
    if err != nil {
	return token,nil
    }
    md := sha3.New512()
    md.Write([]byte(rnd))
    data,err := redis.Bytes(ctx.CRs.Do(
	"GET",
	"ACCESSTOKEN@" + hex.EncodeToString(md.Sum(nil)),
    ))
    if err := json.Unmarshal(data,&token); err != nil {
	return token,err
    }

    return token,nil
}
func AccessLogin(ctx *Context,mail string,passwd string) (string,error) {
    var md hash.Hash

    md = sha3.New512()
    md.Write([]byte(mail))
    key := hex.EncodeToString(md.Sum(nil))

    hash,err := redis.Bytes(ctx.CRs.Do("HGET","ACCT@" + key,"hash"))
    if err != nil {
	return "",StatusError{STATUS_INVALID}
    }
    if bcrypt.CompareHashAndPassword(hash,[]byte(passwd)) != nil {
	return "",StatusError{STATUS_INVALID}
    }

    rnd := make([]byte,512)
    if _,err := rand.Read(rnd); err != nil {
	return "",StatusError{STATUS_INVALID}
    }
    md = sha3.New512()
    md.Write([]byte(rnd))
    token := AccessToken {
	Key: hex.EncodeToString(md.Sum(nil)),
	Timestamp: time.Now().UTC(),
    }
    data,err := json.Marshal(token)
    if err != nil {
	return "",err
    }
    ctx.CRs.Do("SETEX","ACCESSTOKEN@" + token.Key,3600,data)

    return hex.EncodeToString(rnd),nil
}
