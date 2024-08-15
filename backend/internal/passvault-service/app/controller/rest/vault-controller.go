package rest

import (
	"github.com/Novando/pintartek/internal/passvault-service/app/dto/vault"
	"github.com/Novando/pintartek/internal/passvault-service/app/service"
	"github.com/Novando/pintartek/pkg/auth"
	"github.com/Novando/pintartek/pkg/common/structs"
	"github.com/Novando/pintartek/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type VaultRestController struct {
	vaultServ *service.VaultService
}

// NewVaultRestController Initialize Vault controller using REST API
func NewVaultRestController(sv *service.VaultService) *VaultRestController {
	return &VaultRestController{vaultServ: sv}
}

// Create vault for storing credential
func (c *VaultRestController) Create(ctx *fiber.Ctx) error {
	var params vault.VaultRequest
	tokenStr := auth.GetTokenFromBearer(ctx.Get("Authorization"))
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
	res, code := c.vaultServ.Create(tokenStr, params)
	return ctx.Status(code).JSON(res)
}

// GetAll vault for current user
func (c *VaultRestController) GetAll(ctx *fiber.Ctx) error {
	tokenStr := auth.GetTokenFromBearer(ctx.Get("Authorization"))
	res, code := c.vaultServ.GetAll(tokenStr)
	return ctx.Status(code).JSON(res)
}
