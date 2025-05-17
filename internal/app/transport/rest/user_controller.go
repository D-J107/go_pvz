package rest

import (
	"context"
	"my_pvz/internal/db"
	postgresql "my_pvz/internal/db/PostgreSQL"
	"my_pvz/internal/logger"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	repo db.UserRepository
}

func NewUserController(db *db.DB) *UserController {
	return &UserController{repo: postgresql.NewPostgesUserRepositoryImpl(db)}
}

func (uc *UserController) DummyLogin(c *gin.Context) {
	var req DummyLoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token := generateToken(req.Role)
	c.Set("role", req.Role)
	c.Header("Authorization", "Bearer "+token)
	response := &DummyLoginResponse{Token: token}
	c.JSON(http.StatusCreated, response)
}

func (uc *UserController) Register(c *gin.Context) {
	logger.Log.Debug("received POST request to create new user")

	var registerRequest RegisterRequest
	if err := c.BindJSON(&registerRequest); err != nil {
		logger.Log.Error("invalid register request body", "error", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), 7)
	if err != nil {
		logger.Log.Error("cant hash password", "error", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user, err := uc.repo.Create(context.Background(), registerRequest.Email, string(hashedPwd), registerRequest.Role)
	if err != nil && err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)` {
		logger.Log.Error("cant create user", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already occupied by another user."})
		return
	}

	resp := &RegisterResponse{Id: user.ID, Email: user.Email, Role: user.Role}
	c.JSON(http.StatusCreated, resp)
	logger.Log.Info("successfully created new user. ", "userId", user.ID)
}

func (uc *UserController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1) проверяем что пользователь с такой почтой существует и что пароли совпадают
	user, err := uc.repo.GetByEmail(context.Background(), req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user with such email does not exists!"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords are not equal!"})
		return
	}

	// 2) формирует JWT токен и отправляем юзеру
	token := generateToken(user.Role)
	c.Set("role", user.Role)
	c.Header("Authorization", "Bearer "+token)
	resp := &LoginResponse{Token: token}
	c.JSON(http.StatusCreated, resp)
}

func generateToken(role string) string {
	secret := os.Getenv("ACCESS_TOKEN_SECRET")
	claims := jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(24 * 2 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := t.SignedString([]byte(secret))
	return signedToken
}
