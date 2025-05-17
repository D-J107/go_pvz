package pvz

import (
	"context"
	"my_pvz/internal/db"
	postgresql "my_pvz/internal/db/PostgreSQL"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PvzHandler struct {
	UnimplementedPvzServiceServer
	pvzRepo db.PvzRepository
}

func NewPvzHandler(db *db.DB) *PvzHandler {
	return &PvzHandler{pvzRepo: postgresql.NewPostgresPvzRepositoryImpl(db)}
}

func (h *PvzHandler) GetAllPvzs(ctx context.Context, req *GetAllPvzsRequest) (*GetAllPvzsResponse, error) {
	pvzs, err := h.pvzRepo.GetAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	pvzsResponseInfo := make([]*Pvz, 0)
	for _, pvz := range pvzs {
		pvzsResponseInfo = append(pvzsResponseInfo, &Pvz{
			Id:               pvz.ID,
			RegistrationDate: pvz.RegistrationDate.Format(time.RFC3339),
			City:             pvz.City,
		})
	}
	return &GetAllPvzsResponse{Pvzs: pvzsResponseInfo}, nil
}
