package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/resp"
	"loante/tools"
)

type statTraffic struct{}

func NewStatTraffic() *statTraffic {
	return new(statTraffic)
}

type statTrafficList struct {
	StartCreateTime string `query:"StartCreateTime" json:"StartCreateTime"`
	EndCreateTime   string `query:"EndCreateTime" json:"EndCreateTime"`
	MchId           string `query:"mch_id" json:"mch_id"`
	ProductId       string `query:"product_id" json:"product_id"`
	GroupTime       string `query:"group_time" json:"group_time"`
	GroupProduct    string `query:"group_product" json:"group_product"`
	GroupMch        string `query:"group_mch" json:"group_mch"`
	Page            int    `query:"page" json:"page"`
	Size            int    `query:"size" json:"size"`
}

func (a *statTraffic) StatTrafficList(c *fiber.Ctx) error {
	input := new(statTrafficList)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "st.id>0 "
	if input.MchId != "" {
		where += fmt.Sprintf(" and st.mch_id=%d", tools.ToInt(input.MchId))
	}
	if input.ProductId != "" {
		where += fmt.Sprintf(" and st.product_id=%d", tools.ToInt(input.ProductId))
	}
	if input.StartCreateTime != "" {
		where += fmt.Sprintf(" and st.create_time>'%s'", input.StartCreateTime)
		where += fmt.Sprintf(" and st.create_time<'%s'", input.EndCreateTime)
	}

	group := ""
	if input.GroupTime != "" {
		group += "DATE_FORMAT(time,'%Y-%m-%d') desc"
	}
	if input.GroupMch != "" {
		if group != "" {
			group += ","
		}
		group += "mch_id desc"
	}
	if input.GroupProduct != "" {
		if group != "" {
			group += ","
		}
		group += "product_id desc"
	}
	lists, count := new(model.StatTraffic).Page(where, group, input.Page, input.Size)
	//if len(lists) != 0 {
	//	for _, st := range lists {
	//		st.ApplyPassRate = tools.ToFloat64(st.ApplyPass) / tools.ToFloat64(st.Apply)
	//		st.LendingRate = tools.ToFloat64(st.LoanPass) / tools.ToFloat64(st.Apply)
	//		st.LoanPassRate = tools.ToFloat64(st.LoanPass) / tools.ToFloat64(st.ApplyPass)
	//	}
	//}
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}
