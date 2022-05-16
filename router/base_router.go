package router

func (r router) BaseRouter() {
	handler := r.handler.BaseHandler()

	r.route.GET("/", handler.HandleReplaceImage)
	r.route.GET("/url", handler.HandleReplaceImageUrl)
	r.route.GET("/insert", handler.HandleBulkInsertNumber)
}
