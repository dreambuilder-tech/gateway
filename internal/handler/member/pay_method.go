package member

import (
	"gateway/internal/common/auth"
	"slices"
	"wallet/common-lib/app"
	"wallet/common-lib/consts/pay_method"
	"wallet/common-lib/dto/req_dto"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/member_rpcx"
	"wallet/common-lib/utils/stringx"

	"github.com/gin-gonic/gin"
)

type PayMethodListReq struct {
	req_dto.PageArgs
}
type AddPayMethodReq struct {
	Code       pay_method.Code `json:"code"`
	MemberName string          `json:"member_name"`
	Account    string          `json:"account"`
	FileId     int64           `json:"file_id"`
	Bank       string          `json:"bank"`
	Branch     string          `json:"branch"`
}

type DelPayMethodReq struct {
	Id int64 `json:"id"` // id(pk) of member_pay_methods
}

type SetDefPayMethodReq struct {
	Id int64 `json:"id"` // id(pk) of member_pay_methods
}

func PayMethodList(c *gin.Context) {
	req := PayMethodListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}

	resp, err := member_rpcx.PayMethodList(c.Request.Context(), &contracts.PayMethodListReq{
		PageArgs: req.PageArgs,
		MemberID: auth.MemberID(c),
	})
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.ResultPage(c, resp.List, resp.Total)
}

func checkForPayMethod(req *AddPayMethodReq) bool {
	/*
		（1）BANKCARD：姓名（任意中文）、卡号（至少16位任意数宇）、银行（任意中文）、支行（任意中文）
		卡号支持三方扫描识别
		（2） WEPAY：姓名（任意中文）、微信号（字母or数字，不能中文）、收款码收款码支持系统扫描识别是否为微信收款码（若不是，toast报错提示”收款码不正确”）
		（3） ALIPAY：姓名（任意中文）、账号（字母or数字，不能中文）、收款码收款码支持系统扫描识别是否为支付宝收款码（若不是，toast报错提示”收款码不正确”〉
		（4） ECNY：姓名（任意中文）、钣包编号（16位数字
		`code`               VARCHAR(30)  NOT NULL COMMENT '支付方式代码(BANKCARD,WEPAY,ALIPAY,ECNY)',
		`member_name`        VARBINARY(128)  NOT NULL COMMENT '姓名(加密)；普通用户(role=0)是真实姓名，码商是姓名',
		`member_name_digest` CHAR(43)     NOT NULL DEFAULT '' COMMENT '姓名(索引)',
		`account`            VARBINARY(128)  NOT NULL COMMENT '支付账号(加密)[卡号/微信号/支付宝账号/钱包编号]',
		`account_digest`     CHAR(43)     NOT NULL DEFAULT '' COMMENT '支付账号(索引)',
		`bank`               VARCHAR(60)  NOT NULL DEFAULT '' COMMENT '开户银行',
		`branch`             VARCHAR(150) NOT NULL DEFAULT '' COMMENT '开户支行',
	*/
	switch req.Code {

	case pay_method.BankCard:
		// （1）银行卡
		// 姓名（任意中文）
		// 卡号（至少16位数字，支持扫描）
		// 银行（任意中文）
		// 支行（任意中文）
		//如果传入会员名称，就做检查（普通用户可以不传）
		if len(req.MemberName) > 0 && !stringx.IsChinese(req.MemberName, 20) {
			return false
		}
		if !stringx.IsChinese(req.Bank, 60) || !stringx.IsChinese(req.Branch, 150) {
			return false
		}
		if len(req.Account) < 16 || !stringx.IsNumber(req.Account) {
			return false
		}
	case pay_method.Wepay, pay_method.Alipay:
		// （2）微信
		// 姓名（任意中文）
		// 微信号（字母或数字，不能中文）
		// 收款码必须是微信收款码
		if len(req.MemberName) > 0 && !stringx.IsChinese(req.MemberName, 20) {
			return false
		}
		if !stringx.IsNumbersOrLetters(req.Account) {
			return false
		}
	case pay_method.Ecny:
		// （4）数字人民币
		// 姓名（任意中文）
		// 钱包编号（16位数字）
		if len(req.MemberName) > 0 && !stringx.IsChinese(req.MemberName, 20) {
			return false
		}
		if len(req.Account) != 16 || !stringx.IsNumber(req.Account) {
			return false
		}
	}
	return true
}

func AddPayMethod(c *gin.Context) {
	// 1. check binding
	req := new(AddPayMethodReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	// 2. check empty and invalid
	if req.Code == "" || req.Account == "" {
		app.InvalidParams(c, "empty parameters")
		return
	}

	if slices.Contains(pay_method.AllCodeHasQrCode, req.Code) && req.FileId <= 0 {
		app.InvalidParams(c, "empty parameters")
		return
	}

	if req.Code == pay_method.BankCard {
		if req.Bank == "" || req.Branch == "" {
			app.InvalidParams(c, "empty parameters")
			return
		}
	}

	// 3. check logic
	if !slices.Contains(pay_method.AllCode, req.Code) || !checkForPayMethod(req) {
		app.InvalidParams(c, "invalid parameters")
		return
	}

	// 4. invoke rpcx
	addPayMethodReq := &contracts.AddPayMethodReq{
		MemberID:   auth.MemberID(c),
		Code:       req.Code,
		MemberName: req.MemberName,
		Account:    req.Account,
		Bank:       req.Bank,
		Branch:     req.Branch,
		FileId:     req.FileId,
	}
	resp, err := member_rpcx.AddPayMethod(c.Request.Context(), addPayMethodReq)

	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	// 3. return result
	app.Result(c, resp)
}

func DelPayMethod(c *gin.Context) {
	// 1.check binding
	req := new(DelPayMethodReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	// 2.invoke rpcx
	resp, err := member_rpcx.DelPayMethod(c.Request.Context(), auth.MemberID(c), req.Id)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	// 3. return result
	app.Result(c, resp)
}

func SetDefPayMethod(c *gin.Context) {
	// 1.check binding
	req := new(SetDefPayMethodReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	// 2.invoke rpcx
	resp, err := member_rpcx.SetDefPayMethod(c.Request.Context(), &contracts.SetDefPayMethodReq{
		Id:       req.Id,
		MemberId: auth.MemberID(c),
	})

	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	// 3. return result
	app.Result(c, resp)
}
