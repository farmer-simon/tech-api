package model

import (
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/users/token_cache_redis"
	"time"
)

func CreateOauthFactory(sqlType string) *OauthModel {
	return &OauthModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type OauthModel struct {
	BaseModel
	Platform  string `gorm:"column:platform" json:"platform"`
	FrUserId  int64  `gorm:"column:fr_user_id" json:"fr_user_id"`
	Token     string `gorm:"column:token" json:"token"`
	ExpiresAt int64  `gorm:"column:expires_at" json:"expires_at"`
}

func (o *OauthModel) TableName() string {
	return "tech_oauth_access_tokens"
}

//RecordOauthToken 记录用户Token
func (o *OauthModel) RecordOauthToken(platform string, userId int64, token string, expiresAt int64) bool {
	sql := `INSERT   INTO  tech_oauth_access_tokens(platform, fr_user_id, token, expires_at) VALUES(?, ?, ?, ?)`
	//注意：token的精确度为秒，如果在一秒之内，一个账号多次调用接口生成的token其实是相同的，这样写入数据库，第二次的影响行数为0，知己实际上操作仍然是有效的。
	//所以这里只判断无错误即可，判断影响行数的话，>=0 都是ok的
	if o.Exec(sql, platform, userId, token, time.Unix(expiresAt, 0).Format(variable.DateFormat)).Error == nil {
		// 异步缓存用户有效的token到redis
		if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
			go o.ValidTokenCacheToRedis(platform, userId)
		}
		// 清除当前用户失效Token
		go o.delOverMaxOnlineToken(platform, userId)
		return true
	}
	return false
}

//OauthRefreshConditionCheck 用户刷新token,条件检查: 相关token在过期的时间之内，就符合刷新条件
func (o *OauthModel) OauthRefreshConditionCheck(platform string, userId int64, oldToken string) bool {
	// 首先判断旧token在本系统自带的数据库已经存在，才允许继续执行刷新逻辑
	var oldTokenIsExists int
	sql := "SELECT count(*)  as  counts FROM tech_oauth_access_tokens  WHERE platform =? and fr_user_id =? and token=? and NOW()<DATE_ADD(expires_at,INTERVAL  ? SECOND)"
	if o.Raw(sql, platform, userId, oldToken, variable.ConfigYml.GetInt64("Token.JwtTokenRefreshAllowSec")).First(&oldTokenIsExists).Error == nil && oldTokenIsExists == 1 {
		return true
	}
	return false
}

//OauthRefreshToken 用户刷新token
func (o *OauthModel) OauthRefreshToken(userId, expiresAt int64, oldToken, newToken, platform string) bool {
	sql := "UPDATE   tech_oauth_access_tokens   SET  token=? ,expires_at=?,platform=?,updated_at=NOW()  WHERE   fr_user_id=? AND token=?"
	if o.Exec(sql, newToken, time.Unix(expiresAt, 0).Format(variable.DateFormat), platform, userId, oldToken).Error == nil {
		// 异步缓存用户有效的token到redis
		if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
			go o.ValidTokenCacheToRedis(platform, userId)
		}
		// 清除当前用户失效Token
		go o.delOverMaxOnlineToken(platform, userId)
		return true
	}
	return false
}

//Destroy 删除用户以及关联的token记录
func (o *OauthModel) Destroy(id int, platform string) bool {

	// 删除用户时，清除用户缓存在redis的全部token
	if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
		go o.DelTokenCacheFromRedis(int64(id), platform)
	}
	sql := "DELETE FROM  tech_oauth_access_tokens WHERE  fr_user_id=? AND platform=? "
	//判断>=0, 有些没有登录过的用户没有相关token，此语句执行影响行数为0，但是仍然是执行成功
	if o.Exec(sql, id, platform).Error == nil {
		return true
	}
	return false
}

