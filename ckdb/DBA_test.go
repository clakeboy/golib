package ckdb

import (
	"fmt"
	"testing"
	"time"
	"reflect"
	"ck_go_lib/utils"
)

var cfg = &DBConfig{
	DBServer:   "168.168.0.10",
	DBPort:     "3306",
	DBName:     "pcbx_ddb",
	DBUser:     "root",
	DBPassword: "kKie93jgUrn!k",
	DBPoolSize: 200,
	DBIdleSize: 100,
	DBDebug:    true,
}

type CarInfo struct {
	ViId                   int    `json:"vi_id" bson:"vi_id"`                                         //车辆信息表ID
	VlaId                  int    `json:"vla_id" bson:"vla_id"`                                       //车型库表id
	UsrId                  int    `json:"usr_id" bson:"usr_id"`                                       //用户ID
	CityId                 int    `json:"city_id" bson:"city_id"`                                     //所在城市ID
	CityName               string `json:"city_name" bson:"city_name"`                                 //所在城市名称
	WhetherCard            int    `json:"whether_card" bson:"whether_card"`                           //是否上牌(0未上牌，1已上牌)
	VehicleId              string `json:"vehicle_id" bson:"vehicle_id"`                               //车辆识别码(车架号)
	EngineNo               string `json:"engine_no" bson:"engine_no"`                                 //发动机号
	EngineModel            string `json:"engine_model" bson:"engine_model"`                           //发动机型号
	CarName                string `json:"car_name" bson:"car_name"`                                   //品牌型号
	PlateNumber            string `json:"plate_number" bson:"plate_number"`                           //车牌号
	FirstRegisterDate      string `json:"first_register_date" bson:"first_register_date"`             //初登时间
	DeviceStatus           int    `json:"device_status" bson:"device_status"`                         //0未安装，1已安装
	OwnerInfo              string `json:"owner_info" bson:"owner_info"`                               //车主信息
	InsuredInfo            string `json:"insured_info" bson:"insured_info"`                           //被保人信息 序列化存储
	CreatedAt              int    `json:"created_at" bson:"created_at"`                               //创建时间
	CarSeat                int    `json:"car_seat" bson:"car_seat"`                                   //核定座位数
	VehicleOrigin          string `json:"vehicle_origin" bson:"vehicle_origin"`                       //
	JqQueryNo              string `json:"jq_query_no" bson:"jq_query_no"`                             //交强险投保查询码
	SyQueryNo              string `json:"sy_query_no" bson:"sy_query_no"`                             //商业险投保查询码
	SendQueryNo            string `json:"send_query_no" bson:"send_query_no"`                         //发送的查询码
	CInsuredEffectiveDate  int    `json:"c_insured_effective_date" bson:"c_insured_effective_date"`   //生效日期(交强)
	CInsuredExpirationDate int    `json:"c_insured_expiration_date" bson:"c_insured_expiration_date"` //终止日期(交强)
	BInsuredEffectiveDate  int    `json:"b_insured_effective_date" bson:"b_insured_effective_date"`   //生效日期(商业)
	BInsuredExpirationDate int    `json:"b_insured_expiration_date" bson:"b_insured_expiration_date"` //终止日期(商业)
	JqInsuredPdf           string `json:"jq_insured_pdf" bson:"jq_insured_pdf"`                       //交强险pdf保单
	SyInsuredPdf           string `json:"sy_insured_pdf" bson:"sy_insured_pdf"`                       //商业险pdf保单
	OrdFee                 int    `json:"ord_fee" bson:"ord_fee"`                                     //订单金额（分）
	PayDaily               int    `json:"pay_daily" bson:"pay_daily"`                                 //每日支付金额（分）
	DrivingLicenseImg      string `json:"driving_license_img" bson:"driving_license_img"`             //行驶证照片(正本)
	DrivingLicenseCopy     string `json:"driving_license_copy" bson:"driving_license_copy"`           //行驶证照片(副本)
	DrivingLicenseBack     string `json:"driving_license_back" bson:"driving_license_back"`           //行驶证副本反面
	CarCertificate         string `json:"car_certificate" bson:"car_certificate"`                     //车辆合格证照片(非进口新车必传)
	CarInvoice             string `json:"car_invoice" bson:"car_invoice"`                             //购车发票照片(进口新车必传)
	CustomsOrder           string `json:"customs_order" bson:"customs_order"`                         //关单照片(进口新车必传)
	SyInsuranceItem        string `json:"sy_insurance_item" bson:"sy_insurance_item"`                 //商业险返回数据
	DeductionDays          int    `json:"deduction_days" bson:"deduction_days"`                       //总扣除天数
	DeductibleDays         int    `json:"deductible_days" bson:"deductible_days"`                     //抵扣天数，不减应扣天数
	ExchangeDays           int    `json:"exchange_days" bson:"exchange_days"`                         //当前车辆已兑换天数
	BankCard               string `json:"bank_card" bson:"bank_card"`                                 //默认扣款银行卡号
	UsrName                string `json:"usr_name" bson:"usr_name"`                                   //户主姓名
	UsrLicense             string `json:"usr_license" bson:"usr_license"`                             //证件号码
	IsDel                  int    `json:"is_del" bson:"is_del"`                                       //是否删除 1删除
	IdType                 string `json:"id_type" bson:"id_type"`                                     //证件类型
	IdTypeName             string `json:"id_type_name" bson:"id_type_name"`                           //证件类型名称
}

func TestDBA_Insert(t *testing.T) {
	db,err := NewDBA(cfg)
	if err != nil {
		panic(err)
	}
	tab := db.Table("t_vehicle_info")
	list ,err := tab.Limit(10,1).Query().Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v",list[0])
}

func TestDBA_WhereRecursion(t *testing.T) {

	dba, err := NewDBA(cfg)
	if err != nil {
		panic(err)
	}

	var params struct{
		TagName string `json:"tag_name"`
		TagNum   int `json:"tag_num"`
		TagCreatedDate int `json:"tag_created_date"`
	}

	params.TagName = "我去"
	params.TagNum = 3
	params.TagCreatedDate = int(time.Now().Unix())

	data, err := dba.ConvertData(&params)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

//扫描数据到传入的类型
func scanType(scans []interface{},columns []string,i interface{}) interface{} {
	if i == nil {
		return scanMap(scans,columns)
	}
	t := reflect.TypeOf(i)
	switch t.Kind() {
	case reflect.Ptr:
		return scanStruct(t.Elem(),scans,columns)
	case reflect.Struct:
		return scanStruct(t,scans,columns)
	case reflect.Map:
		fallthrough
	default:
		return scanMap(scans,columns)
	}
}
//扫描数据到结构体
func scanStruct(t reflect.Type,scans []interface{},columns []string) interface{} {
	obj := reflect.New(t).Interface()
	objV := reflect.ValueOf(obj)
	for i,colName := range columns {
		idx := findTagOfStruct(t,colName)
		if idx != -1 {
			scans[i] = objV.Field(idx).Interface()
		}
	}
	return obj
}
//在结构体查找TAG值是否存在
func findTagOfStruct(t reflect.Type,colName string) int {
	for i:=0;i<t.NumField();i++ {
		val,ok := t.Field(i).Tag.Lookup(colName)
		if ok && val == colName {
			return i
		}
	}
	return -1
}
//扫描数据到MAP 默认 utils.M
func scanMap(scans []interface{},columns []string) interface{} {
	obj := utils.M{}

	return obj
}