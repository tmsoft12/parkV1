package pdfGenerator

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
)

type ParkInfo struct {
	Operator  string  `json:"operator"`
	Park      string  `json:"park"`
	Money     float64 `json:"money"`
	EntryTime string  `json:"entrytime"`
	ExitTime  string  `json:"exittime"`
}

type RequestData struct {
	Data        []ParkInfo `json:"data"`
	CashierName string     `json:"cashier_name"`
}

func CreatePDF(c *fiber.Ctx) error {
	var requestData RequestData

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	logo := "assets/tm.jpg"
	pdf.Image(logo, 10, 0, 60, 30, false, "", 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetXY(150, 10)
	pdf.Cell(200, 10, fmt.Sprintf("Sene: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(200, 10, "Parkowka Maglumatlary")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Operator")
	pdf.Cell(40, 10, "Parkowka")
	pdf.Cell(40, 10, "Pul Mukdary")
	pdf.Cell(40, 10, "Giren Wagty")
	pdf.Cell(40, 10, "Çykan Wagty")
	pdf.Ln(10)

	totalMoney := 0.0
	pdf.SetFont("Arial", "", 10)

	for _, item := range requestData.Data {
		pdf.Cell(40, 10, item.Operator)
		pdf.Cell(40, 10, item.Park)
		pdf.Cell(40, 10, fmt.Sprintf("%.2f", item.Money))
		pdf.Cell(40, 10, item.EntryTime)
		pdf.Cell(40, 10, item.ExitTime)
		pdf.Ln(10)
		totalMoney += item.Money
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(200, 10, fmt.Sprintf("Jemi: %.2f TMT", totalMoney))
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	if requestData.CashierName == "" {
		pdf.Cell(200, 10, "Kassir: ______________________")
	} else {
		pdf.Cell(200, 10, fmt.Sprintf("Kassir: %s", requestData.CashierName))
	}

	pdfFilePath := "./output.pdf"
	pdf.OutputFileAndClose(pdfFilePath)

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=./pdf/output.pdf")
	c.SendFile(pdfFilePath)

	if err := pdf.Output(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating PDF")
	}

	return nil
}
