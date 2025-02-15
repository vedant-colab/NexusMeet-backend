package controllers

import (
	"fmt"
	"net/url"
	"src/internals/database"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UpdateProfile(c *fiber.Ctx) error {
	// Extract the username from the URL parameters
	username := c.Params("username")
	decodedUsername, err := url.QueryUnescape(username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid username format",
		})
	}

	// Trim any whitespace
	decodedUsername = strings.TrimSpace(decodedUsername)

	// Validate username is not empty after processing
	if decodedUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username cannot be empty",
		})
	}

	// Parse the request body
	var updateData struct {
		Bio            string `json:"bio"`
		Location       string `json:"location"`
		Website        string `json:"website"`
		ProfilePicture string `json:"profile_picture"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update the profile in the database
	err = database.UpdateProfileInDB(decodedUsername, updateData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}

func FetchProfile(c *fiber.Ctx) error {
	username := c.Params("username")
	decodedUsername, err := url.QueryUnescape(username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid username format",
		})
	}

	// Trim any whitespace
	decodedUsername = strings.TrimSpace(decodedUsername)

	// Validate username is not empty after processing
	if decodedUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username cannot be empty",
		})
	}

	fmt.Println("Fetching profile for:", decodedUsername)
	user, err := database.GetUserByUsername(decodedUsername) // Use decodedUsername here
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": nil, "user": user})
}
