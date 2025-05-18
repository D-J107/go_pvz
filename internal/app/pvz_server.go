package app

import (
	"context"
	"fmt"
	"log"
	pvz "my_pvz/internal/app/transport/grpc"
	prometheusMetrics "my_pvz/internal/app/transport/rest/prometheus_metrics"
	"my_pvz/internal/db"
	"my_pvz/internal/logger"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func RunPvzServer() {
	// Logger
	logger.Init()

	// DB
	db := db.NewDb(context.Background())
	db.InitDb(context.Background())
	logger.Log.Info("✅ Database connection successfully established and tables schema defined.")

	router := SetupRoutes(db)
	// REST
	go func() {
		fmt.Println("✅ HTTP serves 8080 port")
		if err := router.Run(":8080"); err != nil {
			fmt.Println("HTTP server error", "err", err)
		}
	}()

	// gRPC
	pvzsHandler := pvz.NewPvzHandler(db)
	grpcServer := grpc.NewServer()
	pvz.RegisterPvzServiceServer(grpcServer, pvzsHandler)
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic("cant establish tcp connection:" + err.Error())
	}
	logger.Log.Info("✅ gRPC server listening on port 3000")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	// prometheus
	prometheusMetrics.Init()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Log.Info("✅ Prometheus metrics server running on 9000 port")
		if err = http.ListenAndServe(":9000", nil); err != nil {
			log.Fatalf("Prometheus server error: %v", err)
		}
	}()

	select {} // лочим текущую горутину чтобы сервер не завершал работу
}
