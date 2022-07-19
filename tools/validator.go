package tools

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"loante/global"
	"reflect"
	"strings"
)

// CUSTOM VALIDATION RULES =============================================

func Validate(payload interface{}) *fiber.Error {
	err := global.Validate.Struct(payload)

	if err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(
				errors,
				//err.Translate(conf.Config.Trans),
				fmt.Sprintf("`%v` with value `%v` doesn't satisfy the `%v` constraint", err.Field(), err.Value(), err.Tag()),
			)
		}
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: strings.Join(errors, ","),
		}
	}

	return nil
}

//ParseBody Password validation rule: required,min=6,max=100
func ParseBody(ctx *fiber.Ctx, body interface{}) *fiber.Error {
	switch ctx.Method() {
	case "GET":
		if err := ctx.QueryParser(body); err != nil {
			return fiber.ErrBadRequest
		}
		if !reflect.ValueOf(body).IsZero() {
			if err := ctx.BodyParser(body); err != nil {
				//return fiber.ErrBadRequest
			}
		}
		//fmt.Println(body)
		////if err := ctx.BodyParser(body); err != nil {
		////	return fiber.ErrBadRequest
		////}
	default:
		if err := ctx.BodyParser(body); err != nil {
			return fiber.ErrBadRequest
		}
	}
	return Validate(body)
}
