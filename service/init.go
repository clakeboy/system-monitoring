package service

import "system-monitoring/models"

func InitSystem() {
	model := models.NewManagerModel(nil)
	model.InitDB()
}
