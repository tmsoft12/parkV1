package getdata

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	resmodel "park/controller/getdata/resModel"
	"park/controller/operator"
	"park/database"
	"park/models/camera"
	modelscar "park/models/modelsCar"
)

const (
	statusInside    = "Inside"
	statusExited    = "Exited"
	statusPending   = "Pending"
	statusUnpaid    = "Unpaid"
	timeFormat      = "2006-01-02 15:04:05"
	defaultImageURL = "testPhoto.jpg"
)

// CreateCarEntry handles the creation of a car entry in the parking lot
// @Summary Create a new car entry in the parking lot
// @Description {"ChannelName": "P41", "EventComment": "BE5084AG"}
// @Tags Car Entry
// @Accept json
// @Produce json
// @Param request body camera.CapturedEventData true "Captured data from the camera"
// @Success 201 {object} resmodel.Response "Car entry created successfully"
// @Failure 400 {object} resmodel.ErrorResponse "Bad request, car is already inside"
// @Failure 500 {object} resmodel.ErrorResponse "Internal server error, failed to save data"
// @Router /api/v1/camera/getdata [post]
func CreateCarEntry(c *fiber.Ctx) error {
	var capturedData camera.CapturedEventData
	var carData modelscar.Car_Model

	if err := c.BodyParser(&capturedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resmodel.ErrorResponse{
			Error:   "Failed to parse request body",
			Details: err.Error(),
		})
	}

	now := time.Now().Format(timeFormat)
	carData.ParkNo = capturedData.ChannelName[:2]
	carData.Car_number = capturedData.EventComment
	carData.ParkNo = capturedData.ChannelName
	carData.Status = statusInside
	carData.Start_time = now
	carData.Image_Url = defaultImageURL
	carData.Reason = "Girdi"
	var existingCar modelscar.Car_Model
	err := database.DB.Order("id desc").First(&existingCar, "car_number = ? AND status = ?", carData.Car_number, statusInside).Error
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Car is already inside the parking lot",
		})
	}

	if err := database.DB.Create(&carData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(resmodel.ErrorResponse{
			Error:   "Failed to save data to the database",
			Details: err.Error(),
		})
	}
	operator.Refresh <- struct{}{}

	return c.Status(fiber.StatusCreated).JSON(resmodel.Response{
		Message: "Car entry created successfully",
		Data:    carData,
	})
}

// CreateCarExit handles the car exit process from the parking lot
// @Summary Create a car exit record in the parking lot
// @Description {"EventComment": "BE5084AG"}
// @Tags Car Entry
// @Accept json
// @Produce json
// @Param request body camera.CapturedEventData true "Captured data from the camera"
// @Success 200 {object} resmodel.Response "Car exit updated successfully"
// @Failure 400 {object} resmodel.ErrorResponse "Bad request, car already exited"
// @Failure 404 {object} resmodel.ErrorResponse "Car not found"
// @Failure 500 {object} resmodel.ErrorResponse "Internal server error, failed to update data"
// @Router /api/v1/camera/getdata [put]
func CreateCarExit(c *fiber.Ctx) error {
	var capturedData camera.CapturedEventData

	if err := c.BodyParser(&capturedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	var carData modelscar.Car_Model
	if err := database.DB.Where("car_number = ?", capturedData.EventComment).Order("id desc").First(&carData).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Car not found",
			"error":   err.Error(),
		})
	}

	if carData.Status == statusExited {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Car already Exited",
		})
	}

	startTimeStr := carData.Start_time
	endTimeStr := time.Now().Format(timeFormat)
	startTime, _ := time.Parse(timeFormat, startTimeStr)

	endTime, _ := time.Parse(timeFormat, endTimeStr)
	duration := endTime.Sub(startTime)
	minutes := duration.Minutes()
	carData.Duration = int(minutes)
	carData.Status = statusPending
	carData.End_time = endTimeStr
	carData.Reason = "Garasylyar"

	if minutes <= 360 {
		carData.Total_payment = 2
	} else if minutes <= 1440 {

		carData.Total_payment = 3
	} else {

		days := int(math.Ceil(minutes / 1440))
		carData.Total_payment = float64(3 * days)
	}
	if err := database.DB.Model(&carData).Updates(carData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database update failed",
			"error":   err.Error(),
		})
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	carData.Image_Url = fmt.Sprintf("http://%s:%s/plate/%s", ip, port, carData.Image_Url)

	operator.Broadcast <- carData

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Car exit updated successfully",
		"car":     carData,
	})
}
