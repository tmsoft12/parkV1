package camfix

import (
	"fmt"
	"park/database"
	"park/models/camera"

	"github.com/gofiber/fiber/v2"
)

// Camfix godoc
// @Summary Create a New CamFix
// @Description Creates a new cam and saves it to the database.
// @Tags CamFix
// @Accept json
// @Produce json
// @Param cam body camera.CamFix true "Cam details to be created"
// @Success 201 {object} camera.CamFix "Successfully created"
// @Router /api/v1/addcam [post]
func AddCam(c *fiber.Ctx) error {
	var fix camera.CamFix
	if err := c.BodyParser(&fix); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Can not parse data",
		})
	}

	// ChannelName ile mevcut bir kayıt var mı kontrol et
	var existing camera.CamFix
	if err := database.DB.Where("channel_name = ?", fix.ChannelName).First(&existing).Error; err == nil {
		// Eğer kayıt bulunursa, hata döndür
		return c.Status(409).JSON(fiber.Map{
			"message": fmt.Sprintf("A camera with ChannelName '%s' already exists", fix.ChannelName),
		})
	}

	// Kayıt yoksa, yeni kaydı ekle
	if err := database.DB.Create(&fix).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Can not create",
		})
	}

	fmt.Println(fix)
	return c.Status(201).JSON(fix)
}
