// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package recapSales

import (
	"fmt"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create recapSales process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	Session         *auth.SessionData
	PartnershipID   string                  `json:"partnership" valid:"required"`
	RecapSalesItems []RequestRecapSalesItem `json:"recap_sales_items" valid:"required"`
}

// RequestRecapSalesItem contain sales order id
type RequestRecapSalesItem struct {
	SalesOrderID string `json:"sales_order" valid:"required"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var partnerID int64
	// check partner recap
	if id, e := common.Decrypt(r.PartnershipID); e != nil {
		o.Failure("partnership", "not valid")
	} else {
		partnerID = id
		partner := &model.Partnership{ID: id, PartnershipType: "customer", IsDeleted: int8(0)}
		if e = partner.Read("ID", "PartnershipType", "IsDeleted"); e != nil {
			o.Failure("partnership", "not valid")
		}
	}

	checkDuplicate := make(map[string]bool)
	// check so
	for i, so := range r.RecapSalesItems {
		if soID, e := common.Decrypt(so.SalesOrderID); e != nil {
			o.Failure(fmt.Sprintf("recap_sales_items.%d.sales_order.invalid", i), "not valid")
		} else {
			sales := &model.SalesOrder{ID: soID, IsDeleted: int8(0)}
			if e = sales.Read("ID", "IsDeleted"); e != nil {
				o.Failure(fmt.Sprintf("recap_sales_items.%d.sales_order.invalid", i), "not valid")
			} else {
				sales.Customer.Read("ID")
				if sales.DocumentStatus == "requested_cancel" || sales.DocumentStatus == "approved_cancel" {
					o.Failure(fmt.Sprintf("recap_sales_items.%d.sales_order.invalid", i), " sales order cant be recap")
				}

				//check duplicate so id
				if checkDuplicate[so.SalesOrderID] == true {
					o.Failure(fmt.Sprintf("recap_sales_items.%d.sales_order.invalid", i), " sales order id duplicate")
				} else {
					checkDuplicate[so.SalesOrderID] = true
				}

				// partnership SO harus sama dengan partner request
				if sales.Customer.ID != partnerID {
					o.Failure(fmt.Sprintf("recap_sales_items.%d.sales_order.invalid", i), "multiple partner")
				}
			}
		}
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *createRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *createRequest) Transform() (model.RecapSales, error) {
	var totalAmount float64
	var code string
	var e error
	var recap model.RecapSales
	var recapItems []*model.RecapSalesItem
	partnerID, _ := common.Decrypt(r.PartnershipID)

	if code, e = util.CodeGen("code_customer_report", "recap_sales"); e == nil {
		// loop recap sales item
		for _, u := range r.RecapSalesItems {
			// read sales order
			soID, _ := common.Decrypt(u.SalesOrderID)
			so := &model.SalesOrder{ID: soID}
			so.Read("ID")
			// calculate total amount
			totalAmount = totalAmount + so.TotalCharge
			// append recap sales item
			recapItem := &model.RecapSalesItem{SalesOrder: so}
			recapItems = append(recapItems, recapItem)
		}

		recap = model.RecapSales{
			Partnership:     &model.Partnership{ID: partnerID},
			Code:            code,
			TotalAmount:     totalAmount,
			CreatedBy:       &model.User{ID: r.Session.User.ID},
			CreatedAt:       time.Now(),
			RecapSalesItems: recapItems,
		}
	}

	return recap, e
}

// cancelRequest data struct that stored request data when requesting an cancel recap sales process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type cancelRequest struct {
	RecapSales *model.RecapSales `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.RecapSales.IsDeleted == int8(1) {
		o.Failure("is_deleted", "document is deleted")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *cancelRequest) Messages() map[string]string {
	return map[string]string{}
}

func (r *cancelRequest) Transform() *model.RecapSales {
	r.RecapSales.IsDeleted = 1

	return r.RecapSales
}
