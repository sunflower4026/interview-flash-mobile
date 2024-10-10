package ginkit

import (
	"context"
	"strings"

	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/web"

	"github.com/gin-gonic/gin"
	en_locale "github.com/go-playground/locales/en"
	id_locale "github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

// Validate validates request body for incoming gin.Context request
// returns web.HTTPError contains field errors (if any).
func Validate(c *gin.Context, req interface{}) *web.HTTPError {
	if err := c.ShouldBind(req); err != nil {
		return web.NewHTTPValidationError(c.Request.Context(), err)
	}

	return nil
}

// ValidatorTranslatorMiddleware adds request body validator's translator
// based on 'Accept-Language' header.
// currently only supports ID & EN locale.
func ValidatorTranslatorMiddleware(v *validator.Validate) gin.HandlerFunc {
	en := en_locale.New()
	id := id_locale.New()
	uni := ut.New(en, en, id)

	transEN, _ := uni.GetTranslator("en")
	transID, _ := uni.GetTranslator("id")

	_ = en_translations.RegisterDefaultTranslations(v, transEN)
	_ = id_translations.RegisterDefaultTranslations(v, transID)

	return func(c *gin.Context) {
		trans := transEN

		lang := c.GetHeader("Accept-Language")
		if strings.ToLower(lang) == "id" ||
			strings.ToLower(lang) == "id-id" {
			trans = transID
		}

		rCtx := context.WithValue(c.Request.Context(), web.ContextKeyTranslator, trans)
		req := c.Request.WithContext(rCtx)
		c.Request = req

		c.Next()
	}
}
