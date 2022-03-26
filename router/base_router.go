package router

func (r router) BaseRouter() {
	handler := r.handler.BaseHandler()

	r.route.GET("/", handler.HandleReplaceImage)
}
