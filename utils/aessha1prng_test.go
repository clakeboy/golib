package utils

import (
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	txt := `2kgrRntvBs0pm7rj{"AREA_NO":"440800","F_DATE":"2026-03-23 14:58:54","F_NO":"POSTAL20260323362523E3","data":{"BBR_LIST":[{"BBR_NAME":"洪中鹏","BBR_MOBILE":"13512321982","BBR_IDTYPE":"1","BBR_IDNO":"44**************79","BBR_EMAIL":"","TAG_FNUM":1}],"BBR_VIN_NO":"","INFOS":[{"KEY":"depart","NAME":"出单部门","VALUE":""}],"INS_END_DATE":"2026-08-31 23:59:59","INS_FOUND_DATE":"2025-09-01 00:00:00","INS_NO":"P2501026440899990000006979","IPS_AMOUNT":480000,"IPS_CODE":"PC-CIC-XPX-B","IPS_NAME":"学平险B款","IPS_YEAR":"1","ORDER_STAT":"2","ORD_DATE":"2025-09-12 23:52:45","ORD_NO":"P2501026440899990000006979","PAY_TOTAL":100,"PAY_TYPE":"2","SCORG":"52406400","SCSALER":"GD000042920","SCSER":"","S_NAME":"广东中华省分机构","S_NO":"gd-cic","TBR_EMAIL":"","TBR_IDNO":"44**************61","TBR_IDTYPE":"2","TBR_MOBILE":"13512321982","TBR_NAME":"陈妃霞"}}CC24840E31F2B9035688A1F1D1C3447F`
	key := "4549cc4890rt9039du52445629e089d45"

	res, err := SHAPRNGEncodeString(txt, key)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(res)

	res2, err := SHAPRNGDecodeString(res, key)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(res2))
}
