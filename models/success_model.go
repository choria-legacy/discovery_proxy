package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// SuccessModel success model
// swagger:model successModel
type SuccessModel struct {

	// detail
	Detail string `json:"detail,omitempty"`

	// HTTP Status Code
	Status int64 `json:"status,omitempty"`
}

// Validate validates this success model
func (m *SuccessModel) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *SuccessModel) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SuccessModel) UnmarshalBinary(b []byte) error {
	var res SuccessModel
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
