package handlers

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"example.com/e/models"
	"github.com/gofiber/fiber/v2"
)

// GetReceiptPoints handles GET /receipts/{id}/points
func GetReceiptPoints(c *fiber.Ctx) error {
	id := c.Params("id")
	receipt, ok := Receipts[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Receipt not found"})
	}
	return c.JSON(fiber.Map{
		"points": calculatePoints(receipt),
	})
}

// **

// These rules collectively define how many points should be awarded to a receipt.
// * One point for every alphanumeric character in the retailer name.
// * 50 points if the total is a round dollar amount with no cents.
// * 25 points if the total is a multiple of `0.25`.
// * 5 points for every two items on the receipt.
// * If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2` and round up to the nearest integer. The result is the number of points earned.
// * 6 points if the day in the purchase date is odd.
// * 10 points if the time of purchase is after 2:00pm and before 4:00pm

// **

func calculatePoints(receipt models.Receipt) int {
	points := 0

	// go through the reciept.Retailer string and ++ for each alphanumeric character
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			points++
		}
	}

	fmt.Println("points from retailer: ", points)

	points += calculateFiscalPoints(receipt.Total)
	points += calculateItemPoints(receipt.Items)
	points += calculateDateTimePoints(receipt.PurchaseDate, receipt.PurchaseTime)

	return points
}

func calculateDateTimePoints(purchaseDate string, purchaseTime string) int {
	points := 0

	// * 6 points if the day in the purchase date is odd.
	// * 10 points if the time of purchase is after 2:00pm and before 4:00pm

	// convert purchaseDate (YYYY-MM-DD) to day of month
	day, err := strconv.Atoi(purchaseDate[8:10])
	if err != nil {
		return 0
	}
	if day%2 != 0 {
		points += 6
	}

	// convert purchaseTime (HH:MM) to hour of day
	hour, err := strconv.Atoi(purchaseTime[:2])
	if err != nil {
		return 0
	}
	fmt.Println("hour: ", hour)
	if hour >= 14 && hour < 16 {
		points += 10
	}

	fmt.Println("points from date and time: ", points)

	return points
}

func calculateItemPoints(items []models.ReceiptItem) int {
	points := 0

	// we process the items in the reciept and calculate the points we award for them
	for _, item := range items {
		// convert price to float64
		priceFloat, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			continue
		}
		// calculate point if trimmed length = % 3
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			points += int(math.Ceil(priceFloat * 0.2))
		}
	}

	// 5 points for every two items on the receipt
	points += (len(items) / 2) * 5

	fmt.Println("points from items: ", points)

	return points
}

func calculateFiscalPoints(total string) int {
	points := 0

	// we only check the total for if it is a multiple of 0.25 or if it is a round dollar amount
	// convert total to float64
	totalFloat, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return 0
	}

	// according to readme it is inclusive so they both apply
	if totalFloat == float64(int(totalFloat)) {
		points += 50
	}

	if int(totalFloat*100)%25 == 0 {
		points += 25
	}

	fmt.Println("points from total: ", points)

	return points
}
