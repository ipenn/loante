package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type Borrow struct {
	bun.BaseModel `bun:"table:borrow,alias:b"`
	Id                 	int			`json:"id"`
	Uid                	int			`json:"uid"`
	MchId              	int			`json:"mch_id"`
	ProductId          	int			`json:"product_id"`
	Postponed          	int			`json:"postponed"`
	LoanAmount         	int			`json:"loan_amount"`
	PostponedPeriod    	int			`json:"postponed_period"`
	BeRepaidAmount    	int			`json:"be_repaid_amount"`
	PostponeValuation 	int			`json:"postpone_valuation"`
	Status           	int			`json:"status"`
	LoanType         	int			`json:"loan_type"`
	Score            	int			`json:"score"`
	RiskModel        	int			`json:"risk_model"`
	ScoreTime        	string		`json:"score_time"`
	CreateTime       	string		`json:"create_time"`
	Payment          	int			`json:"payment"`
	LoanTime         	string		`json:"loan_time"`
	PaymentRequestNo 	string		`json:"payment_request_no"`
	PaymentRespond   	string		`json:"payment_respond"`
	EndTime          	string		`json:"end_time"`
	Remark           	string		`json:"remark"`
	CompleteTime     	string		`json:"complete_time"`
	ActualAmount     	int			`json:"actual_amount"`
}


func (a *Borrow)Insert()  {
	res, _ := global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
	id, _ := res.LastInsertId()
	a.Id = int(id)
}

func (a *Borrow)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *Borrow)Gets(where string) ([]Borrow, int) {
	var datas []Borrow
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Borrow)Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *Borrow)Update(where string)  {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Borrow)Page(where string, page, limit int) ([]User, int) {
	var datas []User
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("b.id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Borrow)Del(where string)  {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}