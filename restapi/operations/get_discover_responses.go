package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/choria-io/pdbproxy/models"
)

// GetDiscoverOKCode is the HTTP code returned for type GetDiscoverOK
const GetDiscoverOKCode int = 200

/*GetDiscoverOK Basic successful discovery request

swagger:response getDiscoverOK
*/
type GetDiscoverOK struct {

	/*
	  In: Body
	*/
	Payload *models.DiscoverySuccessModel `json:"body,omitempty"`
}

// NewGetDiscoverOK creates GetDiscoverOK with default headers values
func NewGetDiscoverOK() *GetDiscoverOK {
	return &GetDiscoverOK{}
}

// WithPayload adds the payload to the get discover o k response
func (o *GetDiscoverOK) WithPayload(payload *models.DiscoverySuccessModel) *GetDiscoverOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get discover o k response
func (o *GetDiscoverOK) SetPayload(payload *models.DiscoverySuccessModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetDiscoverOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetDiscoverBadRequestCode is the HTTP code returned for type GetDiscoverBadRequest
const GetDiscoverBadRequestCode int = 400

/*GetDiscoverBadRequest Standard Error Format

swagger:response getDiscoverBadRequest
*/
type GetDiscoverBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewGetDiscoverBadRequest creates GetDiscoverBadRequest with default headers values
func NewGetDiscoverBadRequest() *GetDiscoverBadRequest {
	return &GetDiscoverBadRequest{}
}

// WithPayload adds the payload to the get discover bad request response
func (o *GetDiscoverBadRequest) WithPayload(payload *models.ErrorModel) *GetDiscoverBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get discover bad request response
func (o *GetDiscoverBadRequest) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetDiscoverBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
