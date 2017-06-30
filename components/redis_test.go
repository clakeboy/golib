package components

import (
	"testing"
	"log"
	"fmt"
)

func TestCKRedis_Set(t *testing.T) {
	cfg := &RedisConfig{
		RDDb:1,
		RDListName:"dudulist",
		RDServer:"168.168.0.10",
		RDPort:"6379",
		RDPassword:"123123",
		RDPoolSize:1000,
		RDIdleSize:50,
	}

	InitRedisPool(cfg)

	rd,err := NewCKRedis()
	if err != nil {
		panic(err)
	}

	val := `{"BeInsured":{"Address":"","CustomerType":"01","IdNo":"500106198403100830","IdType":"01","Mobile":"","Name":"李扬春"},"CarOwner":{"Address":"","CustomerType":"01","IdNo":"500106198403100830","IdType":"01","Mobile":"","Name":"李扬春"},"CommercePolicyBeginDate":"2016-09-02 00:00:00","CommercePolicyEndDate":"2017-09-01 23:59:59","CommercePolicyNo":"11821003980055046126","CommerceTotalPremium":"3320.81","CompulsoryPolicyBeginDate":"2016-08-26 14:00:00","CompulsoryPolicyEndDate":"2017-08-26 13:59:59","CompulsoryPolicyNo":"11821003900160358804","CompulsoryTotalPremium":"855","InsureCompany":"中国平安保险","JQInsurance":"J1","ProductCode":"PAZYCX","SYInsuranceItem":"Z1,136230.4|Z3,500000|Z4,10000|Z5,10000|B1,0|B3,0|F8,136230.4|B11,136230.4|B6,0|B4,0|B5,0|F12,0","ToInsured":{"Address":"","CustomerType":"01","IdNo":"500106198403100830","IdType":"01","Mobile":"","Name":"李扬春"},"TravelInsurance":"CCS","TravelTax":"360","Vehicle":{"BrandCode":"吉利美日MR7183C01轿车","CarTypeCode":"K33","CountryNature":"01","EnergyType":"1","EngineNo":"F7C400630","EnrollDate":"2015-08-30","Exhaust":"1.8","KerbWeight":"0","LicenseFlag":1,"LicenseNo":"渝B6S919","LicenseTypeCode":"02","LoadWeight":"0","LoanBank":"","LoanFlag":0,"NewFlag":0,"NonLocalFlag":0,"Price":"151700","PriceNoTax":"139800","Seat":"5","TransferFlag":0,"TransferFlagTime":"2017-06-15T16:00:00Z","UseNature":"2","Vin":"L6T7944Z8FN010762","Year":"2015"}}`

	err = rd.Set("3e3e371447a001f097f7f27927f525a1",val,300)
	if err != nil {
		log.Println(err)
	}
	defer rd.Close()
	data,err := rd.Get("3e3e371447a001f097f7f27927f525a1")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(data.([]byte)))
}
