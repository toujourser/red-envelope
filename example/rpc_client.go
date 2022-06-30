package main

import (
	"net/rpc"
	"resk/services"

	"github.com/sirupsen/logrus"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func main() {
	c, err := rpc.Dial("tcp", ":8072")
	if err != nil {
		log.Error(err)
		return
	}

	//sendout(c)
	receive(c)

}

func receive(c *rpc.Client) {
	in := services.RedEnvelopeReceiveDTO{
		EnvelopeNo:   "2Apk8rFXJthW4B7IDpLGl3nN33f",
		RecvUserId:   "2Arpj2JdG0mP9s3cHC6u1idLqCu",
		RecvUsername: "2Arpj4QZcjvgiBnqBHBUKtJ0rAW",
		AccountNo:    "",
	}
	out := &services.RedEnvelopeItemDTO{}
	err := c.Call("EnvelopeRpc.Receive", in, out)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("%+v", out)
}

func sendout(c *rpc.Client) {
	in := services.RedEnvelopeSendingDTO{
		Amount:       decimal.NewFromFloat(1),
		UserId:       "1MD35g7HA9aukHZN5VEg2kTNYYx",
		Username:     "测试用户",
		EnvelopeType: services.GeneralEnvelopeType,
		Quantity:     2,
		Blessing:     "",
	}
	out := &services.RedEnvelopeActivity{}
	err := c.Call("EnvelopeRpc.SendOut", in, &out)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", out)
}
