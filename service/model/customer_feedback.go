package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type CustomerFeedBack struct {
	bun.BaseModel `bun:"table:customer_feedback,alias:cf"`
	Id            int    `json:"id"`
	UserId        int    `json:"user_id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Whatsapp      string `json:"whatsapp"`
	Status        int    `json:"status"`
	Subject       string `json:"subject"`
	FeedbackType  int    `json:"feedback_type"`
	Content       string `json:"content"`
	PicUrl        string `json:"pic_url"`
	CreateTime    string `json:"create_time"`
	MarketName    string `json:"market_name"`
}

func (a *CustomerFeedBack) Page(where string, page, limit int) ([]CustomerFeedBack, int) {
	var d []CustomerFeedBack
	count, _ := global.C.DB.NewSelect().Model(&d).
		Where(where).Order("cf.id desc").Offset((page - 1) * limit).Limit(limit).
		ScanAndCount(global.C.Ctx)
	return d, count
}

func (a *CustomerFeedBack) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *CustomerFeedBack) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}
