// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestHasElem(t *testing.T) {
	var IvStockIDS []int64
	ivstock1 := model.DummyItemVariantStock()
	IvStockIDS = append(IvStockIDS, ivstock1.ID)
	bool := HasElem(IvStockIDS, ivstock1.ID)
	assert.Equal(t, true, bool, "seharusnya true karena dalam ivStockIDS terdapat id ivstock1")

	ivstock2 := model.DummyItemVariantStock()
	bool = HasElem(IvStockIDS, ivstock2)
	assert.Equal(t, false, bool, "seharusnya false karena dalam ivStockIDS tidak terdapat id ivstock2")
}

func TestGetLastData(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from stockopname").Exec()

	stockopname := model.DummyStockopname()
	code, e := GetLastData("code", "stockopname")
	assert.NoError(t, e)
	assert.Equal(t, stockopname.Code, code)

	stockopname.Delete()
	code, e = GetLastData("code", "stockopname")
	assert.Error(t, e)
	assert.Empty(t, code)
}

func TestCodeGen(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from stockopname").Exec()

	code, e := CodeGen("code_stockopname", "stockopname")
	assert.NoError(t, e)
	assert.NotEmpty(t, code)
	assert.Equal(t, "WO#SO-00001", code)

	stockopname := model.DummyStockopname()
	stockopname.Code = "WO#SO-00001"
	stockopname.Save("Code")
	code, e = CodeGen("code_stockopname", "stockopname")
	assert.NoError(t, e)
	assert.NotEmpty(t, code)
	assert.Equal(t, "WO#SO-00002", code)
}

func TestGetApplicationSetting(t *testing.T) {
	as, e := GetApplicationSetting("application_setting_name", "code_stockopname")
	assert.NoError(t, e)
	assert.Equal(t, "code_stockopname", as.ApplicationSettingName)
}

func TestGenStringSKU(t *testing.T) {
	abcd := genStringSKU("kURENTO UZUMAKI dEKIL EMANG TERDEKIL ELAHHH")

	assert.Equal(t, "KUDET", abcd)
}

func TestGenerateCodeSKU(t *testing.T) {

	o := orm.NewOrm()
	o.Raw("delete from item_variant_stock").Exec()

	ivx := model.DummyItemVariant()
	ivx.VariantName = "KEREN UZUMAKI DEKIL EMANG TERDEKIL ELAHHH"
	ivx.Save("VariantName")
	ivx.Item.ItemName = "RIZAL"
	ivx.Item.Save("ItemName")

	code, _ := GenerateCodeSKU(ivx.ID)
	assert.Equal(t, "RKUDE-0001", code)

	//test generate SKU CODE where item variant id not nil
	iv := model.DummyItemVariant()
	iv.VariantName = "KEREN UZUMAKI DEKIL EMANG TERDEKIL ELAHHH"
	iv.Save("VariantName")
	iv.Item.ItemName = "RIZAL"
	iv.Item.Save("ItemName")

	ivs := model.DummyItemVariantStock()
	ivs.SkuCode = "RKUDE-0001"
	ivs.ItemVariant = iv
	ivs.Save("SkuCode", "ItemVariant")

	codex, _ := GenerateCodeSKU(iv.ID)
	assert.Equal(t, "RKUDE-0002", codex)
}

func TestGenerateCodeSKUUnique(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from item_variant").Exec()
	o.Raw("delete from item_variant_stock").Exec()

	fakeIVS := model.DummyItemVariantStock()
	fakeIVS.SkuCode = "KUDETS-0005"
	fakeIVS.Save()

	iv := model.DummyItemVariant()
	iv.Item.ItemName = "NARUTO"
	iv.VariantName = "KURENTO UZUMAKI DEKIL EMANG TERDEKIL ELAHHH"
	iv.Item.Save("ItemName")
	iv.Save("VariantName")

	ivs := model.DummyItemVariantStock()
	ivs.SkuCode, _ = GenerateCodeSKU(iv.ID)
	ivs.ItemVariant = iv
	ivs.Save("SkuCode", "ItemVariant")

	ivs1 := model.DummyItemVariantStock()
	ivs1.SkuCode, _ = GenerateCodeSKU(iv.ID)
	ivs1.ItemVariant = iv
	ivs1.Save("SkuCode", "ItemVariant")

	iv2 := model.DummyItemVariant()
	iv2.Item.ItemName = "NOBITA"
	iv2.VariantName = "KAKA UZUMAKI DEKIL EMANG TERDEKIL ELAHHH"
	iv2.Item.Save("ItemName")
	iv2.Save("VariantName")

	fakeIVS2 := model.DummyItemVariantStock()
	fakeIVS2.SkuCode = "KKUDET-0005"
	fakeIVS2.Save()

	ivs2 := model.DummyItemVariantStock()
	ivs2.SkuCode, _ = GenerateCodeSKU(iv2.ID)
	ivs2.ItemVariant = iv2
	ivs2.Save("SkuCode", "ItemVariant")

	assert.Equal(t, "NKUDE-0001", ivs.SkuCode)
	assert.Equal(t, "NKUDE-0002", ivs1.SkuCode)
	assert.Equal(t, "NKUDE-0003", ivs2.SkuCode)
}

func TestConvertDateToTimestamp(t *testing.T) {
	tt, _ := FormatDateToTimestamp(time.RFC3339, "2017-12-25")

	assert.Equal(t, "2017-12-25 00:00:00 +0000 UTC", tt.String())
}

func TestConvertDateToTimestampFalse(t *testing.T) {
	tt, _ := FormatDateToTimestamp(time.RFC3339, "2017-25-01")

	assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", tt.String())
}

func TestMaxDate(t *testing.T) {
	x := MaxDate()
	fmt.Println(x)
}
