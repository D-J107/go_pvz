package app

import (
	"my_pvz/internal/app/transport/rest"
	"my_pvz/internal/app/transport/rest/middleware"
	"my_pvz/internal/db"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *db.DB) *gin.Engine {
	usersController := rest.NewUserController(db)
	pvzController := rest.NewPvzController(db)
	receptionsController := rest.NewReceptionController(db)
	productsController := rest.NewProductController(db)

	router := gin.Default()

	router.Use(middleware.PrometheusMetricsMiddleware())

	// <<__ users __>>
	router.POST("/dummyLogin", usersController.DummyLogin)
	router.POST("/register", usersController.Register)
	router.POST("/login", usersController.Login)

	// всё что ниже этой линии будет проходить через этот middleware
	router.Use(middleware.AuthMiddleware())

	// <<__ pvz __>>
	router.POST("/pvz", middleware.Authorize("moderator"), pvzController.Create)
	router.GET("/pvz", middleware.Authorize("moderator"), pvzController.GetAll)

	// <<__ Receptions __>>
	router.POST("/receptions", middleware.Authorize("employee"), receptionsController.Create)
	router.POST("/pvz/:pvzId/close_last_reception", middleware.Authorize("employee"), receptionsController.CloseLastReception)

	// <<__ Products __>>
	router.POST("/products", middleware.Authorize("employee"), productsController.Create)
	router.POST("/pvz/:pvzId/delete_last_product", middleware.Authorize("employee"), productsController.DeleteLastProduct)

	return router
}
