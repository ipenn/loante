package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

//ReferrerRiskConfig referrer_risk_config
type ReferrerRiskConfig struct {
	bun.BaseModel       `bun:"table:referrer_risk_config,alias:rrc"`
	Id                  int `json:"id" bun:",pk"`
	ReferrerId          int	`json:"referrer_id"`
	StatCompay          int	`json:"stat_compay"`
	RiskModel           int	`json:"risk_model"`
	NewMinScore         int	`json:"new_min_score"`
	NewMaxScore         int	`json:"new_max_score"`
	OldJumpRisk         int	`json:"old_jump_risk"`
	OldMinScore         int	`json:"old_min_score"`
	OldMaxScore         int	`json:"old_max_score"`
	PlatformOldMinScore int	`json:"platform_old_min_score"`
	PlatformOldMaxScore int	`json:"platform_old_max_score"`
	Remark              string	`json:"remark"`
	Referrer          ReferrerConfigLittle  `json:"referrer" bun:"rel:belongs-to,join:referrer_id=id"`

}

func (a *ReferrerRiskConfig) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ReferrerRiskConfig) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ReferrerRiskConfig) Set(col string, value interface{}, where string) error {
	_, err := global.C.DB.NewUpdate().Model(a).SetColumn(col, "?", value).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
		return err
	}
	return nil
}

func (a *ReferrerRiskConfig) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ReferrerRiskConfig) Page(where string, page, limit int) ([]ReferrerRiskConfig, int) {
	var datas []ReferrerRiskConfig
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Referrer").Where(where).Order(fmt.Sprintf("rrc.id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
