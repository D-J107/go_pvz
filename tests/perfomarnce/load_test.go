package perfomarnce

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestLoad_1100RPS_LatencyUnder100ms(t *testing.T) {
	// конфигурируем цель
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/pvz?startDate=2025-04-15T09:27:05.436Z&endDate=2025-06-15T09:27:05.436Z&page=1&limit=100",
		Header: http.Header{"Authorization": []string{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDc2NDYyMDMsInJvbGUiOiJtb2RlcmF0b3IifQ.r_lsk-ALFElHjUBxmMCGKfJ3glWeD8RIgTHRSlGKNfY"}},
	})

	// конфигурируем скорость запросов
	rate := vegeta.Rate{Freq: 1100, Per: time.Second}
	duration := 5 * time.Second
	attacker := vegeta.NewAttacker()

	// "атакуем" (посылаем запросы)
	var results vegeta.Results
	for res := range attacker.Attack(targeter, rate, duration, "pvz-load-test") {
		results = append(results, *res)
	}

	// собираем метрики
	var metrics vegeta.Metrics
	for i := range results {
		metrics.Add(&results[i])
	}
	metrics.Close()

	// проверяем время
	if metrics.Latencies.Max >= 100*time.Millisecond {
		t.Fatalf("max latency = %v; want < 100ms", metrics.Latencies.Max)
	}
	fmt.Printf("max latency = %v\n", metrics.Latencies.Max)
}
