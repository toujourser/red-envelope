package envelops

import (
	"context"
	"resk/infra/base"
	"resk/services"

	"github.com/segmentio/ksuid"
	"github.com/tietang/dbx"
)

type itemDomain struct {
	RedEnvelopeItem
}

//生成itemNo
func (d *itemDomain) createItemNo() {
	d.ItemNo = ksuid.New().Next().String()
}

//创建Item

func (d *itemDomain) Create(item services.RedEnvelopeItemDTO) {
	d.RedEnvelopeItem.FromDTO(&item)
	d.RecvUsername.Valid = true
	d.createItemNo()
}

// Save 保存item数据
func (d *itemDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		id, err = dao.Insert(&d.RedEnvelopeItem)
		return err
	})
	return id, err
}

// GetOne 通过itemNo查询抢红包明细数据
func (d *itemDomain) GetOne(
	ctx context.Context, itemNo string) (dto *services.RedEnvelopeItemDTO) {
	err := base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		po := dao.GetOne(itemNo)
		if po != nil {
			dto = po.ToDTO()
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return dto
}

func (d *itemDomain) GetByUser(userId, envelopeNo string) (dto *services.RedEnvelopeItemDTO) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		po := dao.GetByUser(envelopeNo, userId)
		if po != nil {
			dto = po.ToDTO()
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return dto
}

// FindItems 通过envelopeNo查询已抢到红包列表
func (d *itemDomain) FindItems(envelopeNo string) (itemDtos []*services.RedEnvelopeItemDTO) {
	var items []*RedEnvelopeItem
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		items = dao.FindItems(envelopeNo)
		return nil
	})
	if err != nil {
		return itemDtos
	}

	itemDtos = make([]*services.RedEnvelopeItemDTO, 0)
	if len(items) == 0 {
		return itemDtos
	}
	var luckItem *services.RedEnvelopeItemDTO

	for i, po := range items {
		if po == nil {
			continue
		}
		item := po.ToDTO()

		if i == 0 {
			luckItem = item
		} else {
			if luckItem.Amount.Cmp(po.Amount) < 0 {
				luckItem = item
			}
		}
		itemDtos = append(itemDtos, item)
	}
	luckItem.IsLuckiest = true
	return itemDtos
}
