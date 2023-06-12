package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/qiangxue/fasthttp-routing"

	userDelivery "github.com/t1d333/vk_edu_db_project/internal/user/delivery/http"
	userRepository "github.com/t1d333/vk_edu_db_project/internal/user/repository/postgres"
	userService "github.com/t1d333/vk_edu_db_project/internal/user/service"

	forumDelivery "github.com/t1d333/vk_edu_db_project/internal/forum/delivery/http"
	forumRepository "github.com/t1d333/vk_edu_db_project/internal/forum/repository/postgres"
	forumService "github.com/t1d333/vk_edu_db_project/internal/forum/service"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return
	}

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	conn, err := pgx.Connect(context.Background(), "postgres://"+dbUser+":"+dbPassword+"@"+dbHost+":"+dbPort+"/"+dbName)
	if err != nil {
		logger.Error("Failed to connect to db ", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("Connection to db successfully")

	router := routing.New()

	userRep := userRepository.NewRepository(logger, conn)
	userServ := userService.NewService(logger, userRep)
	userDelivery.RegisterHandlers(router, logger, userServ)

	forumRep := forumRepository.NewRepository(logger, conn)
	forumServ := forumService.NewService(logger, forumRep)
	forumDelivery.RegisterHandlers(router, logger, forumServ)

	logger.Info("Server starting on port: 5000")
	fasthttp.ListenAndServe(":5000", router.HandleRequest)
}
