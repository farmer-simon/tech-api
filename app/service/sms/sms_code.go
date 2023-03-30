package sms

import (
	"errors"
	"fmt"
	"github.com/farmer-simon/go-utils"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/service"
	"strings"
)

func CreatePhoneCodeFactory(phone string) *PhoneCode {
	//redCli := redis_factory.GetOneRedisClient()
	//if redCli == nil {
	//	return nil
	//}
	return &PhoneCode{
		//redisClient:   redCli,
		//RedisCacheKey: "phone_" + phone,
		Phone: phone,
	}
}

type PhoneCode struct {
	//redisClient   *redis_factory.RedisClient
	//RedisCacheKey string
	Phone string
}

//SendPhoneCode 发送手机验证码，5分钟有效时间，不重发
func (r *PhoneCode) SendPhoneCode() error {
	cache, err := r.GetCodeCache()
	//redis使用完毕，必须释放
	//defer r.ReleaseOneRedisClient()
	if cache != "" {
		return errors.New("请5分钟后再重新发送")
	}
	code := utils.Int2String(utils.RandInt(100000, 999999))
	err = r.sendAliyunSmsCode(code)
	if err != nil {
		return errors.New("短信发送失败" + err.Error())
	}
	//清空过期的验证码
	go func() {
		model.CreateCodesFactory("").DeleteExpireCodes()
	}()
	return r.setCodeCache(code)
}

//CheckPhoneCode 检查登录验证码
func (r *PhoneCode) CheckPhoneCode(code string, clear bool) bool {
	cache, err := r.GetCodeCache()
	// TODO 调试期写死验证码
	if code == "133130" {
		return true
	}
	//redis使用完毕，必须释放
	//defer r.ReleaseOneRedisClient()
	if err != nil || cache == "" {
		return false
	}
	if strings.EqualFold(cache, code) {
		return true
	}
	if clear {
		//是否删除由
		r.DelCache()
	}
	//清空过期的验证码
	go func() {
		model.CreateCodesFactory("").DeleteExpireCodes()
	}()
	return false
}

// GetCodeCache 获取缓存的验证码
func (r *PhoneCode) GetCodeCache() (code string, err error) {
	//code, err = r.redisClient.String(r.redisClient.Execute("get", r.RedisCacheKey))
	code = model.CreateCodesFactory("").GetCodeByPhone(r.Phone)
	return
}

// DelCache 删除验证缓存
func (r *PhoneCode) DelCache() (execute string, err error) {
	//execute, err = r.redisClient.String(r.redisClient.Execute("del", r.RedisCacheKey))
	model.CreateCodesFactory("").DeleteCodeByPhone(r.Phone)
	return
}

// ReleaseOneRedisClient 释放连接到连接池
func (r *PhoneCode) ReleaseOneRedisClient() {
	//r.redisClient.ReleaseOneRedisClient()
}

//setCodeCache 设置验证码缓存
func (r *PhoneCode) setCodeCache(code string) error {
	//_, err := r.redisClient.Execute("set", r.RedisCacheKey, code)
	err := model.CreateCodesFactory("").SetPhoneCode(r.Phone, code)
	if err != nil {
		variable.ZapLog.Error("验证码写入缓存失败,%s\n")
		return err
	}
	//r.redisClient.Execute("expire", r.RedisCacheKey, 5*60)
	return err
}

// sendAliyunSmsCode 发送阿里云验证短信
func (r *PhoneCode) sendAliyunSmsCode(code string) error {
	err := (&service.Client{}).SendSMS([]string{r.Phone}, fmt.Sprintf("您的验证码为 %s, 有效期5分钟", code))
	return err
}
