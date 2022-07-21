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
	StartCreateTime     string `query:"StartCreateTime" json:"StartCreateTime"`
	EndCreateTime       string `query:"EndCreateTime" json:"EndCreateTime"`
	MchId               string `query:"mch_id" json:"mch_id"`
	ProductId           string `query:"product_id" json:"product_id"`
	ReferrerConfigId    string `query:"referrer_config_id" json:"referrer_config_id"`
	GroupTime           string `query:"group_time" json:"group_time"`
	GroupProduct        string `query:"group_product" json:"group_product"`
	GroupMch            string `query:"group_mch" json:"group_mch"`
	GroupReferrerConfig string `query:"group_referrer_config" json:"group_referrer_config"`
	Page                int    `query:"page" json:"page"`
	Size                int    `query:"size" json:"size"`
}

// StatTrafficList 借贷日报
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
	if input.ReferrerConfigId != "" {
		where += fmt.Sprintf(" and st.referrer_config_id=%d", tools.ToInt(input.ReferrerConfigId))
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
	if input.GroupReferrerConfig != "" {
		if group != "" {
			group += ","
		}
		group += "referrer_config_id desc"
	}
	lists, count := new(model.StatTraffic).Page(where, group, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type afterReportUnifiedList struct {
	StartCreateTime string `query:"StartCreateTime" json:"StartCreateTime"`
	EndCreateTime   string `query:"EndCreateTime" json:"EndCreateTime"`
	Type            string `query:"type" json:"type"`
	ProductId       string `query:"product_id" json:"product_id"`
	GroupTime       string `query:"group_time" json:"group_time"`
	GroupProduct    string `query:"group_product" json:"group_product"`
	Group_type      string `query:"group_type" json:"group_type"`
	Page            int    `query:"page" json:"page"`
	Size            int    `query:"size" json:"size"`
}

// AfterReportUnifiedList 贷后报表
func (a *statTraffic) AfterReportUnifiedList(c *fiber.Ctx) error {
	input := new(afterReportUnifiedList)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "aru.id>0 "
	group := ""
	if input.Type != "" {
		where += fmt.Sprintf(" and aru.type=%d", tools.ToInt(input.Type))
	}
	if input.ProductId != "" {
		where += fmt.Sprintf(" and aru.product_id=%d", tools.ToInt(input.ProductId))
	}
	if input.StartCreateTime != "" {
		where += fmt.Sprintf(" and aru.create_time>'%s'", input.StartCreateTime)
		where += fmt.Sprintf(" and aru.create_time<'%s'", input.EndCreateTime)
	}

	if input.GroupTime != "" {
		group += "DATE_FORMAT(time,'%Y-%m-%d') desc"
	}
	if input.Group_type != "" {
		if group != "" {
			group += ","
		}
		group += "type desc"
	}
	if input.GroupProduct != "" {
		if group != "" {
			group += ","
		}
		group += "product_id desc"
	}
	lists, count := new(model.AfterReportUnified).Page(where, group, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type reportPayOut struct {
	StartCreateTime string `query:"StartCreateTime" json:"StartCreateTime"`
	EndCreateTime   string `query:"EndCreateTime" json:"EndCreateTime"`
	MchId           string `query:"mch_id" json:"mch_id"`
}

// ReportPayOut 付款监控报表
func (a *statTraffic) ReportPayOut(c *fiber.Ctx) error {
	input := new(reportPayOut)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := ""
	if input.MchId != "" {
		where += fmt.Sprintf(" and b.mch_id=%d", tools.ToInt(input.MchId))
	}
	if input.StartCreateTime != "" {
		where += fmt.Sprintf(" and b.create_time>'%s'", input.StartCreateTime)
		where += fmt.Sprintf(" and b.create_time<'%s'", input.EndCreateTime)
	}

	lists := new(model.Borrow).ForStatistics(where)
	return resp.OK(c, map[string]interface{}{
		"list": lists,
	})
}

type reportPayIn struct {
	StartCreateTime string `query:"StartCreateTime" json:"StartCreateTime"`
	EndCreateTime   string `query:"EndCreateTime" json:"EndCreateTime"`
	MchId           string `query:"mch_id" json:"mch_id"`
}

// ReportPayIn 收款监控报表
func (a *statTraffic) ReportPayIn(c *fiber.Ctx) error {
	input := new(reportPayIn)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := ""
	if input.MchId != "" {
		where += fmt.Sprintf(" and o.mch_id=%d", tools.ToInt(input.MchId))
	}
	if input.StartCreateTime != "" {
		where += fmt.Sprintf(" and o.create_time>'%s'", input.StartCreateTime)
		where += fmt.Sprintf(" and o.create_time<'%s'", input.EndCreateTime)
	}

	lists := new(model.Orders).ForStatistics(where)
	return resp.OK(c, map[string]interface{}{
		"list": lists,
	})
}
