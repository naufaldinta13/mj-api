package setting

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/validation"
)

type bankAccount struct {
	ID         int64  `json:"id"`
	IsDefault  int8   `json:"is_default"`
	BankName   string `json:"bank_name"`
	BankNumber string `json:"bank_number"`
}

// updateRequest data struct that stored request data when requesting an update setting process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateRequest struct {
	ApplicationSettingName string        `json:"application_setting_name" valid:"required"`
	Value                  string        `json:"value" valid:"required"`
	BankAccounts           []bankAccount `json:"bank_accounts"`

	SessionData        *auth.SessionData         `json:"-"`
	ApplicationSetting *model.ApplicationSetting `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	// yang boleh input bank account hanya usergroup 2 yaitu owner
	if len(r.BankAccounts) == 0 && r.SessionData.User.Usergroup.ID == 2 {
		o.Failure("bank_accounts.valid", "Bank Accounts is required")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *updateRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *updateRequest) Transform() {
	return
}
