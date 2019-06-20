package measurement

import (
	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
)

// GetMeasurement get all data Measurement that matched with query request parameters.
// returning slices of Measurement, total data without limit and error.
func GetMeasurement(rq *orm.RequestQuery) (m *[]model.Measurement, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.Measurement))

	// get total data
	if total, err = q.Filter("is_deleted", 0).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.Measurement
	if _, err = q.Filter("is_deleted", 0).All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// ShowMeasurement find a single data Measurement using field and value condition.
func ShowMeasurement(field string, values ...interface{}) (*model.Measurement, error) {
	m := new(model.Measurement)
	o := orm.NewOrm().QueryTable(m)
	if err := o.Filter(field, values...).Filter("is_deleted", 0).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// getMeasurementByName fungsi untuk mendapatkan measurement berdasarkan nama
func getMeasurementByName(name string) (*model.Measurement, error) {
	var measurement model.Measurement
	o := orm.NewOrm()
	err := o.Raw("SELECT * FROM measurement WHERE measurement_name = ? AND is_deleted = ?", name, 0).QueryRow(&measurement)

	return &measurement, err
}

// isMeasurementUsed untuk melakukan cek apabila measurement sedang digunakan pada salah satu item variant
func isMeasurementUsed(id int64) bool {
	var total int64
	o := orm.NewOrm()
	o.Raw("SELECT COUNT(*) AS TOTAL FROM item_variant WHERE item_variant.measurement_id = ?;", id).QueryRow(&total)

	if total > 0 {
		return true
	}

	return false
}
