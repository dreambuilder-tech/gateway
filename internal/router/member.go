package router

import (
	"gateway/internal/common/rule"
	"gateway/internal/handler/member"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func memberRouter(r *gin.RouterGroup) {
	m := r.Group("/")
	{
		m.POST("/send-sms", withLimit(rule.SendSMS), member.SendSms)
		m.POST("/register", withLimit(rule.Register), member.Register)
		m.POST("/login", withLimit(rule.Login), member.Login)
		m.POST("/reset-pwd", member.ResetPwd)
		m.POST("/change-pwd", member.ChangePwd)
	}
	auth := r.Group("/", middleware.Auth())
	{
		auth.GET("/assets", member.Assets)
		ap := auth.Group("/profile")
		{
			ap.GET("/info", member.Profile)
			ap.POST("/change-avatar", member.ChangeAvatar)
			ap.POST("/change-nickname", member.ChangeNickname)
		}
		mp := auth.Group("/pin")
		{
			mp.POST("/commit", member.CommitPin)
			mp.POST("/reset", member.ResetPin)
			mp.POST("/change", withLimit(rule.ChangePIN), member.ChangePIN)
		}
		mPay := auth.Group("/payment") //支付
		{
			mPayMethod := mPay.Group("/method") //支付/方式
			{
				mPayMethod.GET("/list", member.PayMethodList)
				mPayMethod.POST("/add", member.AddPayMethod)
				mPayMethod.POST("/del", member.DelPayMethod)
				mPayMethod.POST("/set-def", member.SetDefPayMethod)
			}
		}
		mSecurity := auth.Group("/security") //支付
		{
			mSecurity.GET("/center", member.DeviceLoginRecordList)
		}

	}
	p := auth.Group("/phone")
	{
		p.POST("/verify", member.VerifyPhone)
		p.POST("/bind", member.PhoneBind)
	}
	realName := auth.Group("/realname")
	{
		realName.POST("/verify", member.VerifyRealName)
	}
}
