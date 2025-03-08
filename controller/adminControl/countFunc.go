package admincontrol

import (
	"park/database"
	"park/models/camera"
	modelscar "park/models/modelsCar"
	modelsuser "park/models/modelsUser"

	"github.com/gofiber/fiber/v2"
)

const statusInside = "Inside"

// UsersCount retrieves the total number of users and the count of users by role
// @Summary Get total number of users and the count of users by role
// @Description Retrieves the count of all users and optionally filtered by role
// @Tags Users Count
// @Produce json
// @Router /api/v1/userCount [get]
func UsersCount(c *fiber.Ctx) error {
	role := "operator"

	var totalCount int64
	var roleCount int64
	var totalCars int64
	var cameraCount int64
	var cars []modelscar.Car_Model

	if err := database.DB.Where("status = ?", statusInside).Find(&cars).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error retrieving cars with status 'Inside'",
		})
	}

	if err := database.DB.Model(&modelsuser.User{}).Count(&totalCount).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error counting total users",
		})
	}

	if err := database.DB.Model(&modelsuser.User{}).Where("role = ?", role).Count(&roleCount).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error counting users by role",
		})
	}
	if err := database.DB.Model(&camera.CamFix{}).Count(&cameraCount).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error counting camera",
		})
	}

	totalCars = int64(len(cars))

	return c.Status(200).JSON(fiber.Map{
		"totalUsers": totalCount,
		"operator":   roleCount,
		"totalCars":  totalCars,
		"camera":     cameraCount,
	})
}
