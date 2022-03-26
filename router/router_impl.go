package router

type Router interface {
	All()
	BaseRouter()
	NotFound()
}
