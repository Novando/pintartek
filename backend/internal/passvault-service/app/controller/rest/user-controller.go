package rest

import (
	"github.com/Novando/pintartek/internal/passvault-service/app/dto/user"
	"github.com/Novando/pintartek/internal/passvault-service/app/service"
	"github.com/Novando/pintartek/pkg/auth"
	"github.com/Novando/pintartek/pkg/common/structs"
	"github.com/Novando/pintartek/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type UserRestController struct {
	userServ *service.UserService
}

// NewUserRestController Initialize User controller using REST API
func NewUserRestController(su *service.UserService) *UserRestController {
	return &UserRestController{userServ: su}
}

// Register the entry point for user creation
func (c *UserRestController) Register(ctx *fiber.Ctx) error {
	var params user.RegisterRequest
	if err := ctx.BodyParser(&params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(structs.StdResponse{
			Message: "PAYLOAD_ERROR",
			Data:    err.Error(),
		})
	}
	if err := validator.Validate(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(structs.StdResponse{
			Message: "VALIDATION_ERROR",
			Data:    err.Error(),
		})
	}
	res, code := c.userServ.Register(params)
	return ctx.Status(code).JSON(res)
}

// Login the entry point for session creation
func (c *UserRestController) Login(ctx *fiber.Ctx) error {
	var params user.LoginRequest
	if err := ctx.BodyParser(&params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(structs.StdResponse{
			Message: "PAYLOAD_ERROR",
			Data:    err.Error(),
		})
	}
	if err := validator.Validate(params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(structs.StdResponse{
			Message: "VALIDATION_ERROR",
			Data:    err.Error(),
		})
	}
	res, code := c.userServ.Login(params)
	return ctx.Status(code).JSON(res)
}

// Logout delete the session for current user
func (c *UserRestController) Logout(ctx *fiber.Ctx) error {
	tokenStr := auth.GetTokenFromBearer(ctx.Get("Authorization"))
	res, code := c.userServ.Logout(tokenStr)
	return ctx.Status(code).JSON(res)
}
