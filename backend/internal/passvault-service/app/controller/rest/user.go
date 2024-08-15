package rest

import (
	"github.com/Novando/pintartek/internal/passvault-service/app/dto"
	"github.com/Novando/pintartek/internal/passvault-service/app/service"
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
	var params dto.RegisterRequest
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
