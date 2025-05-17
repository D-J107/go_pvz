package rest

import (
	prometheusMetrics "my_pvz/internal/app/transport/rest/prometheus_metrics"
	"my_pvz/internal/db"
	postgresql "my_pvz/internal/db/PostgreSQL"
	"my_pvz/internal/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productRepo   db.ProductRepository
	receptionRepo db.ReceptionRepository
}

func NewProductController(db *db.DB) *ProductController {
	return &ProductController{
		productRepo:   postgresql.NewPosgresProductRepositoryImpl(db),
		receptionRepo: postgresql.NewPostgesReceptionRepositoryImpl(db),
	}
}

func (pc *ProductController) Create(c *gin.Context) {
	logger.Log.Debug("received POST request to create new product")

	var req CreateProductRequest
	if err := c.BindJSON(&req); err != nil {
		logger.Log.Error("invalid product create request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentReception, err := pc.receptionRepo.GetLastByPvzID(c.Request.Context(), req.PvzId)
	if err != nil && err.Error() == "no rows in result set" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current PVZ does not contains any receptions!"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if currentReception.Status == "close" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current PVZ does not have open any open recepitons now!"})
		return
	}
	createdProduct, err := pc.productRepo.Create(c.Request.Context(), req.Type, currentReception.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := CreateProductResponse{
		Id:          createdProduct.ID,
		DateTime:    createdProduct.DateTime.Format(time.RFC3339),
		Type:        createdProduct.Type,
		ReceptionId: currentReception.ID,
	}
	prometheusMetrics.ProductsCreated.Inc()
	c.JSON(http.StatusCreated, resp)
}

func (pc *ProductController) DeleteLastProduct(c *gin.Context) {
	pvzId := c.Param("pvzId")
	// сначала пытаемся получить последнюю приёмку
	lastReception, err := pc.receptionRepo.GetLastByPvzID(c.Request.Context(), pvzId)
	if err != nil && err.Error() == "no rows in result set" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current pvz does not contain any receptions!"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if lastReception.Status == "close" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current pvz does not have in-progress reception!"})
		return
	}

	// пытаемся удалить товар
	err = pc.productRepo.DeleteLastProductInReception(c.Request.Context(), lastReception.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Товар удалён")
}
