package intf

type IInterceptor interface {
	Handle(responseData any)
}
