package passvaultService

import (
	"context"
	"github.com/Novando/pintartek/internal/passvault-service/app/controller/rest"
	"github.com/Novando/pintartek/internal/passvault-service/app/service"
	"github.com/Novando/pintartek/pkg/logger"
	"github.com/Novando/pintartek/pkg/postgresql/pgx"
	"github.com/Novando/pintartek/pkg/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPassvaultService(
	app fiber.Router,
	db *pgx.Queries,
	pool *pgxpool.Pool,
	rds *redis.Redis,
	log *logger.Logger,
) {
	ctx := context.Background()

	su := service.NewUserService(
		service.WithUserPostgres(ctx, db, pool, log),
		service.WithUserRedis(rds),
	)
	sv := service.NewVaultService(
		service.WithVaultPostgres(ctx, db, pool, log),
		service.WithVaultRedis(rds),
	)

	cu := rest.NewUserRestController(su)
	cv := rest.NewVaultRestController(sv)

	user := app.Group("/user")
	user.Get("/logout", cu.Logout)
	user.Post("/register", cu.Register)
	user.Post("/login", cu.Login)

	vault := app.Group("/vault")
	vault.Get("/", cv.GetAll)
	vault.Get("/:vaultId", cv.GetOne)
	vault.Post("/", cv.Create)
	vault.Post("/:vaultId", cv.CreateCredential)
	vault.Put("/:vaultId", cv.UpdateVaultName)
	vault.Put("/:vaultId/:credentialId", cv.UpdateCredential)
	//vault.Delete("/:vaultId", cv.Create)
	//vault.Delete("/:vaultId/:credentialId", cv.Create)
}
