package runtime

import (
	"thirteen-admin/app/services"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// 进程级共享依赖（连接池与无请求状态组件）。
var (
	db      *gorm.DB
	rdb     *redis.Client
	jwt     *services.JWTManager
	online  *services.OnlineService
	captcha *services.CaptchaService
	file    *services.FileService
)

// SetDB 设置全局 DB 连接池。
func SetDB(v *gorm.DB) { db = v }

// DB 返回全局 DB 连接池。
func DB() *gorm.DB { return db }

// SetRedis 设置全局 Redis 客户端。
func SetRedis(v *redis.Client) { rdb = v }

// Redis 返回全局 Redis 客户端。
func Redis() *redis.Client { return rdb }

// SetJWT 设置 JWT 管理器。
func SetJWT(v *services.JWTManager) { jwt = v }

// JWT 返回 JWT 管理器。
func JWT() *services.JWTManager { return jwt }

// SetOnline 设置在线用户服务。
func SetOnline(v *services.OnlineService) { online = v }

// Online 返回在线用户服务。
func Online() *services.OnlineService { return online }

// SetCaptcha 设置验证码服务。
func SetCaptcha(v *services.CaptchaService) { captcha = v }

// Captcha 返回验证码服务。
func Captcha() *services.CaptchaService { return captcha }

// SetFile 设置文件服务。
func SetFile(v *services.FileService) { file = v }

// File 返回文件服务。
func File() *services.FileService { return file }
