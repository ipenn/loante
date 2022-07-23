package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type ProductPrecept struct {
	bun.BaseModel `bun:"table:product_precept,alias:ppt" bun:",pk"`
	Id            int    `json:"id"`
	ProductId     int    `json:"product_id"`
	Status        int    `json:"status"`
	Amount        float64    `json:"amount"`
	MinCount      int    `json:"min_count"`
	CreateTime    string `json:"create_time"`
	Product 	  ProductLittle  `json:"product" bun:"rel:belongs-to,join:product_id=id"`
}

func (a *ProductPrecept) Insert() {
	a.CreateTime = tools.GetFormatTime()
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ProductPrecept) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ProductPrecept) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ProductPrecept) Page(where string, page, limit int) ([]ProductPrecept, int) {
	var datas []ProductPrecept
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Product").Where(where).Order(fmt.Sprintf("min_count asc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *ProductPrecept) Del(where string) {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}