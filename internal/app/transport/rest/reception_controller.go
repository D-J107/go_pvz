package rest

import (
	prometheusMetrics "my_pvz/internal/app/transport/rest/prometheus_metrics"
	"my_pvz/internal/db"
	postgresql "my_pvz/internal/db/PostgreSQL"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ReceptionController struct {
	receptionRepo db.ReceptionRepository
}

func NewReceptionController(db *db.DB) *ReceptionController {
	return &ReceptionController{receptionRepo: postgresql.NewPostgesReceptionRepositoryImpl(db)}
}

func (rc *ReceptionController) Create(c *gin.Context) {
	var req ReceptionCreationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lastReception, err := rc.receptionRepo.GetLastByPvzID(c.Request.Context(), req.PvzId)
	if err == nil && lastReception != nil && lastReception.Status == "in_progress" {
		// у этого ПВЗ уже есть текущая незакрытая приёмка, => кидаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "pvz has already in-progress reception!"})
		return
	}
	created, err := rc.receptionRepo.Create(c.Request.Context(), req.PvzId, "in_progress")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pvz with such ID does not exists!"})
		return
	}
	response := ReceptionCreationResponse{
		ReceptionId: created.ID,
		DateTime:    created.DateTime.Format(time.RFC3339),
		PvzId:       created.PvzId,
		Status:      created.Status,
	}
	prometheusMetrics.ReceptionsCreated.Inc()
	c.JSON(http.StatusCreated, response)
}

func (rc *ReceptionController) CloseLastReception(c *gin.Context) {
	pvzID := c.Param("pvzId")
	if pvzID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing pvzId parameter!"})
		return
	}
	lastReception, err := rc.receptionRepo.GetLastByPvzID(c.Request.Context(), pvzID)
	if err != nil && err.Error() == "no rows in result set" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "selected PVZ does not have any Receptions yet"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if lastReception.Status == "close" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "last Reception already closed!"})
		return
	}
	updatingErr := rc.receptionRepo.UpdateReceptionStatus(c.Request.Context(), lastReception.ID, "close")
	if updatingErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updatingErr.Error()})
		return
	}

	resp := ReceptionCloseResponse{
		ReceptionId: lastReception.ID,
		DateTime:    lastReception.DateTime.Format(time.RFC3339),
		PvzId:       lastReception.PvzId,
		Status:      lastReception.Status,
	}
	c.JSON(http.StatusOK, resp)
}
