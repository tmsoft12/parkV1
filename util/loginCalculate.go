package util

import (
	"fmt"
	"park/config"
	"park/database"
	modelsuser "park/models/modelsUser"
	modeloperator "park/models/operatorModel"
	"time"
)

func LoginMath(username string, role string, park string) error {
	now := time.Now().Format(config.TimeFormat)

	fmt.Printf("User: %s, Role: %s, Park: %s\n", username, role, park)

	login := modeloperator.Operator{
		Operator: username,
		Park:     park,
		LoginAt:  now,
	}
	if role == string(modelsuser.OperatorRole) {
		if err := database.DB.Create(&login).Error; err != nil {
			return err
		}
	}

	return nil
}

func LoginOut(username string, role string) error {
	now := time.Now().Format(config.TimeFormat)
	if role == string(modelsuser.OperatorRole) {
		var lastLogin modeloperator.Operator

		if err := database.DB.Where("operator = ? ", username).Find(&lastLogin).Error; err != nil {
			return err
		}

		lastLogin.LogoutAt = now
		lastLogin.Money = 0
		fmt.Println("salam")
		if err := database.DB.Save(&lastLogin).Error; err != nil {
			return err
		}
	}

	return nil
}
