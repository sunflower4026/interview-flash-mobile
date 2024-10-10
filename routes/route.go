package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.com/sunflower4026/interview-flash-mobile/common/httpservice"
	userController "gitlab.com/sunflower4026/interview-flash-mobile/controller/user"
	"gitlab.com/sunflower4026/interview-flash-mobile/middleware"
	userRepository "gitlab.com/sunflower4026/interview-flash-mobile/model/repository/user"
	userService "gitlab.com/sunflower4026/interview-flash-mobile/service/user"

	authController "gitlab.com/sunflower4026/interview-flash-mobile/controller/auth"
	authService "gitlab.com/sunflower4026/interview-flash-mobile/service/auth"

	profileController "gitlab.com/sunflower4026/interview-flash-mobile/controller/profile"
	profileService "gitlab.com/sunflower4026/interview-flash-mobile/service/profile"

	transactionController "gitlab.com/sunflower4026/interview-flash-mobile/controller/transaction"
	transactionRepository "gitlab.com/sunflower4026/interview-flash-mobile/model/repository/transaction"
	transactionService "gitlab.com/sunflower4026/interview-flash-mobile/service/transaction"

	jwtService "gitlab.com/sunflower4026/interview-flash-mobile/service/jwt"
)

func Router(ctx context.Context, gin *gin.Engine, cfg *viper.Viper, svc *httpservice.Service) {

	var (
		jwtService = jwtService.NewJWTService(cfg)

		userRespository userRepository.UserRepository = userRepository.NewUserRepository()
		userService     userService.UserService       = userService.NewUserService(userRespository, svc.GetPostgreConf())
		userController  userController.UserController = userController.NewUserController(userService)

		authService    authService.AuthService       = authService.NewAuthService(userRespository, jwtService, svc.GetPostgreConf())
		authController authController.AuthController = authController.NewAuthController(authService)

		profileService    profileService.ProfileService = profileService.NewProfileService(userRespository, svc.GetPostgreConf())
		profileController                               = profileController.NewProfileController(profileService)

		transactionRepository transactionRepository.TransactionRepository = transactionRepository.NewTransactionRepository()
		transactionService    transactionService.TransactionService       = transactionService.NewTransactionService(transactionRepository, svc.GetPostgreConf())
		transactionController transactionController.TransactionController = transactionController.NewTransactionController(transactionService)
	)

	apiVersion := gin.Group("/api/v1")
	{
		authGroup := apiVersion.Group("/auth")
		{
			authGroup.POST("/register", authController.Register)
			authGroup.POST("/login", authController.Login)
		}

		profileGroup := apiVersion.Group("/profile")
		{
			profileGroup.Use(middleware.AuthMiddleware(jwtService, userRespository, svc.GetPostgreConf()))
			profileGroup.PATCH("/update", profileController.Update)
		}

		userGroup := apiVersion.Group("/users")
		{
			userGroup.Use(middleware.AuthMiddleware(jwtService, userRespository, svc.GetPostgreConf()))
			userGroup.GET("", userController.FindAll)
		}

		transactionGroup := apiVersion.Group("/transactions")
		{
			transactionGroup.Use(middleware.AuthMiddleware(jwtService, userRespository, svc.GetPostgreConf()))
			transactionGroup.POST("/topup", transactionController.Topup)
			transactionGroup.POST("/payment", transactionController.Payment)
			transactionGroup.POST("/transfer", transactionController.Transfer)
			transactionGroup.GET("", transactionController.FindAll)
		}
	}
}
