package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetSetSetHandlerFunc turns a function with the right signature into a get set set handler
type GetSetSetHandlerFunc func(GetSetSetParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetSetSetHandlerFunc) Handle(params GetSetSetParams) middleware.Responder {
	return fn(params)
}

// GetSetSetHandler interface for that can handle valid get set set params
type GetSetSetHandler interface {
	Handle(GetSetSetParams) middleware.Responder
}

// NewGetSetSet creates a new http.Handler for the get set set operation
func NewGetSetSet(ctx *middleware.Context, handler GetSetSetHandler) *GetSetSet {
	return &GetSetSet{Context: ctx, Handler: handler}
}

/*GetSetSet swagger:route GET /set/{set} getSetSet

Retrieves the query or nodes for a set

*/
type GetSetSet struct {
	Context *middleware.Context
	Handler GetSetSetHandler
}

func (o *GetSetSet) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewGetSetSetParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}