package operator

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"park/database"
	modelscar "park/models/modelsCar"
)

const statusExited = "Exited"

type UpdateCarResponse struct {
	Message string              `json:"message"`
	Car     modelscar.Car_Model `json:"car"`
}

// GetCars godoc
// @Summary Get list of cars
// @Description Get list of cars with pagination
// @Tags cars
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(5)
// @Success 200 {object} GetCarsResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/getallcars [get]
func GetCars(c *fiber.Ctx) error {
	var cars []modelscar.Car_Model
	var totalCount int64
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "5")
	parkno := c.Locals("parkno")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid page number",
		})
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid limit number",
		})
	}

	query := database.DB.Model(&modelscar.Car_Model{})
	if parkno != "" {
		query = query.Where("park_no = ?", parkno)
	}

	query.Count(&totalCount)
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	offset := (page - 1) * limit
	query.Order("id desc").Limit(limit).Offset(offset).Find(&cars)

	if len(cars) == 0 {
		cars = []modelscar.Car_Model{}
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	for i := range cars {
		cars[i].Image_Url = fmt.Sprintf("http://%s:%s/plate/%s", ip, port, cars[i].Image_Url)
	}
	return c.Status(200).JSON(fiber.Map{
		"cars":       cars,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
		"hasNext":    hasNext,
		"hasPrev":    hasPrev,
	})
}

// GetCar godoc
// @Summary Get a car by ID
// @Description Get a car by ID
// @Tags cars
// @Accept  json
// @Produce  json
// @Param id path int true "Car ID"
// @Success 200 {object} modelscar.Car_Model
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/getcar/{id} [get]
func GetCar(c *fiber.Ctx) error {
	id := c.Params("id")
	var car modelscar.Car_Model
	database.DB.Where("id = ?", id).First(&car)
	if car.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Car not found",
		})
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	car.Image_Url = fmt.Sprintf("http://%s:%s/plate/%s", ip, port, car.Image_Url)

	c.Status(200)
	return c.JSON(car)
}

// UpdateCar godoc
// @Summary Update a car by plate number
// @Description Updates a car's status and calculates payment and duration based on start and end times.
// @Tags cars
// @Accept  json
// @Produce  json
// @Param plate path string true "Car plate number"
// @Param car body modelscar.CarUpdate true "Car details to update"
// @Success 200 {object} map[string]interface{} "Updated car details"
// @Failure 400 {object} ErrorResponse "Car already exited or invalid request"
// @Failure 404 {object} ErrorResponse "Car not found"
// @Failure 500 {object} ErrorResponse "Error parsing time"
// @Router /api/v1/camera/updatecar/{plate} [put]
func UpdateCar(c *fiber.Ctx) error {
	plate := c.Params("plate")
	userIDVal := c.Locals("username")

	var car modelscar.Car_Model
	if err := database.DB.Order("id desc").Where("car_number = ?", plate).First(&car).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Car not found", "error": err.Error()})
	}
	if car.Status == statusExited {
		return c.Status(400).JSON("Car already Exited")
	}

	var updatedCar modelscar.Car_Model
	if err := c.BodyParser(&updatedCar); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request", "error": err.Error()})
	}

	updatedCar.Status = statusExited

	if updatedCar.Reason == "" {
		updatedCar.Reason = "Toleg edildi"
		updatedCar.Total_payment = car.Total_payment
	} else {
		updatedCar.Total_payment = 0
	}

	car.Reason = updatedCar.Reason
	car.Total_payment = updatedCar.Total_payment
	updatedCar.End_time = car.End_time

	userID, ok := userIDVal.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error - Invalid user ID type",
		})
	}
	updatedCar.User_id = userID

	if err := database.DB.Model(&car).Updates(map[string]interface{}{
		"reason":        updatedCar.Reason,
		"total_payment": updatedCar.Total_payment,
		"user_id":       updatedCar.User_id,
		"status":        updatedCar.Status,
		"end_time":      updatedCar.End_time,
	}).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Database update failed", "error": err.Error()})
	}

	updatedCar.ID = car.ID
	updatedCar.Car_number = car.Car_number
	updatedCar.Start_time = car.Start_time
	updatedCar.ParkNo = car.ParkNo
	updatedCar.End_time = car.End_time
	updatedCar.Image_Url = car.Image_Url

	return c.Status(200).JSON(fiber.Map{
		"message": "Car updated successfully",
		"car":     updatedCar},
	)
}

// SearchCar godoc
// @Summary Search for cars
// @Description Retrieve a paginated list of cars with optional filtering by car number, enter time, end time, park number, and status.
// @Tags cars
// @Accept json
// @Produce json
// @Param car_number query string false "Filter by car plate number (partial match allowed)"
// @Param enter_time query string false "Filter by enter time (YYYY-MM-DD)"
// @Param end_time query string false "Filter by end time (YYYY-MM-DD)"
// @Param parkno query string false "Filter by parking spot number"
// @Param status query string false "Filter by car status (Inside, Exited)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(5)
// @Success 200 {object} GetCarsResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/searchcar [get]
func SearchCar(c *fiber.Ctx) error {
	var cars []modelscar.Car_Model
	var totalCount int64

	carNumber := c.Query("car_number")
	enterTime := c.Query("enter_time")
	endTime := c.Query("end_time")
	parkNo := c.Locals("parkno")
	status := c.Query("status")
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "5")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid page number"})
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid limit number"})
	}

	query := database.DB.Model(&modelscar.Car_Model{})

	if carNumber != "" {
		query = query.Where("car_number LIKE ?", "%"+carNumber+"%")
	}
	if enterTime != "" {
		if _, err := time.Parse("2006-01-02", enterTime); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Invalid enter_time format. Use YYYY-MM-DD."})
		}
		query = query.Where("DATE(start_time) = ?", enterTime)
	}
	if endTime != "" {
		if _, err := time.Parse("2006-01-02", endTime); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Invalid end_time format. Use YYYY-MM-DD."})
		}
		query = query.Where("DATE(end_time) = ?", endTime)
	}
	if parkNo != "" {
		query = query.Where("park_no = ?", parkNo)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Error counting cars", "error": err.Error()})
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1
	offset := (page - 1) * limit

	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&cars).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Error retrieving cars", "error": err.Error()})
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	for i := range cars {
		cars[i].Image_Url = fmt.Sprintf("http://%s:%s/plate/%s", ip, port, cars[i].Image_Url)
	}
	return c.Status(200).JSON(GetCarsResponse{
		Cars:       cars,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	})
}

type GetCarsResponse struct {
	Cars       []modelscar.Car_Model `json:"cars"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"totalPages"`
	HasNext    bool                  `json:"hasNext"`
	HasPrev    bool                  `json:"hasPrev"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
