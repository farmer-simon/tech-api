package routers

import (
	"github.com/farmer-simon/go-utils"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/controller/api/admin"
	"goskeleton/app/http/controller/api/home"
	"goskeleton/app/http/controller/captcha"
	"goskeleton/app/http/middleware/authorization"
	"goskeleton/app/http/middleware/cors"
	validatorFactory "goskeleton/app/http/validator/core/factory"
	"goskeleton/app/utils/gin_release"
	"time"
)

// 该路由主要设置门户类网站等前台路由

func InitApiRouter() *gin.Engine {
	var router *gin.Engine
	// 非调试模式（生产模式） 日志写到日志文件
	if variable.ConfigYml.GetBool("AppDebug") == false {
		//1.gin自行记录接口访问日志，不需要nginx，如果开启以下3行，那么请屏蔽第 34 行代码
		//gin.DisableConsoleColor()
		//f, _ := os.Create(variable.BasePath + variable.ConfigYml.GetString("Logs.GinLogName"))
		//gin.DefaultWriter = io.MultiWriter(f)

		//【生产模式】
		// 根据 gin 官方的说明：[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
		// 如果部署到生产环境，请使用以下模式：
		// 1.生产模式(release) 和开发模式的变化主要是禁用 gin 记录接口访问日志，
		// 2.go服务就必须使用nginx作为前置代理服务，这样也方便实现负载均衡
		// 3.如果程序发生 panic 等异常使用自定义的 panic 恢复中间件拦截、记录到日志
		router = gin_release.ReleaseRouter()
	} else {
		// 调试模式，开启 pprof 包，便于开发阶段分析程序性能
		router = gin.Default()
		pprof.Register(router)
	}
	// 设置可信任的代理服务器列表,gin (2021-11-24发布的v1.7.7版本之后出的新功能)
	if variable.ConfigYml.GetInt("HttpServer.TrustProxies.IsOpen") == 1 {
		if err := router.SetTrustedProxies(variable.ConfigYml.GetStringSlice("HttpServer.TrustProxies.ProxyServerList")); err != nil {
			variable.ZapLog.Error(consts.GinSetTrustProxyError, zap.Error(err))
		}
	} else {
		_ = router.SetTrustedProxies(nil)
	}

	//根据配置进行设置跨域
	if variable.ConfigYml.GetBool("HttpServer.AllowCrossDomain") {
		router.Use(cors.Next())
	}

	router.GET("/", func(context *gin.Context) {
		//context.String(http.StatusOK, "Api 模块接口 hello word！")
		context.JSON(200, gin.H{
			"version":  "20201023.1508",
			"datetime": time.Now().Format(utils.TimeFormat),
		})
	})

	//处理静态资源（不建议gin框架处理静态资源，参见 Public/readme.md 说明 ）
	router.Static("/public", "./public")             //  定义静态资源路由与实际目录映射关系
	router.StaticFile("/abcd", "./public/readme.md") // 可以根据文件名绑定需要返回的文件名

	front := router.Group("/api/v1")
	{
		passport := front.Group("/passport")
		{
			//注册
			passport.POST("/register", validatorFactory.Create(consts.ValidatorPrefix+"Register"))
			//登录
			passport.POST("/login", validatorFactory.Create(consts.ValidatorPrefix+"PhoneLogin"))

			//passport.POST("/logout", (&home.Passport{}).PhoneCode)
			//
			passport.POST("/code", (&home.Passport{}).PhoneCode)
			passport.GET("/sso", (&home.Passport{}).Sso)

		}

		//需要校验登录的部分
		front.Use(authorization.BindContextAuthToken())
		{
			//主页
			front.GET("/index", (&home.Home{}).Index)
			front.GET("/category", (&home.Home{}).Category)
			services := front.Group("/services")
			{
				services.GET("/index", validatorFactory.Create(consts.ValidatorPrefix+"ServicesIndex"))
				services.GET("/info", (&home.Services{}).Info)
			}
			needs := front.Group("/needs")
			{
				needs.GET("/index", validatorFactory.Create(consts.ValidatorPrefix+"NeedIndex"))
				needs.GET("/info", (&home.Needs{}).Info)
			}

			front.Use(authorization.CheckTokenAuth(), authorization.CheckAccessToApi())
			{
				//contact := front.Group("/contact")
				//{
				//	contact.POST("/phone_code", (&home.Members{}).GetPhoneCode)
				//	contact.POST("/services", validatorFactory.Create(consts.ValidatorPrefix+"ContactServices"))
				//	contact.POST("/needs", validatorFactory.Create(consts.ValidatorPrefix+"ContactNeeds"))
				//}
				post := front.Group("/post")
				{
					post.POST("/services", validatorFactory.Create(consts.ValidatorPrefix+"ServicesRecordPost"))
					post.POST("/needs", validatorFactory.Create(consts.ValidatorPrefix+"NeedsRecordPost"))
				}

				//用户中心相关
				my := front.Group("/center")
				{
					//文件上传公共路由
					uploadFiles := my.Group("/upload")
					{
						uploadFiles.POST("/files", validatorFactory.Create(consts.ValidatorPrefix+"UploadFiles"))
					}

					myInfo := my.Group("/member")
					{
						//个人信息查询
						myInfo.GET("/info", (&home.Members{}).Info)
						//个人信息保存
						myInfo.POST("/save", validatorFactory.Create(consts.ValidatorPrefix+"UpdateInfo"))

						myInfo.POST("/change_phone", validatorFactory.Create(consts.ValidatorPrefix+"ChangePhone"))
					}

					myNeeds := my.Group("/needs")
					{
						//我的需求
						myNeeds.GET("/index", validatorFactory.Create(consts.ValidatorPrefix+"MyNeedIndex"))

						myNeeds.POST("/create", validatorFactory.Create(consts.ValidatorPrefix+"MyNeedCreate"))

						myNeeds.POST("/edit", validatorFactory.Create(consts.ValidatorPrefix+"MyNeedEdit"))

						myNeeds.GET("/info", (&home.Needs{}).MyInfo)

						myNeeds.GET("/delete", (&home.Needs{}).Delete)

						//Record
						myNeeds.GET("/record/index", validatorFactory.Create(consts.ValidatorPrefix+"MyNeedsRecordIndex"))
						myNeeds.POST("/record/close", validatorFactory.Create(consts.ValidatorPrefix+"MyNeedsRecordClose"))
						myNeeds.POST("/record/accept", validatorFactory.Create(consts.ValidatorPrefix+"MyNeedsRecordAccept"))

						myNeeds.GET("/record/list", validatorFactory.Create(consts.ValidatorPrefix+"MyNeedsRecordList"))
					}

					myServices := my.Group("/services")
					{
						//我的服务
						myServices.GET("/index", validatorFactory.Create(consts.ValidatorPrefix+"MyServicesIndex"))

						myServices.POST("/create", validatorFactory.Create(consts.ValidatorPrefix+"MyServicesCreate"))

						myServices.POST("/edit", validatorFactory.Create(consts.ValidatorPrefix+"MyServicesEdit"))

						myServices.GET("/info", (&home.Services{}).MyInfo)

						myServices.GET("/delete", (&home.Services{}).Delete)
						//Record
						myServices.GET("/record/index", validatorFactory.Create(consts.ValidatorPrefix+"MyServicesRecordIndex"))
						myServices.POST("/record/close", validatorFactory.Create(consts.ValidatorPrefix+"MyServicesRecordClose"))
						myServices.POST("/record/accept", validatorFactory.Create(consts.ValidatorPrefix+"MyServicesRecordAccept"))

						myServices.GET("/record/list", validatorFactory.Create(consts.ValidatorPrefix+"MyServicesRecordList"))
					}
				}
			}

		}

	}

	// Admin 相关的API
	backend := router.Group("/adminapi")
	{

		// 验证码
		passport := backend.Group("/passport")
		{
			// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
			passport.GET("/captcha/", (&captcha.Captcha{}).GenerateId)                          //  获取验证码ID
			passport.GET("/captcha/:captcha_id", (&captcha.Captcha{}).GetImg)                   // 获取图像地址
			passport.GET("/captcha/:captcha_id/:captcha_value", (&captcha.Captcha{}).CheckCode) // 校验验证码
			// 后台登录
			passport.Use(authorization.CheckCaptchaAuth()).POST("/login", validatorFactory.Create(consts.ValidatorPrefix+"AdminLogin"))

		}
		backend.Use(authorization.BindContextAuthToken(), authorization.CheckTokenAuth(), authorization.CheckAccessToAdmin())
		{
			// 退出登录
			passport.GET("/logout", (&admin.Users{}).LoginOut)
			//文件上传公共路由
			uploadFiles := backend.Group("/upload")
			{
				uploadFiles.POST("/files", validatorFactory.Create(consts.ValidatorPrefix+"UploadFiles"))
			}
			backend.GET("/settings", (&admin.Settings{}).GetValues)
			backend.POST("/settings", validatorFactory.Create(consts.ValidatorPrefix+"Settings"))
			// 方向分类组路由
			category := backend.Group("/category")
			{
				// 查询
				category.GET("/index", validatorFactory.Create(consts.ValidatorPrefix+"CategoryIndex"))
				// 新增
				category.POST("/create", validatorFactory.Create(consts.ValidatorPrefix+"CategoryCreate"))
				// 更新
				category.POST("/edit", validatorFactory.Create(consts.ValidatorPrefix+"CategoryEdit"))
				// 删除
				category.GET("/delete", (&admin.Category{}).Destroy)
				category.GET("/info", (&admin.Category{}).Info)
			}

			users := backend.Group("/users")
			{
				// 查询
				users.GET("/index", validatorFactory.Create(consts.ValidatorPrefix+"UsersIndex"))
				// 新增
				users.POST("/create", validatorFactory.Create(consts.ValidatorPrefix+"UsersCreate"))
				// 更新
				users.POST("/edit", validatorFactory.Create(consts.ValidatorPrefix+"UsersEdit"))
				// 删除
				users.GET("/delete", (&admin.Users{}).Destroy)
				// 详情
				users.GET("/info", (&admin.Users{}).Info)
				// 改密码
				users.POST("/change_passwd", validatorFactory.Create(consts.ValidatorPrefix+"UsersChangePasswd"))
			}

			members := backend.Group("/members")
			{
				members.GET("/index", validatorFactory.Create(consts.ValidatorPrefix+"MembersIndex"))
				members.GET("/info", (&admin.Members{}).Info)
				members.POST("/verify", validatorFactory.Create(consts.ValidatorPrefix+"MembersVerify"))
			}

			needs := backend.Group("/needs")
			{
				needs.GET("/index", validatorFactory.Create(consts.ValidatorPrefix+"AdminNeedIndex"))
				needs.GET("/info", (&admin.Needs{}).Info)
				needs.POST("/verify", validatorFactory.Create(consts.ValidatorPrefix+"AdminNeedVerify"))

				needs.GET("/record", validatorFactory.Create(consts.ValidatorPrefix+"AdminNeedRecord"))
			}

			services := backend.Group("/services")
			{
				services.GET("/index", validatorFactory.Create(consts.ValidatorPrefix+"AdminServicesIndex"))
				services.GET("/info", (&admin.Services{}).Info)
				services.POST("/verify", validatorFactory.Create(consts.ValidatorPrefix+"AdminServicesVerify"))

				services.GET("/record", validatorFactory.Create(consts.ValidatorPrefix+"AdminServicesRecord"))
			}
		}
		// 登录
	}
	return router
}
