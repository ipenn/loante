package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type ProductDelayConfig struct {
	bun.BaseModel `bun:"table:product_delay_config,alias:pdc"`
	Id            int       `json:"id" bun:",pk"`
	MchId         int       `json:"mch_id"`
	ProductId     int       `json:"product_id"`
	DelayDay      int       `json:"delay_day"`
	DelayRate     float64   `json:"delay_rate"`
	Status        int       `json:"status"`
	IsShowDelay   int       `json:"is_show_delay"`
	CreateTime    string    `json:"create_time"`
	Merchant      *Merchant `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
	Product       *Product  `json:"product" bun:"rel:belongs-to,join:product_id=id"`
}

func (a *ProductDelayConfig) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ProductDelayConfig) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ProductDelayConfig) Page(where string, page, limit int) ([]ProductDelayConfig, int) {
	var d []ProductDelayConfig
	count, _ := global.C.DB.NewSelect().Model(&d).Relation("Merchant").Relation("Product").Where(where).Order(fmt.Sprintf("pdc.id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return d, count
}

func (a *ProductDelayConfig) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}