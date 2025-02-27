package testz

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ZReport struct {
	Total_payment int    `json:"total_payment"`
	Username      string `json:"username"`
}

// CreateTarif godoc
// @Summary Create a New Tarif
// @Description Creates a new Report and saves it to the database.
// @Tags Zreport
// @Accept json
// @Produce json
// @Param tarif body ZReport true "Tarif details to be created"
// @Success 201 {object} ZReport "Successfully created"
// @Failure 400 {object} resmodel.ErrorResponse "Invalid request data"
// @Failure 500 {object} resmodel.ErrorResponse "Failed to save data to the database"
// @Router /zreport [post]
func GetZdata(c *fiber.Ctx) error {
	var data ZReport
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON("Prese err")
	}
	fmt.Println("payment:", data.Total_payment)
	fmt.Println("Username:", data.Username)
	return c.JSON(data)
}
