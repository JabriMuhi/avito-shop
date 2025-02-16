package tools

import "avito-shop/src/models"

func SliceToMap(merchSlice []models.Merch) map[string]models.Merch {
	merchMap := make(map[string]models.Merch)
	for _, merch := range merchSlice {
		merchMap[merch.Name] = merch
	}
	return merchMap
}
