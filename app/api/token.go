package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/coocood/freecache"
	"github.com/revel/revel"
)

const CACHE_SIZE = 10 * 1024 * 1024 // 10MB
var LocalCache *freecache.Cache

type Post struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Token   string `json:"app_access_token"`
	Expire  int    `json:"expire"`
}

func GetToken() string {

	var tokenstr string
	tokenByte, err := LocalCache.Get([]byte("app_access_token"))
	if err != nil {
		if err.Error() == "Entry not found" {
			tokenstr = refreshToken()
		} else {
			panic(err)
		}
	} else {
		tokenstr = string(tokenByte)
	}
	return tokenstr
}

func refreshToken() string {
	revel.AppLog.Debug("开始刷新飞书app_access_token")
	posturl := "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal"
	body := []byte(`{
		"app_id": "cli_a4d8741035f8d00d",
		"app_secret": "vn0yNEDf2XVxLzEWWMfj5fOMOwzcvweN"
	}`)

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	post := &Post{}
	derr := json.NewDecoder(res.Body).Decode(post)
	if derr != nil {
		panic(derr)
	}

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}
	token := post.Token
	LocalCache.Set([]byte("app_access_token"), []byte(token), post.Expire-30)
	revel.AppLog.Debug("成功刷新飞书app_access_token:", post)
	return token
}

func InitCache() {
	revel.AppLog.Debug("初始化本地缓存。")
	LocalCache = freecache.NewCache(CACHE_SIZE)
}
