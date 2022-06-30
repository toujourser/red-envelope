package accounts

import (
	"github.com/kataras/iris/v12/x/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"resk/infra/base"
	"resk/services"
	"sync"
)

var _ services.AccountService = new(accountService)

var once sync.Once

func init() {
	once.Do(func() {
		services.IAccountService = new(accountService)
	})
}

type accountService struct {
}

func (a *accountService) CreateAccount(dto services.AccountCreatedDTO) (*services.AccountDTO, error) {
	domain := accountDomain{}
	// 验证输入参数
	err := base.Validate().Struct(&dto)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range errs {
				logrus.Error(e.Translate(base.Transtate()))
			}
		}
		return nil, err
	}
	// 执行账户创建的业务逻辑
	amount, err := decimal.NewFromString(dto.Amount)
	if err != nil {
		return nil, err
	}
	account := services.AccountDTO{
		AccountName:  dto.AccountName,
		AccountType:  dto.AccountType,
		CurrencyCode: dto.CurrencyCode,
		UserId:       dto.UserId,
		Username:     dto.Username,
		Balance:      amount,
		Status:       1,
	}
	return domain.Create(account)
}

func (a *accountService) Transfer(dto services.AccountTransferDTO) (services.TransferedStatus, error) {
	// 验证输入参数
	err := base.Validate().Struct(&dto)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range errs {
				logrus.Error(e.Translate(base.Transtate()))
			}
		}
		return services.TransferedStatusFailure, err
	}
	// 执行转账逻辑
	amount, err := decimal.NewFromString(dto.AmountStr)
	if err != nil {
		return services.TransferedStatusFailure, err
	}
	dto.Amount = amount
	if dto.ChangeFlag == services.FlagTransferOut {
		if dto.ChangeType > 0 {
			return services.TransferedStatusFailure, errors.New("如果changeFlag为支出， 那么changeType必须小于零")
		}
	} else {
		if dto.ChangeType < 0 {
			return services.TransferedStatusFailure, errors.New("如果changeFlag为收入，那么changeTYpe必须大于零")
		}
	}
	domain := accountDomain{}

	return domain.Transfer(dto)
}

// 储值操作，交易主体和交易对象为同一人，直接将主体赋值为对象
func (a *accountService) StoreValue(dto services.AccountTransferDTO) (services.TransferedStatus, error) {
	dto.TradeTarget = dto.TradeBody
	dto.ChangeFlag = services.FlagTransferIn
	dto.ChangeType = services.AccountStoreValue
	return a.Transfer(dto)
}

func (a *accountService) GetEnvelopeAccountByUserId(userId string) *services.AccountDTO {
	domain := accountDomain{}
	return domain.GetEnvelopeAccountByUserId(userId)
}

func (a *accountService) GetAccount(accountNo string) *services.AccountDTO {
	domain := accountDomain{}
	return domain.GetAccount(accountNo)
}
