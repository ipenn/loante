package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type AdminMenu struct {
	bun.BaseModel `bun:"table:admin_menu,alias:am"`
	Id	int	`json:"id" bun:",pk"`
	Name	string	`json:"name"`
	Path	string	`json:"path"`
	Rights	string	`json:"rights"`
	Icon	string	`json:"icon"`
	ParentId	int	`json:"parent_id"`
}

func (a *AdminMenu)Gets(where string) []AdminMenu {
	var data []AdminMenu
	err := global.C.DB.NewSelect().Model(&data).Where(where).Order(fmt.Sprintf("id asc")).Scan(global.C.Ctx)
	if err != nil{
		fmt.Println(err.Error())
	}
	return data
}
func (a *AdminMenu)GetIds(ids []int) []AdminMenu {
	var data []AdminMenu
	err := global.C.DB.NewSelect().Model(&data).Where("id IN (?)", bun.In(ids)).Order(fmt.Sprintf("id asc")).Scan(global.C.Ctx)
	if err != nil{
		fmt.Println(err.Error())
	}
	return data
}