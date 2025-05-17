package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"my_pvz/internal/app/transport/rest"
	"my_pvz/internal/app/transport/rest/middleware"
	"my_pvz/internal/db"
	"my_pvz/internal/logger"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	testDB     *db.DB
	testRouter *gin.Engine
)

func TestMain(m *testing.M) {
	logger.Init()

	os.Setenv("DATABASE_URL", "postgresql://localhost/avito_test_pvz_db?user=avito_test_user&password=test_qdmkio231")
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()
	testDB = db.NewDb(ctx)
	if err := testDB.InitDb(ctx); err != nil {
		panic("failed to initialized test database " + err.Error())
	}

	// Очищаём все перед тестами чтобы можно было каждый раз запускать тесты
	// с вставкой одинаковых данных и не получать ошибку нарушения уникальности
	cleanupQuery := `
		DELETE FROM products;
		DELETE FROM receptions;
		DELETE FROM pvz;
		DELETE FROM users;
	`
	_, err := testDB.Pool.Exec(ctx, cleanupQuery)
	if err != nil {
		panic("failed to clean up test database: " + err.Error())
	}

	gin.SetMode(gin.TestMode)
	testRouter = gin.New()

	// <<__ users __>>
	usersController := rest.NewUserController(testDB)
	testRouter.POST("/dummyLogin", usersController.DummyLogin)
	testRouter.POST("/register", usersController.Register)
	testRouter.POST("/login", usersController.Login)

	// всё что ниже этой линии будет проходить через этот middleware
	testRouter.Use(middleware.AuthMiddleware())

	// <<__ pvz __>>
	pvzController := rest.NewPvzController(testDB)
	testRouter.POST("/pvz", middleware.Authorize("moderator"), pvzController.Create)
	testRouter.GET("/pvz", middleware.Authorize("moderator"), pvzController.GetAll)

	// <<__ Receptions __>>
	receptionsController := rest.NewReceptionController(testDB)
	testRouter.POST("/receptions", middleware.Authorize("employee"), receptionsController.Create)
	testRouter.POST("/pvz/:pvzId/close_last_reception", middleware.Authorize("employee"), receptionsController.CloseLastReception)

	// <<__ Products __>>
	productsController := rest.NewProductController(testDB)
	testRouter.POST("/products", middleware.Authorize("employee"), productsController.Create)
	testRouter.POST("/pvz/:pvzId/delete_last_product", middleware.Authorize("employee"), productsController.DeleteLastProduct)

	os.Exit(m.Run())
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RegisterResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type PvzCreationRequest struct {
	Id               string `json:"id"`
	RegistrationDate string `json:"registrationDate"`
	City             string `json:"city"`
}

type PvzCreationResponse struct {
	Id               string `json:"id"`
	RegistrationDate string `json:"registrationDate"`
	City             string `json:"city"`
}

func TestRegisterEndpoint(t *testing.T) {
	// 1 reg moder
	modEmail := randomEmail()
	modPass := "pass12"
	b, _ := json.Marshal(rest.RegisterRequest{Email: modEmail, Password: modPass, Role: "moderator"})
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(rec, req)
	// 1.1
	assert.Equal(t, http.StatusCreated, rec.Code)

	// 2 login moder
	b, _ = json.Marshal(rest.LoginRequest{Email: modEmail, Password: modPass})
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(rec, req)
	// 2.1
	assert.Equal(t, http.StatusCreated, rec.Code)
	var loginReponse rest.LoginResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &loginReponse))
	modToken := loginReponse.Token
	// fmt.Println("modToken:", modToken)

	//3 create pvz
	pvzId := "3fa85f64-5717-4562-b3fc-2c963f66afa2"
	b, _ = json.Marshal(rest.PvzCreationRequest{Id: pvzId, RegistrationDate: time.Now().Format(time.RFC3339), City: "Казань"})
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/pvz", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+modToken)
	testRouter.ServeHTTP(rec, req)
	// 3.1
	assert.Equal(t, http.StatusCreated, rec.Code)
	var pvzCreateResponse rest.PvzCreationResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &pvzCreateResponse))
	assert.Equal(t, pvzId, pvzCreateResponse.Id)
	assert.Equal(t, "Казань", pvzCreateResponse.City)

	// 4 register employee
	emplEmail := randomEmail()
	emplPass := "pass123"
	b, _ = json.Marshal(rest.RegisterRequest{Email: emplEmail, Password: emplPass, Role: "employee"})
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(rec, req)
	// 4.1
	assert.Equal(t, http.StatusCreated, rec.Code)

	// 5 login employee
	b, _ = json.Marshal(rest.LoginRequest{Email: emplEmail, Password: emplPass})
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(rec, req)
	// 5.1
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &loginReponse))
	emplToken := loginReponse.Token

	// 6 create reception
	b, _ = json.Marshal(rest.ReceptionCreationRequest{PvzId: pvzId})
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/receptions", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+emplToken)
	testRouter.ServeHTTP(rec, req)
	// 6.1
	assert.Equal(t, http.StatusCreated, rec.Code)
	var receptionCreateResponse rest.ReceptionCreationResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &receptionCreateResponse))
	assert.Equal(t, pvzId, receptionCreateResponse.PvzId)
	assert.Equal(t, "in_progress", receptionCreateResponse.Status)

	// 7 add 50 products
	types := []string{"электроника", "обувь", "одежда"}
	typesI := 0
	for i := 0; i < 50; i++ {
		b, _ = json.Marshal(rest.CreateProductRequest{Type: types[typesI], PvzId: pvzId})
		typesI = (typesI + 1) % 3
		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/products", bytes.NewBuffer(b))
		req.Header.Set("Content-type", "application/json")
		req.Header.Set("Authorization", "Bearer "+emplToken)
		testRouter.ServeHTTP(rec, req)
		// 7.1
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	// 8 close reception
	b, _ = json.Marshal(rest.ReceptionCloseRequest{PvzId: pvzId})
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", fmt.Sprintf("/pvz/%s/close_last_reception", pvzId), bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+emplToken)
	testRouter.ServeHTTP(rec, req)
	// 8.1
	assert.Equal(t, http.StatusOK, rec.Code)

	// 9 get all pvz
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/pvz?startDate=2023-04-15T09:27:05.436Z&endDate=2026-06-15T09:27:05.436Z&page=1&limit=100", nil)
	req.Header.Set("Authorization", "Bearer "+modToken)
	testRouter.ServeHTTP(rec, req)
	// 9.1
	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println("GET body:", string(rec.Body.Bytes())) // {"Pvzs":null}

	var listResponse rest.GetAllFilterResponse
	err := json.Unmarshal(rec.Body.Bytes(), &listResponse)
	assert.NoError(t, err)

	// 10 assert pvz
	assert.Len(t, listResponse.Pvzs, 1)
	pvz := listResponse.Pvzs[0]
	assert.Len(t, pvz.ReceptionInfo, 1)
	reception := pvz.ReceptionInfo[0]
	assert.Equal(t, "close", reception.Reception.Status)
	assert.Len(t, reception.Products, 50)
}

func randomEmail() string {
	return fmt.Sprintf("user_%d@example.com", rand.Int63())
}

func randomString(length int) string {
	charset := "qwertyuiopasdfghjklzxcvbnmASDFGHJKLZXCVBNMQWERTYUIO"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
