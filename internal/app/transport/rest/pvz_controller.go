package rest

import (
	"context"
	prometheusMetrics "my_pvz/internal/app/transport/rest/prometheus_metrics"
	"my_pvz/internal/db"
	postgresql "my_pvz/internal/db/PostgreSQL"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PvzController struct {
	pvzRepo       db.PvzRepository
	receptionRepo db.ReceptionRepository
	productRepo   db.ProductRepository
}

func NewPvzController(db *db.DB) *PvzController {
	return &PvzController{
		pvzRepo:       postgresql.NewPostgresPvzRepositoryImpl(db),
		receptionRepo: postgresql.NewPostgesReceptionRepositoryImpl(db),
		productRepo:   postgresql.NewPosgresProductRepositoryImpl(db),
	}
}

func (pc *PvzController) Create(c *gin.Context) {
	var createRequest PvzCreationRequest
	if err := c.BindJSON(&createRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !isValidUUID(createRequest.Id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pvz id has wrong format! must be UUID."})
		return
	}
	created, err := pc.pvzRepo.Create(context.Background(), createRequest.Id, createRequest.RegistrationDate, createRequest.City)
	if err != nil && strings.Contains(err.Error(), `duplicate key value violates unique constraint "pvz_pkey"`) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pvz with such id already exists!"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	createResponse := PvzCreationResponse{Id: created.ID, City: createRequest.City, RegistrationDate: createRequest.RegistrationDate}
	prometheusMetrics.PvzCreated.Inc()
	c.JSON(http.StatusCreated, createResponse)
}

func (pc *PvzController) GetAll(c *gin.Context) {
	startStr := c.Query("startDate")
	endStr := c.Query("endDate")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	// дальше внутренняя логика
	startDate, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startDate, must be RFC3339"})
		return
	}
	endDate, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endDate, must be RFC3339"})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit number"})
		return
	}

	// Получаем все что надо с фильтром
	pvzs, err := pc.pvzRepo.GetAllWithFilter(context.Background(), startDate, endDate, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var out []PvzsResponseStruct
	for _, p := range pvzs {
		// 1) сопоставляем инфу о пвз
		info := PvzInfoResponseStruct{
			Id:               p.ID,
			RegistrationDate: p.RegistrationDate.Format(time.RFC3339),
			City:             p.City,
		}

		// 2) загружаем приёмки
		recs, err := pc.receptionRepo.GetAllByPvzID(c.Request.Context(), p.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "loading receptions: " + err.Error()})
			return
		}

		// 3) сопоставляем приемки
		var recDTOs []ReceptionsInfoResponseStruct
		for _, r := range recs {
			currentProducts, err := pc.productRepo.GetAllByReceptionId(c.Request.Context(), r.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "loading products: " + err.Error()})
				return
			}
			currentProductsStructure := make([]ProductsInfoResponseStruct, 0)
			for _, product := range currentProducts {
				currentProductsStructure = append(currentProductsStructure, ProductsInfoResponseStruct{
					Id:          product.ID,
					DateTime:    product.DateTime.Format(time.RFC3339),
					Type:        product.Type,
					ReceptionId: r.ID,
				})
			}
			// 4) сопоставляем продукты в рамках приёмок в текущем ПВЗ

			recDTOs = append(recDTOs, ReceptionsInfoResponseStruct{
				// 5) и саму текущую приемку в рамках этого ПВЗ сопоставляем
				Reception: ReceptionInfoResponseStruct{
					Id:       r.ID,
					DateTime: r.DateTime.Format(time.RFC3339),
					PvzId:    r.PvzId,
					Status:   r.Status,
				},
				Products: currentProductsStructure,
			})
		}

		out = append(out, PvzsResponseStruct{
			PvzInfo:       info,
			ReceptionInfo: recDTOs,
		})
	}

	response := GetAllFilterResponse{Pvzs: out}
	c.JSON(http.StatusOK, response)
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
