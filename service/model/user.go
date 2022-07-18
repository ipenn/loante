package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	Id	int		`json:"id"`
	Token	string		`json:"token"`
	Phone       string		`json:"phone"`
	AadhaarName string		`json:"aadhaar_name"`
	PanName     string		`json:"pan_name"`
	Gender      int		`json:"gender"`
	New			int		`json:"new"`
	Birthday   string		`json:"birthday"`
	CreateTime string		`json:"create_time"`
	Blacklist     int		`json:"blacklist"`
	EditAadName string		`json:"edit_aad_name"`
	EditPanName string		`json:"edit_pan_name"`
	Province    string		`json:"province"`
	Pin       string		`json:"pin"`
	IdAddress string		`json:"id_address"`
	Area      string		`json:"area"`
	Email                           string		`json:"email"`
	EmergencyContactName1         string		`json:"emergency_contact_name_1"`
	EmergencyContactRelationship1 string		`json:"emergency_contact_relationship_1"`
	EmergencyContactPhone1          string		`json:"emergency_contact_phone_1"`
	EmergencyContactName2         string		`json:"emergency_contact_name_2"`
	EmergencyContactRelationship2 string		`json:"emergency_contact_relationship_2"`
	EmergencyContactPhone2        string		`json:"emergency_contact_phone_2"`
	Education                       int		`json:"education"`
	Religion	int		`json:"religion"`
	Remark             string		`json:"remark"`
	AadhaarPhotoFront string		`json:"aadhaar_photo_front"`
	AadhaarPhotoBack string		`json:"aadhaar_photo_back"`
	PanPhotoFront   string		`json:"pan_photo_front"`
	PanPhotoBack        string		`json:"pan_photo_back"`
	FacialPortrait      string		`json:"facial_portrait"`
	KycVerifyBankcard    int		`json:"kyc_verify_bankcard"`
	KycVerifyBaseinfo       int		`json:"kyc_verify_baseinfo"`
	KycVerifyEmergency     int		`json:"kyc_verify_emergency"`
	KycVerifyFaceCompare int		`json:"kyc_verify_face_compare"`
	KycVerifyFaceSearch int		`json:"kyc_verify_face_search"`
	KycVerifyIdFront      int		`json:"kyc_verify_id_front"`
	KycVerifyIdBack    int		`json:"kyc_verify_id_back"`
	KycVerifyIdNumeric int		`json:"kyc_verify_id_numeric"`
	KycVerifyFacial        int		`json:"kyc_verify_facial"`
	KycVerifyPanOcr        int		`json:"kyc_verify_pan_ocr"`
	KycVerifyPanNumeric int		`json:"kyc_verify_pan_numeric"`
	KycVerifyUpdateTime string		`json:"kyc_verify_update_time"`
	Bankcard            string		`json:"bankcard"`
	Ifsc                 string		`json:"ifsc"`
	BankName           string		`json:"bank_name"`
	BankDidTransferred int		`json:"bank_did_transferred"`
	SocialAccounts     string		`json:"social_accounts"`
}

func (a *User)Insert()  {
	res, _ := global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
	id, _ := res.LastInsertId()
	a.Id = int(id)
}

func (a *User)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *User)Gets(where string) ([]User, int) {
	var datas []User
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *User)Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *User)Update(where string)  {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *User)Page(where string, page, limit int) ([]User, int) {
	var datas []User
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("u.id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *User)Del(where string)  {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}

//其他表  user_privacy_app,user_privacy_call,user_privacy_contact,user_privacy_sms

type UserPrivacyApp struct {
	bun.BaseModel `bun:"table:user_privacy_app,alias:upa"`
	Id	int	`json:"id"`
	Uid	int	`json:"uid"`
	App	string	`json:"app"`
	Version	string	`json:"version"`
}
func (a *UserPrivacyApp)Gets(where string) ([]UserPrivacyApp, int) {
	var datas []UserPrivacyApp
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("upa.id desc")).ScanAndCount(global.C.Ctx)
	return datas, count
}

type UserPrivacyCall struct {
	bun.BaseModel `bun:"table:user_privacy_call,alias:upc"`
	Id	int	`json:"id"`
	Uid	int	`json:"uid"`
	Who	string	`json:"who"`
	Phone	string	`json:"phone"`
	Type	int	`json:"type"`
	Duration   string	`json:"duration"`
	CreateTime string	`json:"create_time"`
}
func (a *UserPrivacyCall)Gets(where string) ([]UserPrivacyCall, int) {
	var datas []UserPrivacyCall
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("upc.id desc")).ScanAndCount(global.C.Ctx)
	return datas, count
}

type UserPrivacyContact struct {
	bun.BaseModel `bun:"table:user_privacy_contact,alias:upc"`
	Id	int	`json:"id"`
	Uid	int	`json:"uid"`
	Name	string	`json:"name"`
	Phone       string	`json:"phone"`
	CreateTime string	`json:"create_time"`
	UpdateTime string	`json:"update_time"`
}
func (a *UserPrivacyContact)Gets(where string) ([]UserPrivacyContact, int) {
	var datas []UserPrivacyContact
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("upc.id desc")).ScanAndCount(global.C.Ctx)
	return datas, count
}

type UserPrivacySms struct {
	bun.BaseModel `bun:"table:user_privacy_sms,alias:ups"`
	Id	int	`json:"id"`
	Uid	int	`json:"uid"`
	Type	int	`json:"type"`
	Read	int	`json:"read"`
	Title	string	`json:"title"`
	Content	string	`json:"content"`
	Status	int	`json:"status"`
	Who         string `json:"who"`
	WhoPhone   string `json:"who_phone"`
	CreateTime string `json:"create_time"`
}
func (a *UserPrivacySms)Gets(where string) ([]UserPrivacySms, int) {
	var datas []UserPrivacySms
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("ups.id desc")).ScanAndCount(global.C.Ctx)
	return datas, count
}

