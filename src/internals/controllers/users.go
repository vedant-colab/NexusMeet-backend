package controllers

import (
	"crypto/rand"
	"fmt"
	"log"
	"src/internals/database"
	"src/internals/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type (
	UserBody struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	SignIn struct {
		Name     string `json:"name" validate:"required"`
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

func Signin(c *fiber.Ctx) error {
	var signin SignIn
	if err := c.BodyParser(&signin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "1",
			"message": "Invalid input",
		})
	}

	validator := utils.NewValidator()
	errors := validator.Validate(signin)

	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"message": utils.FormatValidationErrors(errors),
		})
	}

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}
	// keyString := hex.EncodeToString(key)
	encryptedUser, err := utils.Encrypt(signin.Name, key)
	if err != nil {
		return fmt.Errorf("error in encrypting %v", err)
	}
	fmt.Println(encryptedUser)

	tokenString, err := utils.CreateToken(encryptedUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error_code": "1",
			"message":    fmt.Errorf("error creating token: %v", err),
		})
	}

	_, dberr := database.SaveUserToken(tokenString, signin.Name, signin.Password)
	if dberr != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "1",
			"message": "Problem signing user :" + dberr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error_code": "0",
		"message":    "Token created sucessfully",
		"token":      tokenString,
	})
}

func Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(400).JSON(fiber.Map{
			"error_code": "1",
			"message":    "Token missing",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Authorization header format")
	}

	token := parts[1]
	claims, err := utils.GetPayloadFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error_code": "1",
			"message":    "Invalid token",
		})
	}
	encryptedUsername := claims["username"].(string)
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}
	decryptedUsername, err := utils.Decrypt(encryptedUsername, string(key))
	if err != nil {
		return fmt.Errorf("error in decrypting %v", err)
	}
	err = database.DeactivateToken(token, decryptedUsername)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error_code": "1",
			"message":    "Problem logging out user",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"error_code": "0",
		"message":    "User logged out successfully",
	})

}

func Load(c *fiber.Ctx) error {
	sum := 0
	for i := 0; i < 1e5; i++ {
		sum += i
	}
	return c.Status(200).JSON(fiber.Map{
		"Sum": sum,
	})
}
