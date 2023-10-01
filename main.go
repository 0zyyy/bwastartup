package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// dsn := "root:root@tcp(127.0.0.1:3306)/bwastart?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:@tcp(127.0.0.1:3306)/bwastart?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	//Migrate db
	// check if table exist
	db.AutoMigrate(&user.User{}, &campaign.Campaign{}, &campaign.CampaignImage{}, &transaction.Transaction{})

	// repositories
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	//services
	campaignService := campaign.NewService(campaignRepository)
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository)

	//handlers
	campaignHandler := handler.NewCampaignHandler(campaignService)
	userHandler := handler.NewUserHandler(userService, authService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	//
	router := gin.Default()
	router.Static("/images", "./img")
	api := router.Group("/api/v1")

	// users
	api.POST("/user", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email-checkers", userHandler.CheckEmail)
	api.POST("/avatar", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.SaveCampaign)

	// campaigns
	api.GET("/campaigns", authMiddleware(authService, userService), campaignHandler.GetCampaigns)
	api.GET("/campaign/:id", campaignHandler.GetCampaign)
	api.PUT("/campaign/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	//save images
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadCampaignImage)

	// transaction [TODO]
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		var tokenString string
		token := strings.Split(authHeader, " ")
		if len(token) == 2 {
			tokenString = token[1]
		}
		validatedToken, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claim, ok := validatedToken.Claims.(jwt.MapClaims)
		if !ok || !validatedToken.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userId := int(claim["user_id"].(float64))
		user, err := userService.GetUserById(userId)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		c.Set("currentUser", user)
	}
}
