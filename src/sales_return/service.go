package salesReturn

import (
	"strings"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
)

// GetSalesReturn get all data Sales Return that matched with query request parameters.
// returning slices of Sales Return, total data without limit and error.
func GetSalesReturn(rq *orm.RequestQuery) (m *[]model.SalesReturn, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.SalesReturn))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.SalesReturn
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// ShowSalesReturn find a single data Sales Return using field and value condition.
func ShowSalesReturn(field string, values ...interface{}) (*model.SalesReturn, error) {
	m := new(model.SalesReturn)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	o.LoadRelated(m, "SalesReturnItems", 3)
	return m, nil
}

// UpdateSalesReturnItem will update data and other data that has foreign key to this data in database
// child data ID that did not exist will be create and child data ID that doesn't exist will be delete in database
// other data will be update which match with their ID in database
func UpdateSalesReturnItem(oldData []*model.SalesReturnItem, newData []*model.SalesReturnItem) (e error) {
	savedNewData := []int64{}
	oldDataID := []int64{}
	sliceOfQuestion := []string{}
	// save all new data and take new data id to savedNewData
	for _, newSRitem := range newData {
		if e = newSRitem.Save(); e != nil {
			return e
		}
		savedNewData = append(savedNewData, newSRitem.ID)
	}
	// filter old and new data and take an id of unsaved old data to oldDataID
	for _, oldSRitem := range oldData {
		updatedFlag := false
		for _, newID := range savedNewData {
			// this mean, oldSRItem id is saved as new
			if oldSRitem.ID == newID {
				updatedFlag = true
			}
		}
		// this mean, oldSRItem id is not save yet
		if updatedFlag == false {
			oldDataID = append(oldDataID, oldSRitem.ID)
			sliceOfQuestion = append(sliceOfQuestion, "?")
		}
	}

	//check if there is any old data left
	if len(oldDataID) > 0 {
		questionMarks := strings.Join(sliceOfQuestion, ",")
		o := orm.NewOrm()
		_, e = o.Raw("DELETE FROM sales_return_item WHERE id IN ("+questionMarks+")", oldDataID).Exec()
	}
	return e
}

//CancelSalesReturn to change document status to cancelled
func CancelSalesReturn(sr *model.SalesReturn) (err error) {
	sr.DocumentStatus = "cancelled"
	o := orm.NewOrm()
	if _, err = o.Update(sr, "document_status"); err != nil {
		return err
	}
	return nil
}

//UpdateExpense to update is_deleted finance expense
func UpdateExpense(sr *model.SalesReturn) (err error) {

	o := orm.NewOrm()

	_, err = o.Raw("update finance_expense set is_deleted = 1 where ref_type = 'sales_return' AND ref_id = ?; ", sr.ID).Exec()

	return err
}
