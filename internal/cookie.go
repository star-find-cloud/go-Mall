package internal

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/star-find-cloud/star-mall/conf"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
)

type Cookie struct {
	secureCookie *securecookie.SecureCookie
}

func NewSecureCookie() *Cookie {
	// 使用 securecookie 库中的函数生成随机 64 位的密钥
	hashKey := securecookie.GenerateRandomKey(64)
	// 生成随机 32 位的密钥
	blockKey := securecookie.GenerateRandomKey(32)
	return &Cookie{
		secureCookie: securecookie.New(hashKey, blockKey),
	}
}

// 数据编码
func (c *Cookie) encodeValue(value interface{}) (string, error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return c.secureCookie.Encode("", jsonData)
}

// 数据解码
func (c *Cookie) decodeValue(encoded string, obj interface{}) error {
	var decoded []byte
	if err := c.secureCookie.Decode("", encoded, &decoded); err != nil {
		log.AppLogger.Errorf("http cookie error: %s", err)
		log.HttpLogger.Errorf("decode cookie value error: %s", err)
		return err
	}
	return json.Unmarshal(decoded, obj)
}

func (c *Cookie) Set(ctx *gin.Context, key string, value interface{}) {
	config := conf.GetConfig()
	encoded, err := c.encodeValue(value)
	if err != nil {
		log.AppLogger.Errorf("http cookie error: %s", err)
		log.HttpLogger.Errorf("encode cookie value error: %s", err)
		return
	}

	ctx.SetCookie(
		key,
		encoded,
		config.Cookie.MaxAge,
		config.Cookie.Path,
		config.Cookie.Domain,
		config.Cookie.Secure,
		config.Cookie.HttpOnly,
	)
}

func (c *Cookie) Delete(ctx *gin.Context, key string) {
	config := conf.GetConfig()
	ctx.SetCookie(
		key,
		"",
		-1,
		"/",
		config.Cookie.Domain,
		config.Cookie.Secure,
		config.Cookie.HttpOnly,
	)
}

func (c *Cookie) Get(ctx *gin.Context, key string, obj interface{}) bool {
	cookieValue, err := ctx.Cookie(key)
	if err != nil {
		log.AppLogger.Errorf("http cookie error: %s", err)
		log.HttpLogger.Errorf("get cookie error: %s", err)
		return false
	}
	return c.decodeValue(cookieValue, obj) == nil
}
