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

	user := app.Group("/user")

	su := service.NewUserService(
		service.WithPostgres(ctx, db, pool, log),
		service.WithRedis(rds),
	)

	cu := rest.NewUserRestController(su)

	user.Post("/register", cu.Register)
	user.Post("/login", cu.Login)
}
