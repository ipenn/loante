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

//StructToMapReflect 结构体转Map
func StructToMapReflect(in interface{}, tagName string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}

