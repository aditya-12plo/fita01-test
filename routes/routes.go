package routes

import (
	"github.com/gin-gonic/gin"

	checkoutsController "fita-test-01/controllers/checkoutsController"
	homepageController "fita-test-01/controllers/homepageController"
	productsController "fita-test-01/controllers/productsController"
	promotionsController "fita-test-01/controllers/promotionsController"
)

type routes struct {
	router *gin.Engine
}

func (r routes) Run(addr ...string) error {
	return r.router.Run()
}

func InitializeRoute() routes {
	r := routes{
		router: gin.Default(),
	}

	r.router.Use(gin.Recovery())

	r.router.GET("/", homepageController.Index)

	product := r.router.Group("/products")
	r.routesProducts(product)

	promotion := r.router.Group("/promotions")
	r.routesPromotions(promotion)

	checkout := r.router.Group("/checkouts")
	r.routesCheckouts(checkout)

	return r
}

func (r routes) routesProducts(rg *gin.RouterGroup) {
	rg.GET("/index", productsController.Index)
}

func (r routes) routesPromotions(rg *gin.RouterGroup) {
	rg.GET("/index", promotionsController.Index)
}

func (r routes) routesCheckouts(rg *gin.RouterGroup) {
	rg.GET("/index", checkoutsController.Index)
	rg.POST("/data", checkoutsController.CheckoutData)
}
