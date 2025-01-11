package controllers

import (
	"fmt"
	"src/internals/database"
	"src/internals/utils"

	"github.com/gofiber/fiber/v2"
)

type (
	UserBody struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
)

func Signup(c *fiber.Ctx) error {
	var userBody UserBody
	err := c.BodyParser(&userBody)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "1",
			"message": "Invalid input",
		})
	}

	validator := utils.NewValidator()
	errors := validator.Validate(userBody)

	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"message": utils.FormatValidationErrors(errors),
		})
	}

	hashPassword, err := utils.HashPassword(userBody.Password)
	if err != nil {
		fmt.Println("Cannot hashpassword", err)
	}

	_, dberr := database.CreateUser(userBody.Name, userBody.Email, hashPassword)
	if dberr != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "1",
			"message": "Cannot create user: " + dberr.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   "0",
		"message": "User created",
	})
}
