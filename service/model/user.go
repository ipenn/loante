package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type UserLittle struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	Id            int    `json:"id"  bun:",pk"`
	Token         string `json:"token"`
	Phone         string `json:"phone"`
	AadhaarName   string `json:"aadhaar_name"`
}
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	UserLittle
	PanName                       string `json:"pan_name"`
	Gender                        int    `json:"gender"`
	New                           int    `json:"new"`
	Birthday                      string `json:"birthday"`
	CreateTime                    string `json:"create_time"`
	Blacklist                     int    `json:"blacklist"`
	EditAadName                   string `json:"edit_aad_name"`
	EditPanName                   string `json:"edit_pan_name"`
	Province                      string `json:"province"`
	Pin                           string `json:"pin"`
	IdAddress                     string `json:"id_address"`
	Area                          string `json:"area"`
	Email                         string `json:"email"`
	EmergencyContactName1         string `json:"emergency_contact_name_1"`
	EmergencyContactRelationship1 string `json:"emergency_contact_relationship_1"`
	EmergencyContactPhone1        string `json:"emergency_contact_phone_1"`
	EmergencyContactName2         string `json:"emergency_contact_name_2"`
	EmergencyContactRelationship2 string `json:"emergency_contact_relationship_2"`
	EmergencyContactPhone2        string `json:"emergency_contact_phone_2"`
	Education                     int    `json:"education"`
	Religion                      int    `json:"religion"`
	Remark                        string `json:"remark"`
	AadhaarPhotoFront             string `json:"aadhaar_photo_front"`
	AadhaarPhotoBack              string `json:"aadhaar_photo_back"`
	PanPhotoFront                 string `json:"pan_photo_front"`
	FacialPortrait                string `json:"facial_portrait"`
	KycVerifyBankcard             int    `json:"kyc_verify_bankcard"`
	KycVerifyBaseinfo             int    `json:"kyc_verify_baseinfo"`
	KycVerifyEmergency            int    `json:"kyc_verify_emergency"`
	KycVerifyFaceCompare          int    `json:"kyc_verify_face_compare"`
	KycVerifyFaceSearch           int    `json:"kyc_verify_face_search"`
	KycVerifyIdFront              int    `json:"kyc_verify_id_front"`
	KycVerifyIdBack               int    `json:"kyc_verify_id_back"`
	KycVerifyIdNumeric            int    `json:"kyc_verify_id_numeric"`
	KycVerifyFacial               int    `json:"kyc_verify_facial"`
	KycVerifyPanOcr               int    `json:"kyc_verify_pan_ocr"`
	KycVerifyPanNumeric           int    `json:"kyc_verify_pan_numeric"`
	KycVerifyUpdateTime           string `json:"kyc_verify_update_time"`
	Bankcard                      string `json:"bankcard"`
	Ifsc                          string `json:"ifsc"`
	BankName                      string `json:"bank_name"`
	BankDidTransferred            int    `json:"bank_did_transferred"`
	SocialAccounts 				  string `json:"social_accounts"`
	IdNo                          string `json:"id_no"`
}

type UserList struct {
	User
	ApplyingLoad				int		`json:"applying_load" bun:"applying_load"` //申请中的贷款
	Loading						int		`json:"loading" bun:"loading"` //贷款中的贷款
	LoanRepaid					int		`json:"loan_repaid" bun:"loan_repaid"` //已还款的贷款
}

func (a *User) Insert() {
	res, _ := global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
	id, _ := res.LastInsertId()
	a.Id = int(id)
}

func (a *User) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *User) Gets(where string) ([]User, int) {
	var datas []User
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *User) Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *User) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *User) Page(where string, page, limit int) ([]UserList, int) {
	var datas []UserList
	count, _ := global.C.DB.NewSelect().Model((*User)(nil)).
		ColumnExpr("u.*").
		ColumnExpr("b1.applying_load, b2.loading,b3.loan_repaid").
		Join("LEFT JOIN (select uid,count(*) as applying_load from borrow where status < 5 and status > 0 group by uid) as b1").JoinOn("b1.uid = u.id").
		Join("LEFT JOIN (select uid,count(*) as loading from borrow where status = 5 group by uid) as b2").JoinOn("b2.uid = u.id").
		Join("LEFT JOIN (select uid,count(*) as loan_repaid from borrow where status > 7 group by uid) as b3").JoinOn("b3.uid = u.id").
		Where(where).Order(fmt.Sprintf("u.id desc")).
		Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx, &datas)
	return datas, count
}

func (a *User) Del(where string) {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}

//其他表  user_privacy_app,user_privacy_call,user_privacy_contact,user_privacy_sms

