package helper

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
	"strings"
)

// Parse 用于解析公共参数中间件
func ParseParams(params interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		typ := reflect.TypeOf(&params).Elem()
		val := reflect.ValueOf(&params).Elem()

		for i := 0; i < typ.NumField(); i++ {
			typeField := typ.Field(i)

			inputFieldName := typeField.Tag.Get("form")
			if inputFieldName == "" {
				inputFieldName = typeField.Tag.Get("json")
				if inputFieldName == "" {
					continue
				}
			}

			keys := strings.Split(inputFieldName, ",")
			for _, key := range keys {
				v := pickValue(key, c)
				if v != "" {
					if val.Field(i).Kind() == reflect.Int64 {
						intVal, _ := strconv.ParseInt(v, 10, 64)
						val.Field(i).SetInt(intVal)
					} else if val.Field(i).Kind() == reflect.Float64 {
						floatVal, _ := strconv.ParseFloat(v, 64)
						val.Field(i).SetFloat(floatVal)
					} else {
						val.Field(i).SetString(v)
					}
					break
				}
			}
		}

		c.Set("params", &params)

		c.Next()
	}
}

func pickValue(key string, c *gin.Context) string {
	v := c.Query(key)
	if v != "" {
		return v
	}

	pv := c.PostForm(key)
	if pv != "" {
		return pv
	}

	hv := c.Request.Header.Get(key)
	if hv != "" {
		return hv
	}
	cv, err := c.Request.Cookie(key)
	if err == nil && cv.Value != "" {
		return cv.Value
	}
	return ""
}
