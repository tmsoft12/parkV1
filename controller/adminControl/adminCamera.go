package admincontrol

import (
	"park/database"
	"park/models/camera"
	"park/util"

	"github.com/gofiber/fiber/v2"
)

// CreateCamera creates a new camera
// @Summary Create a new camera
// @Description Creates a new camera in the database with validation for camera type
// @Tags Cameras
// @Accept json
// @Produce json
// @Param camera body camera.Cameras true "Camera data"
// @Success 201 {object} camera.Cameras
// @Failure 400 {string} string "Invalid camera type"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/cameras [post]
func CreateCamera(c *fiber.Ctx) error {
	var cam camera.Cameras

	if err := c.BodyParser(&cam); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Failed to parse camera data"})
	}

	if !util.IsValidCamera(cam.Type) {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid camera type"})
	}

	if err := database.DB.Create(&cam).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Internal server error"})
	}

	return c.Status(201).JSON(cam)
}

// UpdateCamera updates an existing camera by ID
// @Summary Update a camera by ID
// @Description Updates the camera data in the database using its unique ID
// @Tags Cameras
// @Accept json
// @Produce json
// @Param id path int true "Camera ID"
// @Param camera body camera.Cameras true "Updated camera data"
// @Success 200 {object} camera.Cameras
// @Failure 400 {string} string "Invalid camera data"
// @Failure 404 {string} string "Camera not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/cameras/{id} [put]
func UpdateCamera(c *fiber.Ctx) error {
	id := c.Params("id")
	var datacam camera.Cameras

	if err := database.DB.Where("id = ?", id).First(&datacam).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Camera not found",
		})
	}

	var updateCamera camera.Cameras
	if err := c.BodyParser(&updateCamera); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid camera data",
		})
	}

	if updateCamera.Type != "" && !util.IsValidCamera(updateCamera.Type) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid camera type",
		})
	}

	if updateCamera.Type != "" {
		datacam.Type = updateCamera.Type
	}
	if updateCamera.Name != "" {
		datacam.Name = updateCamera.Name
	}

	if err := database.DB.Save(&datacam).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error while updating camera",
		})
	}

	return c.Status(200).JSON(datacam)
}

// DeleteCamera deletes a camera by ID
// @Summary Delete a camera by ID
// @Description Deletes the camera from the database using its unique ID
// @Tags Cameras
// @Param id path int true "Camera ID"
// @Success 200 {string} string "Camera deleted successfully"
// @Failure 404 {string} string "Camera not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/cameras/{id} [delete]
func DeleteCamera(c *fiber.Ctx) error {
	id := c.Params("id")
	var camera camera.Cameras

	if err := database.DB.Where("id = ?", id).First(&camera).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Camera not found",
		})
	}

	if err := database.DB.Delete(&camera).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error while deleting camera",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Camera deleted successfully",
	})
}

// GetCameraByID retrieves a camera by ID
// @Summary Get a camera by ID
// @Description Retrieves the camera from the database using its unique ID
// @Tags Cameras
// @Param id path int true "Camera ID"
// @Success 200 {object} camera.Cameras
// @Failure 404 {string} string "Camera not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/cameras/{id} [get]
func GetCameraByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var camera camera.Cameras

	if err := database.DB.Where("id = ?", id).First(&camera).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Camera not found",
		})
	}

	return c.Status(200).JSON(camera)
}

// GetCameras retrieves a list of cameras with pagination
// @Summary Get cameras with pagination
// @Description Retrieves a list of cameras from the database with pagination
// @Tags Cameras
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} camera.Cameras
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/cameras [get]
func GetCameras(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	var cameras []camera.Cameras
	var totalCameras int64

	if err := database.DB.Model(&camera.Cameras{}).Count(&totalCameras).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	totalPages := int(totalCameras) / limit
	if totalCameras%int64(limit) > 0 {
		totalPages++
	}

	hasNext := page < totalPages
	hasPrev := page > 1

	if err := database.DB.Order("id desc").Offset((page - 1) * limit).Limit(limit).Find(&cameras).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
		"hasNext":    hasNext,
		"hasPrev":    hasPrev,
		"cameras":    cameras,
	})
}
