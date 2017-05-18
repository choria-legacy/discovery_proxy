package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/choria-io/pdbproxy/models"
)

// DeleteSetSetOKCode is the HTTP code returned for type DeleteSetSetOK
const DeleteSetSetOKCode int = 200

/*DeleteSetSetOK Node Set

swagger:response deleteSetSetOK
*/
type DeleteSetSetOK struct {

	/*
	  In: Body
	*/
	Payload *models.Set `json:"body,omitempty"`
}

// NewDeleteSetSetOK creates DeleteSetSetOK with default headers values
func NewDeleteSetSetOK() *DeleteSetSetOK {
	return &DeleteSetSetOK{}
}

// WithPayload adds the payload to the delete set set o k response
func (o *DeleteSetSetOK) WithPayload(payload *models.Set) *DeleteSetSetOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete set set o k response
func (o *DeleteSetSetOK) SetPayload(payload *models.Set) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteSetSetOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteSetSetNotFoundCode is the HTTP code returned for type DeleteSetSetNotFound
const DeleteSetSetNotFoundCode int = 404

/*DeleteSetSetNotFound Not found

swagger:response deleteSetSetNotFound
*/
type DeleteSetSetNotFound struct {
}

// NewDeleteSetSetNotFound creates DeleteSetSetNotFound with default headers values
func NewDeleteSetSetNotFound() *DeleteSetSetNotFound {
	return &DeleteSetSetNotFound{}
}

// WriteResponse to the client
func (o *DeleteSetSetNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
}
