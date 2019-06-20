// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package setting

import (
	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
)

// GetApplicationSetting get all data Application Setting that matched with query request parameters.
// returning slices of GetApplicationSetting, total data without limit and error.
func GetApplicationSetting(rq *orm.RequestQuery) (m *[]model.ApplicationSetting, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.ApplicationSetting))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.ApplicationSetting
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// ShowApplicationSetting find a single data Application Setting using field and value condition.
func ShowApplicationSetting(field string, values ...interface{}) (*model.ApplicationSetting, error) {
	m := new(model.ApplicationSetting)
	o := orm.NewOrm().QueryTable(m)
	err := o.Filter(field, values...).RelatedSel().Limit(1).One(m)

	return m, err
}

func UpdateApplicationSetting(r *updateRequest) (*model.ApplicationSetting, error) {
	var e error
	as := new(model.ApplicationSetting)
	as = r.ApplicationSetting
	as.Value = r.Value
	as.ApplicationSettingName = r.ApplicationSettingName
	if e = as.Save("Value", "ApplicationSettingName"); e == nil {
		// simpan account bank
		if len(r.BankAccounts) > 0 {
			for _, i := range r.BankAccounts {
				var ba *model.BankAccount
				if i.ID != 0 {
					ba = &model.BankAccount{
						ID:         i.ID,
						BankName:   i.BankName,
						BankNumber: i.BankNumber,
						IsDefault:  i.IsDefault,
					}
				} else {
					ba = &model.BankAccount{
						BankName:   i.BankName,
						BankNumber: i.BankNumber,
						IsDefault:  i.IsDefault,
					}
				}

				ba.Save()

			}
		}

		return as, nil

	}
	return nil, e
}
