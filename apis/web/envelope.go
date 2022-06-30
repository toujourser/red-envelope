package web

import (
	"resk/infra"
	"resk/infra/base"
	"resk/services"

	"github.com/kataras/iris/v12"
)

func init() {
	infra.RegisterApi(&EnvelopeApi{})
}

type EnvelopeApi struct {
	service services.RedEnvelopeService
}

func (e *EnvelopeApi) Init() {
	e.service = services.GetRedEnvelopeService()
	groupRouter := base.Iris().Party("/v1/envelope")
	groupRouter.Post("/sendout", e.sendOutHandler)
	groupRouter.Post("/receive", e.receiveHandler)

}

// {
// 	"envelopeNo": "",
//	"recvUsername": "",
//	"recvUserId": ""
//	}
func (e *EnvelopeApi) receiveHandler(ctx iris.Context) {
	dto := services.RedEnvelopeReceiveDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		return
	}
	item, err := e.service.Receive(dto)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		return
	}
	r.Data = item
	_, _ = ctx.JSON(r)
}

//{
//	"envelopeType": 0,
//	"username": "",
//	"userId": "",
//	"blessing": "",
//	"amount": "0",
//	"quantity": 0
//}
func (e *EnvelopeApi) sendOutHandler(ctx iris.Context) {
	dto := services.RedEnvelopeSendingDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		return
	}
	activity, err := e.service.SendOut(dto)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		return
	}
	r.Data = activity
	_, _ = ctx.JSON(r)

}
