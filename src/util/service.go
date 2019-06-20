// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"git.qasico.com/mj/api/datastore/model"

	"time"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

// GenerateEanCode fungsi ini membuat random string dan valid ean 13
func GenerateEanCode() (code string) {
	random := common.RandomNumeric(9)
	code = fmt.Sprintf("211%09s", random)
	codeSplite := strings.Split(code, "")
	weightflag := true
	var sum int
	for i := len(code) - 1; i >= 0; i-- {
		totalweightflag := 1
		if weightflag == true {
			totalweightflag = 3
		}
		sum += common.ToInt(codeSplite[i]) * totalweightflag
		weightflag = !weightflag
	}

	total := (10 - (sum % 10)) % 10
	code = fmt.Sprintf("%s%d", code, total)

	return code
}

// CodeGen will generated a code base on code prefix and last code
// e.g appSettingName: "code_sales_order"
// NOTE untuk sekarang code hanya bisa sampai 99999, mis : WO#SO-00001 s/d WO#SO-99999
func CodeGen(appSettingName string, table string) (code string, e error) {
	// codeIndexs for number part in code
	var codeIndex, max int
	var min, ruleDigit string

	// get an initial data
	if code, max, min, ruleDigit, e = InitCode(appSettingName); e == nil {
		// check whether last data exist or not
		if lastCode, err := GetLastData("code", table); err == nil {
			// get a number part of code in lastCode
			number := regexp.MustCompile(`[\d]+$`).FindString(lastCode)
			codeIndex = common.ToInt(number)

			// check whether a codeIndex already maximum number or not
			if codeIndex == max {
				// then change a number to initial number e.g:000001
				number = min
				code = code + number
			} else {
				// if code index is not maximum then increment a code number
				codeIndex = codeIndex + 1
				number = fmt.Sprintf("%0"+ruleDigit, codeIndex)
				code = code + number
			}
		} else {
			// generated a new code if last data not exist
			code = fmt.Sprintf(code+"%s", min)
		}
	}
	//fmt.Println("code", code)
	return code, e
}

// InitCode to set initial data for codeGen
func InitCode(settingName string) (code string, maxDigit int, min string, ruleDigit string, e error) {
	var prefixValue map[string]string
	var appSett *model.ApplicationSetting

	// get a value code from application_setting_name in database
	if appSett, e = GetApplicationSetting("application_setting_name", settingName); e == nil {
		byteCode := []byte(appSett.Value)
		if e = json.Unmarshal(byteCode, &prefixValue); e == nil {

			// get prefix code
			rule := strings.Split(prefixValue["code_prefix"], "%")
			code = rule[0]

			// get last code digit e.g:[xxx-%5d] ---> 5d
			ruleDigit = rule[(len(rule) - 1)]

			// measure length number in code e.g:%5d ----> 5
			total := regexp.MustCompile(`[\d]`).FindString(rule[(len(rule) - 1)])
			totalDigit := common.ToInt(total)

			// get a minimum number for code digit e.g:00001
			min = fmt.Sprintf("%0"+rule[(len(rule) - 1)], 1)

			// get a number of 9 with length of code e.g:%5d ----> 99999 in int
			var max string
			for i := 1; i <= totalDigit; i++ {
				max = max + "9"
			}
			maxDigit = common.ToInt(max)

			return code, maxDigit, min, ruleDigit, nil
		}
	}
	return code, maxDigit, min, ruleDigit, e
}

// GetLastData will get lst data of the table in database
func GetLastData(field string, table string) (m string, e error) {
	o := orm.NewOrm()
	if e = o.Raw("SELECT " + field + " FROM " + table + " ORDER BY id DESC LIMIT 1").QueryRow(&m); e == nil {
		return m, nil
	}
	return
}

// GetApplicationSetting find a single data application_setting using field and value condition.
func GetApplicationSetting(field string, values ...interface{}) (*model.ApplicationSetting, error) {
	m := new(model.ApplicationSetting)
	o := orm.NewOrm().QueryTable(m)
	if err := o.Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// HasElem code for check elemen in array
// dalam stockopname fungsi ini diperlukan untuk mengecek item_variant_stock_id
// item_variant_stock_id tidak boleh sama dalam 1 x proses create stockopname
func HasElem(s interface{}, elem interface{}) bool {
	arrV := reflect.ValueOf(s)

	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {

			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}

	return false
}

// FormatDateToTimestamp formating date to timestamp.
func FormatDateToTimestamp(layout string, date string) (time.Time, error) {
	ds, _ := time.Parse("2006-01-02", date)
	dt := ds.Format(layout)

	return time.Parse(layout, dt)
}

// GenStringSKU generate variant for get prefix for SKU CODE
func genStringSKU(VarName string) string {
	var variantName string

	variantx := regexp.MustCompile(`\b[a-zA-Z]`).FindAllString(VarName, 5)
	variantName = strings.Join(variantx[:], "")
	prefix := strings.ToUpper(variantName)
	return prefix
}

// getLastCodeSKUprefix will get last sku code data table in database
func getLastCodeSKUprefix(prefix string) (m string, e error) {
	o := orm.NewOrm()
	if e = o.Raw("SELECT sku_code FROM item_variant_stock where sku_code like '" + prefix + "%' ORDER BY id DESC LIMIT 1").QueryRow(&m); e == nil {
		return m, nil
	}
	return
}

// GenerateCodeSKU for generate SKU Code
func GenerateCodeSKU(variantID int64) (code string, e error) {
	var m, varName string
	var firstInt int

	iv := &model.ItemVariant{ID: variantID}
	iv.Read("ID")
	iv.Item.Read("ID")

	varName = iv.Item.ItemName + " " + iv.VariantName
	prefix := genStringSKU(varName)

	if m, e = getLastCodeSKUprefix(prefix + "-"); e != nil {
		code = prefix + "-" + fmt.Sprintf("%04d", firstInt+1)
	} else {
		splitCode := strings.Split(m, "-")
		lastCode := common.ToInt(splitCode[1])
		number := fmt.Sprintf("%04d", lastCode+1)
		code = prefix + "-" + number
	}
	return
}

// MaxDate get maximum date
func MaxDate() (r string) {
	n := time.Now()
	var date string
	if n.Format("02") > "07" {
		date = n.Format("2006-01-") + "01"
	} else {
		date = n.AddDate(0, -1, 0).Format("2006-01-") + "01"
	}
	return date
}
