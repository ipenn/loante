package global

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/wonderivan/logger"
)

type Config struct {
	DB *bun.DB
	Ctx context.Context
	Mysql struct{
		Host		string
		Port		int
		User	string
		Password	string
		DbName		string
	}
	Http		struct{
		Port		int
		Host		string
		Debug 		bool
	}
	//安全配置的
	Safety struct{
		Secret 		string
		Expired 	int
	}
	//Maps
	Maps struct{
		TestMap map[string]string 	`json:"test_map"`
		ServiceDeductType []string 	`json:"service_deduct_type"`
		ServiceType []string 	`json:"service_type"`
		SmsCompanys []string 	`json:"sms_companys"`
		SmsTypes []string 	`json:"sms_types"`
		BlankTypes []string 	`json:"blank_types"`
		SystemSettingTypes []string 	`json:"system_setting_types"`
		StatisticsCompany []string 	`json:"statistics_company"`
		RiskModel []string 	`json:"risk_model"`
		MchFundType []string 	`json:"mch_fund_type"`
	}  `json:"maps"`
}

var C Config
var Log *logger.LocalLogger
var AdminRights map[string]string
var Validate 		*validator.Validate

func init() {
	configInit()
	dbInit()
	logInit()
}

func configInit()  {
	if _, err := toml.DecodeFile("conf/config.toml", &C); err != nil {
		fmt.Println("toml.DecodeFile", err.Error())
	}
	AdminRights = map[string]string{}
	Validate = validator.New()
}

func dbInit()  {
	sqldb, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", C.Mysql.User, C.Mysql.Password, C.Mysql.Host, C.Mysql.Port, C.Mysql.DbName))
	if err != nil {
		fmt.Println(err)
		return
	}
	C.Ctx = context.Background()
	C.DB  = bun.NewDB(sqldb, mysqldialect.New())
	if C.Http.Debug {
		C.DB.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}

}

func logInit()  {
	Log = logger.NewLogger()
	Log.SetLogger("file", `{"filename":"./log/info.log","maxlines":10000,"maxsize":2,"append":true}`)
}
