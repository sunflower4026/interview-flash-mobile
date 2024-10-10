package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/sunflower4026/interview-flash-mobile/common/constants"
	"gitlab.com/sunflower4026/interview-flash-mobile/common/httpservice"
	"gorm.io/gorm"

	userRepository "gitlab.com/sunflower4026/interview-flash-mobile/model/repository/user"
	userWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/user"
	jwtService "gitlab.com/sunflower4026/interview-flash-mobile/service/jwt"
)

// AuthMiddleware checks if the user is logged in
func AuthMiddleware(JWTService jwtService.JWTService, userRepository userRepository.UserRepository, DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			httpservice.ResponseData(c, nil, errors.New(httpservice.ERR_UNAUTHORIZED))
			c.Abort()
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			httpservice.ResponseData(c, nil, errors.New(httpservice.ERR_UNAUTHORIZED))
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		// Validate the token
		parsedToken, err := JWTService.ValidateToken(token, true) // Assuming it is an access token
		if err != nil {
			httpservice.ResponseData(c, nil, errors.New(httpservice.ERR_UNAUTHORIZED))
			c.Abort()
			return
		}

		if claims, ok := parsedToken.Claims.(*jwtService.JWTCustomClaim); ok && parsedToken.Valid {
			user, err := userRepository.Find(c.Request.Context(), DB, fmt.Sprintf("id = '%s'", claims.UserID))
			if err != nil {
				httpservice.ResponseData(c, nil, errors.New(httpservice.ERR_UNAUTHORIZED))
				c.Abort()
				return
			}

			c.Set(constants.ContextKeyUserID, claims.UserID)
			c.Set(constants.ContextKeyUserDetail, userWeb.ToUserResponse(user))
			ctx := context.WithValue(c.Request.Context(), constants.ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, constants.ContextKeyUserDetail, userWeb.ToUserResponse(user))
			c.Request = c.Request.WithContext(ctx)
		} else {
			httpservice.ResponseData(c, nil, errors.New(httpservice.ERR_UNAUTHORIZED))
			c.Abort()
			return
		}

		c.Next()
	}
}
