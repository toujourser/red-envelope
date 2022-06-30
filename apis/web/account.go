package web

import (
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"resk/infra"
	"resk/infra/base"
	"resk/services"
)

// 定义web api 的时候。对每一个子业务，定义统一的前缀
// 资金账户的根路径定义为 /account
// 版本号 ：/v1/account

const (
	ResCodeBizTransferFailure = base.ResCode(6010)
)

func init() {
	infra.RegisterApi(new(AccountApi))
}

type AccountApi struct {
	service services.AccountService
}

func (a *AccountApi) Init() {
	a.service = services.GetAccountService()
	groupRouter := base.Iris().Party("/v1/account")
	groupRouter.Post("/create", a.createHandler)
}

// 账户创建的接口：/v1/account/create
func (a *AccountApi) createHandler(ctx iris.Context) {
	account := services.AccountCreatedDTO{}
	err := ctx.ReadJSON(&account)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	//执行创建账户的代码

	dto, err := a.service.CreateAccount(account)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		logrus.Error(err)
	}
	r.Data = dto
	ctx.JSON(&r)

}

// 转账的接口：/v1/account/transfer
func (a *AccountApi) transferHandler(ctx iris.Context) {
	account := services.AccountTransferDTO{}
	err := ctx.ReadJSON(&account)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	//执行转账逻辑
	status, err := a.service.Transfer(account)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		logrus.Error(err)
	}
	r.Data = status
	if status != services.TransferedStatusSuccess {
		r.Code = ResCodeBizTransferFailure
		r.Message = err.Error()
	}
	ctx.JSON(&r)
}

//查询红包账户的web接口: /v1/account/envelope/get
func (a *AccountApi) getEnvelopeAccountHandler(ctx iris.Context) {
	userId := ctx.URLParam("userId")
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if userId == "" {
		r.Code = base.ResCodeRequestParamsError
		r.Message = "用户ID不能为空"
		ctx.JSON(&r)
		return
	}
	account := a.service.GetEnvelopeAccountByUserId(userId)
	r.Data = account
	ctx.JSON(&r)
}

//查询账户信息的web接口：/v1/account/get
func (a *AccountApi) getAccountHandler(ctx iris.Context) {
	accountNo := ctx.URLParam("accountNo")
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if accountNo == "" {
		r.Code = base.ResCodeRequestParamsError
		r.Message = "账户编号不能为空"
		ctx.JSON(&r)
		return
	}
	account := a.service.GetAccount(accountNo)
	r.Data = account
	ctx.JSON(&r)
}
