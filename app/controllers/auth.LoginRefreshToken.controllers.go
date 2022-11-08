package controllers

import (
	"fmt"
	"net/http"
	"os"
	"sejuta-cita/app/models"
	"sejuta-cita/app/requests"
	"sejuta-cita/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt"
	"github.com/gookit/validate"
	"golang.org/x/crypto/bcrypt"
)

var refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))

func LoginRefrehToken(c *fiber.Ctx) error {
	var userClaim models.UserClaim

	p := new(requests.LoginForm)
	if err := c.BodyParser(p); err != nil {
		return err
	}
	v := validate.Struct(p)
	if !v.Validate() {
		return c.JSON(fiber.Map{
			"Message": v.Errors.One(),
		})
	}
	user := new(models.User)

	db := config.GetDBInstance()

	if res := db.Where("Email = ?", p.Email).Preload("Role").First(&user); res.RowsAffected <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error":   true,
			"Message": "Invalid Email!",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error":   true,
			"Message": "Password is incorrect!",
		})
	}
	userClaim.Issuer = utils.UUIDv4()
	userClaim.Id = int(user.ID)
	userClaim.Email = user.Email
	userClaim.Role = user.Role.Name
	accessToken, refreshToken := models.GenerateTokens(&userClaim, false)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"AccessToken":  accessToken,
		"RefreshToken": refreshToken,
		"User": fiber.Map{
			"ID":    user.ID,
			"Email": user.Email,
			"Role":  user.Role.Name,
		},
	})
}

func RefreshToken(c *fiber.Ctx) error {
	var userClaim models.UserClaim
	// var response models.Response

	refreshToken := c.FormValue("refresh_token")

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return refreshSecret, nil
	})

	if err != nil {
		fmt.Println("the error from parse: ", err)
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"Message": "Invalid token"})
	}

	//is token valid?
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"Message": "StatusUnauthorized",
		})
	}

	user_id := claims["user_id"]
	user_id_int := int(user_id.(float64))
	userClaim.Issuer = fmt.Sprintf("%v", claims["issuer"])
	userClaim.Id = user_id_int
	userClaim.Email = fmt.Sprintf("%v", claims["email"])
	userClaim.Role = fmt.Sprintf("%v", claims["role"])

	// if fail refresh token
	accessToken, refreshToken := models.GenerateTokens(&userClaim, true)
	if len(accessToken) < 1 || len(refreshToken) < 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Message": "something went wrong",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"AccessToken":  accessToken,
		"RefreshToken": refreshToken,
		"User": fiber.Map{
			"ID":    userClaim.Id,
			"Email": userClaim.Email,
			"Role":  userClaim.Role,
		},
	})
}
