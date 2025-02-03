package util

import (
	"fmt"
	"log"
	"park/config"
	"park/database"
	modelscar "park/models/modelsCar"
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
	now := time.Now().Add(time.Minute).Format(config.TimeFormat)

	if role == string(modelsuser.OperatorRole) {
		var lastLogin modeloperator.Operator

		if err := database.DB.Where("operator = ? ", username).Order("id DESC").First(&lastLogin).Error; err != nil {
			return err
		}
		lastLogin.LogoutAt = now
		if err := database.DB.Save(&lastLogin).Error; err != nil {
			return err
		}
	}

	return nil
}

func CalculateV2(username string, role string) (int, error) {
	now := time.Now().Format(config.TimeFormat)

	var calculations []modelscar.Car_Model
	var totalPayment float64

	if role == string(modelsuser.OperatorRole) {
		if err := database.DB.Where("user_id = ? AND pay_status = true", username).Find(&calculations).Error; err != nil {
			log.Println("Error fetching car models for user:", username, "Error:", err)
			return 0, err
		}

		if len(calculations) == 0 {
			return 0, fmt.Errorf("no records found for user %s", username)
		}

		for _, car := range calculations {
			totalPayment += car.Total_payment
		}

		if err := database.DB.Model(&modelscar.Car_Model{}).
			Where("user_id = ? AND pay_status = true", username).
			Update("pay_status", false).Error; err != nil {
			log.Println("Failed to update paystatus for user:", username, "Error:", err)
			return 0, err
		}

		var operator modeloperator.Operator
		if err := database.DB.Where("operator = ?", username).Order("id DESC").First(&operator).Error; err != nil {
			log.Println("Operator not found for user:", username, "Error:", err)
			return 0, fmt.Errorf("operator not found for user %s", username)
		}

		if err := database.DB.Model(&operator).
			Update("money", totalPayment).Error; err != nil {
			log.Println("Failed to update operator money for user:", username, "Error:", err)
			return 0, err
		}

		operator.LogoutAt = now
		if err := database.DB.Save(&operator).Error; err != nil {
			log.Println("Failed to update operator logout time for user:", username, "Error:", err)
			return 0, err
		}
	}

	return int(totalPayment), nil
}
