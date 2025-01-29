package accountant

import (
	"fmt"
	"park/database"
	modelscar "park/models/modelsCar"
	modeloperator "park/models/operatorModel"

	"time"

	"github.com/gofiber/fiber/v2"
)

// CalculateMoney handles the request to calculate money based on car entry/exit times
// @Summary Calculate cars based on start and end time
// @Description Fetch cars that are within the specified time range. 2025-01-29 13:07:31 2025-01-29 14:09:19
// @Tags Accountant
// @Accept json
// @Produce json
// @Param start query string false "Start Time" example("2025-01-29 10:13:51")
// @Param end query string false "End Time" example("2025-01-29 12:15:10")
// @Success 200 {array} modelscar.Car_Model "List of cars"
// @Router /api/v1/accountant/calculateMoney [get]
func CalculateMoney(c *fiber.Ctx) error {
	start := c.Query("start")
	end := c.Query("end")
	user := c.Locals("username")
	var cars []modelscar.Car_Model
	query := database.DB

	if start != "" && end != "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", start)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid start time format"})
		}

		endTime, err := time.Parse("2006-01-02 15:04:05", end)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid end time format"})
		}

		startTimeStr := startTime.Format("2006-01-02 15:04:05")
		endTimeStr := endTime.Format("2006-01-02 15:04:05")

		query = query.Where("start_time >= ? AND end_time <= ?", startTimeStr, endTimeStr)
	}

	if err := query.Where("user_id = ?", user).Order("id DESC").Find(&cars).Error; err != nil {
		fmt.Println("Database error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch cars"})
	}
	var totalPayment float64

	for _, car := range cars {
		totalPayment += car.Total_payment
	}
	fmt.Println(user)
	return c.JSON(fiber.Map{
		"cars":          cars,
		"total_payment": totalPayment,
	})
}

// @Summary Get all operators with pagination
// @Description Retrieve a list of operators with pagination support
// @Tags Accountant
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} []modeloperator.Operator
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/accountant/operators [get]
func GetOperators(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	pageNum := 1
	limitNum := 10
	fmt.Sscanf(page, "%d", &pageNum)
	fmt.Sscanf(limit, "%d", &limitNum)

	var operators []modeloperator.Operator
	var totalRecords int64

	database.DB.Model(&modeloperator.Operator{}).Count(&totalRecords)

	offset := (pageNum - 1) * limitNum
	database.DB.Offset(offset).Limit(limitNum).Find(&operators)

	totalPages := (totalRecords + int64(limitNum) - 1) / int64(limitNum)
	hasNext := pageNum < int(totalPages)
	hasPrev := pageNum > 1

	return c.JSON(fiber.Map{
		"page":       pageNum,
		"limit":      limitNum,
		"totalPages": totalPages,
		"hasNext":    hasNext,
		"hasPrev":    hasPrev,
		"data":       operators,
	})
}
