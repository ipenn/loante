package req

type PageReq struct {
	Page	int	`query:"page" json:"page"`
	Size	int	`query:"size" json:"size"`
}

type IdReq struct {
	Id	int	`query:"id" json:"id"`
}

type ModifyReq struct {
	Id	int	`json:"id" query:"id"`
	Key	string	`json:"key" query:"key"`
	Value	string	`json:"value" query:"value"`
}
