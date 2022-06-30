package accounts

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"resk/infra/base"
	_ "resk/testx"
	"testing"
)

func TestAccountDao_GetOne(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{runner: runner}
		Convey("通过编号查询账户数据", t, func() {
			a := &Account{
				AccountNo:    ksuid.New().Next().String(),
				AccountName:  "测试资金账户",
				AccountType:  0,
				CurrencyCode: "",
				UserId:       ksuid.New().Next().String(),
				Username:     sql.NullString{String: "测试用户", Valid: true},
				Balance:      decimal.NewFromFloat(100),
				Status:       1,
			}
			id, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)
			na := dao.GetOne(a.AccountNo)
			So(na, ShouldNotBeNil)
			So(na.Balance.String(), ShouldEqual, a.Balance.String())
			So(na.CreatedAt, ShouldNotBeNil)
			So(na.UpdatedAt, ShouldNotBeNil)
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}

func TestAccountDao_GetByUserId(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{runner: runner}
		Convey("通过用户ID和账户类型查询账户数据", t, func() {
			a := &Account{
				AccountNo:    ksuid.New().Next().String(),
				AccountName:  "测试资金账户",
				AccountType:  2,
				CurrencyCode: "",
				UserId:       ksuid.New().Next().String(),
				Username:     sql.NullString{String: "测试用户", Valid: true},
				Balance:      decimal.NewFromFloat(100),
				Status:       1,
			}
			id, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)
			na := dao.GetByUserId(a.UserId, a.AccountType)
			So(na, ShouldNotBeNil)
			So(na.Balance.String(), ShouldEqual, a.Balance.String())
			So(na.CreatedAt, ShouldNotBeNil)
			So(na.UpdatedAt, ShouldNotBeNil)
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}

func TestAccountDao_UpdateBalance(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{runner: runner}
		blance := decimal.NewFromFloat(100)
		Convey("更新账户余额", t, func() {
			a := &Account{
				AccountNo:    ksuid.New().Next().String(),
				AccountName:  "测试资金账户",
				AccountType:  2,
				CurrencyCode: "",
				UserId:       ksuid.New().Next().String(),
				Username:     sql.NullString{String: "测试用户", Valid: true},
				Balance:      blance,
				Status:       1,
			}
			id, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)

			// 增加余额
			Convey("增加余额", func() {
				amount := decimal.NewFromFloat(10)
				rows, err := dao.UpdateBalance(a.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 1)
				na := dao.GetOne(a.AccountNo)
				newBalance := blance.Add(amount)
				So(na, ShouldNotBeNil)
				So(na.Balance.String(), ShouldEqual, newBalance.String())
				So(na.CreatedAt, ShouldNotBeNil)
				So(na.UpdatedAt, ShouldNotBeNil)
			})
			// 扣减余额，余额充足
			Convey("扣减余额，余额充足", func() {
				a1 := dao.GetOne(a.AccountNo)
				So(a1, ShouldNotBeNil)
				t.Log(a1.Balance)
				amount := decimal.NewFromFloat(-30)
				rows, err := dao.UpdateBalance(a.AccountNo, amount)
				t.Log(a1.Balance)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 1)
				a2 := dao.GetOne(a.AccountNo)
				So(a2, ShouldNotBeNil)
				So("70", ShouldEqual, a2.Balance.String())
			})

			// 扣减余额，余额不足
			Convey("扣减余额，余额不够", func() {
				a1 := dao.GetOne(a.AccountNo)
				So(a1, ShouldNotBeNil)
				amount := decimal.NewFromFloat(-300)
				rows, err := dao.UpdateBalance(a.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 0)
				a2 := dao.GetOne(a.AccountNo)
				So(a2, ShouldNotBeNil)
				So(a1.Balance.String(), ShouldEqual, a2.Balance.String())
			})

		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}