//delOverMaxOnlineToken 删除超过最大在线人数Token和过期Token
func (o *OauthModel) delOverMaxOnlineToken(platform string, userId int64) bool {
	//删除此用户已过期的Token
	sql := "DELETE FROM  tech_oauth_access_tokens WHERE  fr_user_id=? AND platform=? AND expires_at<=NOW() "
	if o.Exec(sql, userId, platform).Error != nil {
		return false
	}
	var tokenCounts int64
	o.Model(o).Where("fr_user_id", userId).Where("platform", platform).Count(&tokenCounts)
	maxOnlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	if int(tokenCounts) > maxOnlineUsers {
		sql = "DELETE FROM  `tech_oauth_access_tokens` WHERE fr_user_id=? AND platform=? ORDER BY expires_at ASC , updated_at ASC LIMIT ?"
		//判断>=0, 有些没有登录过的用户没有相关token，此语句执行影响行数为0，但是仍然是执行成功
		if o.Exec(sql, userId, platform, int(tokenCounts)-maxOnlineUsers).Error != nil {
			return false
		}
	}

	return true
}

// OauthCheckTokenIsOk 判断用户token是否在数据库存在+状态OK
func (o *OauthModel) OauthCheckTokenIsOk(userId int64, token, platform string) bool {
	sql := "SELECT   token  FROM  `tech_oauth_access_tokens`  WHERE   fr_user_id=?  AND  platform=?  AND  expires_at>NOW() ORDER  BY  expires_at  DESC , updated_at  DESC  LIMIT ?"
	maxOnlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	rows, err := o.Raw(sql, userId, platform, maxOnlineUsers).Rows()
	defer func() {
		//  凡是查询类记得释放记录集
		_ = rows.Close()
	}()
	if err == nil && rows != nil {
		for rows.Next() {
			var tempToken string
			err := rows.Scan(&tempToken)
			if err == nil {
				if tempToken == token {
					return true
				}
			}
		}
	}
	return false
}

// ValidTokenCacheToRedis 后续两个函数专门处理用户 token 缓存到 redis 逻辑
func (o *OauthModel) ValidTokenCacheToRedis(platform string, userId int64) {
	tokenCacheRedisFact := token_cache_redis.CreateUsersTokenCacheFactory(userId, platform)
	if tokenCacheRedisFact == nil {
		variable.ZapLog.Error("redis连接失败，请检查配置")
		return
	}
	defer tokenCacheRedisFact.ReleaseRedisConn()

	sql := "SELECT   token,expires_at  FROM  `tech_oauth_access_tokens`  WHERE   fr_user_id=?  AND  platform=?  AND  expires_at>NOW() ORDER  BY  expires_at  DESC , updated_at  DESC  LIMIT ?"
	maxOnlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	rows, err := o.Raw(sql, userId, platform, maxOnlineUsers).Rows()
	defer func() {
		//  凡是获取原生结果集的查询，记得释放记录集
		_ = rows.Close()
	}()

	var tempToken, expires string
	if err == nil && rows != nil {
		for i := 1; rows.Next(); i++ {
			err = rows.Scan(&tempToken, &expires)
			if err == nil {
				if ts, err := time.ParseInLocation(variable.DateFormat, expires, time.Local); err == nil {
					tokenCacheRedisFact.SetTokenCache(ts.Unix(), tempToken)
					// 因为每个用户的token是按照过期时间倒叙排列的，第一个是有效期最长的，将该用户的总键设置一个最大过期时间，到期则自动清理，避免不必要的数据残留
					if i == 1 {
						tokenCacheRedisFact.SetUserTokenExpire(ts.Unix())
					}
				} else {
					variable.ZapLog.Error("expires_at 转换位时间戳出错", zap.Error(err))
				}
			}
		}
	}
	// 缓存结束之后删除超过系统设置最大在线数量的token
	go tokenCacheRedisFact.DelOverMaxOnlineCache()
}

// DelTokenCacheFromRedis 用户密码修改后，删除redis所有的token
func (o *OauthModel) DelTokenCacheFromRedis(userId int64, platform string) {
	tokenCacheRedisFact := token_cache_redis.CreateUsersTokenCacheFactory(userId, platform)
	if tokenCacheRedisFact == nil {
		variable.ZapLog.Error("redis连接失败，请检查配置")
		return
	}
	tokenCacheRedisFact.ClearUserToken()
	tokenCacheRedisFact.ReleaseRedisConn()
}