type UserPrivacyApp struct {
	bun.BaseModel `bun:"table:user_privacy_app,alias:upa"`
	Id            int    `json:"id"`
	Uid           int    `json:"uid"`
	App           string `json:"app"`
	Version       string `json:"version"`
}
func (a *UserPrivacyApp) Page(where string, page, limit int) ([]UserPrivacyApp, int) {
	var datas []UserPrivacyApp
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("upa.id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
func (a *UserPrivacyApp) Gets(where string) ([]UserPrivacyApp, int) {
	var datas []UserPrivacyApp
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("upa.id desc")).ScanAndCount(global.C.Ctx)
	return datas, count
}

type UserPrivacyCall struct {
	bun.BaseModel `bun:"table:user_privacy_call,alias:upc"`
	Id            int    `json:"id"`
	Uid           int    `json:"uid"`
	Who           string `json:"who"`
	Phone         string `json:"phone"`
	Type          int    `json:"type"`
	Duration      string `json:"duration"`
	CreateTime    string `json:"create_time"`
}
func (a *UserPrivacyCall) Page(where string, page, limit int) ([]UserPrivacyCall, int) {
	var datas []UserPrivacyCall
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("upc.id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
func (a *UserPrivacyCall) Gets(where string) ([]UserPrivacyCall, int) {
	var datas []UserPrivacyCall
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("upc.id desc")).ScanAndCount(global.C.Ctx)
	return datas, count
}

type UserPrivacyContact struct {
	bun.BaseModel `bun:"table:user_privacy_contact,alias:upc"`
	Id            int    `json:"id"`
	Uid           int    `json:"uid"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	CreateTime    string `json:"create_time"`
	UpdateTime    string `json:"update_time"`
}
func (a *UserPrivacyContact) Page(where string, page, limit int) ([]UserPrivacyContact, int) {
	var datas []UserPrivacyContact
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("upc.id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
func (a *UserPrivacyContact) Gets(where string) ([]UserPrivacyContact, int) {
	var datas []UserPrivacyContact
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("upc.id desc")).ScanAndCount(global.C.Ctx)
	return datas, count
}

type UserPrivacySms struct {
	bun.BaseModel `bun:"table:user_privacy_sms,alias:ups"`
	Id            int    `json:"id"`
	Uid           int    `json:"uid"`
	Type          int    `json:"type"`
	Read          int    `json:"read"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Status        int    `json:"status"`
	Who           string `json:"who"`
	WhoPhone      string `json:"who_phone"`
	CreateTime    string `json:"create_time"`
}
func (a *UserPrivacySms) Page(where string, page, limit int) ([]UserPrivacySms, int) {
	var datas []UserPrivacySms
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("ups.id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
func (a *UserPrivacySms) Gets(where string) ([]UserPrivacySms, int) {
	var datas []UserPrivacySms
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("ups.id desc")).ScanAndCount(global.C.Ctx)
	return datas, count
}


type UserQuota struct {
	bun.BaseModel `bun:"table:user_quota,alias:uq"`
	Id            int    	`json:"id"`
	ProductId     int	`json:"product_id"`
	UserId        int	`json:"user_id"`
	Quota         int	`json:"quota"`
	CreateTime    string	`json:"create_time"`
	Remark        string	`json:"remark"`
}

func (a *UserQuota) Insert() {
	a.CreateTime = tools.GetFormatTime()
	res, _ := global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
	id, _ := res.LastInsertId()
	a.Id = int(id)
}

//Increase 产品额度提升
//19 repay.overdue.dont.increase.amount	1	逾期结清不提额（0：否，1：是）
//20 repay.overdue.reset.increase.amount	0	逾期结清时重置提额度（0：否，1：是）
//21 repay.fast.dont.increase.amount		1	还款间隔未满24小时不提额（0：否，1：是）
//pid 产品Id
//uid 用户Id
//overdue 是否逾期 9 逾期 8 正常
func (a *UserQuota) Increase(pid, uid, overdue int)  {
	//获取用户当前额度
	uquote := new(UserQuota)
	uquote.One(fmt.Sprintf("user_id = %d", uid))
	//获取规则
	config19 := new(SystemSetting)
	config19.One("id = 19") //逾期结清不提额（0：否，1：是）
	rules, _ := new(ProductPrecept).Page(fmt.Sprintf("ppt.status = 1 and ppt.product_id = %d", pid),1,1000)
	where := fmt.Sprintf("b.uid = %d and b.product_id = %d and (b.status = 8", uid, pid)
	remark := ""
	if config19.ParamValue == "0"{
		where += " or b.status = 9)"
		remark += "逾期结清也提额"
	}else{
		where += ")"
		remark += "逾期结清不提额"
	}
	config21 := new(SystemSetting)
	config21.One("id = 21")
	if config19.ParamValue == "1"{
		where += " and TIMESTAMPDIFF(HOUR,loan_time, complete_time) > 24"  //24小时不提额
		remark += " 24小时不提额"
	}else{
		remark += " 不满24小时也提额"
	}
	//统计个数
	borrow := new(Borrow)
	count := borrow.Count(where)
	amount := 0.0
	for _, item := range rules{
		if item.MinCount <= count{
			amount = item.Amount
		}
	}
	remark += fmt.Sprintf(" count = %d", count)
	quotaAmount := 0
	if amount > 0{
		// 插入提额数据
		uq := new(UserQuota)
		uq.ProductId = pid
		uq.UserId = uid
		uq.Quota = int(amount)
		uq.Remark = remark
		uq.Insert()
		quotaAmount = int(amount)
		global.Log.Info("插入提额数据 %v", uq)
	}
	//处理 逾期结清时重置提额度
	config20 := new(SystemSetting)
	config20.One("id = 20")
	if config20.ParamValue == "1" && overdue == 9{
		//获取产品初始额度
		product := new(Product)
		product.One(fmt.Sprintf("id = %d", pid))
		if product.Id > 0{
			uq := new(UserQuota)
			uq.ProductId = pid
			uq.UserId = uid
			uq.Quota = product.StartAmount
			uq.Remark = "逾期重置额度"
			uq.Insert()
			quotaAmount = product.StartAmount
			global.Log.Info("逾期重置额度 %v", uq)
		}
	}
	if uquote.Quota < quotaAmount { //提额发送短信通知
		//查询号码
		user := new(User)
		user.One(fmt.Sprintf("id = %d", uid))
		var ps []string
		if new(SmsTemplate).Send(6, user.Phone, ps) {

		}
	}
}
func (a *UserQuota) One(where string) {
	global.C.DB.NewSelect().Model(a).Where(where).Order("id desc").Scan(global.C.Ctx)
}