package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PutSetSetHandlerFunc turns a function with the right signature into a put set set handler
type PutSetSetHandlerFunc func(PutSetSetParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PutSetSetHandlerFunc) Handle(params PutSetSetParams) middleware.Responder {
	return fn(params)
}

// PutSetSetHandler interface for that can handle valid put set set params
type PutSetSetHandler interface {
	Handle(PutSetSetParams) middleware.Responder
}

// NewPutSetSet creates a new http.Handler for the put set set operation
func NewPutSetSet(ctx *middleware.Context, handler PutSetSetHandler) *PutSetSet {
	return &PutSetSet{Context: ctx, Handler: handler}
}

/*PutSetSet swagger:route PUT /set/{set} putSetSet

Update a set

*/
type PutSetSet struct {
	Context *middleware.Context
	Handler PutSetSetHandler
}

func (o *PutSetSet) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewPutSetSetParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}