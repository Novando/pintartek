package passvaultService

import (
	"github.com/Novando/pintartek/internal/passvault-service/app/controller/rest"
	"github.com/Novando/pintartek/internal/passvault-service/app/service"
	"github.com/Novando/pintartek/pkg/logger"
	"github.com/Novando/pintartek/pkg/postgresql/pgx"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPassvaultService(
	app fiber.Router,
	db *pgx.Queries,
	pool *pgxpool.Pool,
	log *logger.Logger,
) {
	user := app.Group("/user")

	su := service.NewUserService(service.WithPostgres(log, db, pool))

	cu := rest.NewUserRestController(su)

	user.Post("/register", cu.Register)
}
