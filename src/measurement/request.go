package measurement

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create setting process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	MeasurementName string            `json:"measurement_name" valid:"required"`
	Note            string            `json:"note"`
	SessionData     *auth.SessionData `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	_, err := getMeasurementByName(r.MeasurementName)

	if err == nil {
		o.Failure("measurement_name", "Measurement name already exists")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *createRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *createRequest) Transform() *model.Measurement {

	mx := &model.Measurement{
		MeasurementName: r.MeasurementName,
		Note:            r.Note,
	}

	return mx
}

// updateRequest data struct that stored request data when requesting an update setting process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateRequest struct {
	ID              int64
	MeasurementName string            `json:"measurement_name" valid:"required"`
	Note            string            `json:"note"`
	SessionData     *auth.SessionData `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	measurement, err := getMeasurementByName(r.MeasurementName)

	if err == nil && measurement.ID != r.ID {
		o.Failure("measurement_name", "Can't use this measurement name, already exists")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *updateRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *updateRequest) Transform(m *model.Measurement) *model.Measurement {

	mx := &model.Measurement{
		ID:              m.ID,
		MeasurementName: r.MeasurementName,
		Note:            r.Note,
	}

	return mx
}

// deleteRequest data struct that stored request data when requesting delete measurement process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type deleteRequest struct {
	Measurement *model.Measurement
}

// Validate implement validation.Requests interfaces.
func (r *deleteRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.Measurement.IsDeleted == int8(1) {
		o.Failure("is_deleted", "already deleted")
	}

	total := isMeasurementUsed(r.Measurement.ID)
	if total == true {
		o.Failure("measurement", "Can't delete this measurement because it's being used")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *deleteRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *deleteRequest) Transform() {
	r.Measurement.IsDeleted = int8(1)
}
