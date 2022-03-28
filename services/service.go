package services

import (
	"reflect"
	"sort"
	"strings"

	"github.com/shinhagunn/Shop-Watches/backend/models"
)

func UpperFirstLetter(s string) string {
	return strings.Title(strings.ToLower(s))
}

func getValueByKeyFloat(object interface{}, key string) float64 {
	return reflect.ValueOf(object).FieldByName(UpperFirstLetter(key)).Float()
}

func getValueByKeyString(object interface{}, key string) string {
	return reflect.ValueOf(object).FieldByName(UpperFirstLetter(key)).String()
}

func getValueByKeyInt(object interface{}, key string) int64 {
	return reflect.ValueOf(object).FieldByName(UpperFirstLetter(key)).Int()
}

func SortProductByOrderby(products []models.Product, order string, orderby string) []models.Product {
	result := products
	if order == "asc" {
		if reflect.ValueOf(result[0]).FieldByName(UpperFirstLetter(orderby)).Kind() == reflect.Float64 {
			sort.SliceStable(result, func(i int, j int) bool {
				return getValueByKeyFloat(result[i], orderby) > getValueByKeyFloat(result[j], orderby)
			})
		}
		if reflect.ValueOf(result[0]).FieldByName(UpperFirstLetter(orderby)).Kind() == reflect.String {
			sort.SliceStable(result, func(i int, j int) bool {
				return getValueByKeyString(result[i], orderby) > getValueByKeyString(result[j], orderby)
			})
		}
		if reflect.ValueOf(result[0]).FieldByName(UpperFirstLetter(orderby)).Kind() == reflect.Int64 {
			sort.SliceStable(result, func(i int, j int) bool {
				return getValueByKeyInt(result[i], orderby) > getValueByKeyInt(result[j], orderby)
			})
		}
	}
	if order == "desc" {
		if reflect.ValueOf(result[0]).FieldByName(UpperFirstLetter(orderby)).Kind() == reflect.Float64 {
			sort.SliceStable(result, func(i int, j int) bool {
				return getValueByKeyFloat(result[i], orderby) < getValueByKeyFloat(result[j], orderby)
			})
		}
		if reflect.ValueOf(result[0]).FieldByName(UpperFirstLetter(orderby)).Kind() == reflect.String {
			sort.SliceStable(result, func(i int, j int) bool {
				return getValueByKeyString(result[i], orderby) < getValueByKeyString(result[j], orderby)
			})
		}
		if reflect.ValueOf(result[0]).FieldByName(UpperFirstLetter(orderby)).Kind() == reflect.Int64 {
			sort.SliceStable(result, func(i int, j int) bool {
				return getValueByKeyInt(result[i], orderby) < getValueByKeyInt(result[j], orderby)
			})
		}
	}

	return result
}
