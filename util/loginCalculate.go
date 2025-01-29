package util

import (
	"fmt"
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
	now := time.Now().Format(config.TimeFormat)

	if role == string(modelsuser.OperatorRole) {
		var lastLogin modeloperator.Operator

		if err := database.DB.Where("operator = ? ", username).Order("id DESC").First(&lastLogin).Error; err != nil {
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

func CalculateTotalPayment(user string, start string, end string) (float64, []modelscar.Car_Model, error) {
	var cars []modelscar.Car_Model
	query := database.DB

	// Kullanıcının arabalarını getir
	if err := query.Where("user_id = ?", user).Order("id DESC").Find(&cars).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to fetch cars: %v", err)
	}

	// Tarih aralığı varsa filtrele
	if start == "" && end == "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", start)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid start time format")
		}

		endTime, err := time.Parse("2006-01-02 15:04:05", end)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid end time format")
		}

		// Tarih aralığına göre arabaları filtrele
		if err := query.Where("user_id = ? AND start_time >= ? AND end_time <= ?", user, startTime, endTime).Find(&cars).Error; err != nil {
			return 0, nil, fmt.Errorf("failed to fetch cars: %v", err)
		}
	}

	// Toplam ödemeyi hesapla
	totalPayment := 0.0
	for _, car := range cars {
		totalPayment += car.Total_payment
	}

	// Para modelini güncelle
	var money modeloperator.Operator
	// Veritabanında kullanıcıya ait bir "Operator" kaydı var mı diye sorgula
	if err := database.DB.Where("operator = ?", user).First(&money).Error; err != nil {
		return 0, nil, fmt.Errorf("operator record not found for user %s", user)
	}

	// Eğer kayıt varsa, sadece "money" değerini güncelle
	if err := database.DB.Model(&money).Update("money", int(totalPayment)).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to update total payment in database: %v", err)
	}

	// Sonuçları döndür
	return totalPayment, cars, nil
}
