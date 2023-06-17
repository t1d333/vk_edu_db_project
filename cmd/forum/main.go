package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qiangxue/fasthttp-routing"

	userDelivery "github.com/t1d333/vk_edu_db_project/internal/user/delivery/http"
	userRepository "github.com/t1d333/vk_edu_db_project/internal/user/repository/postgres"
	userService "github.com/t1d333/vk_edu_db_project/internal/user/service"

	forumDelivery "github.com/t1d333/vk_edu_db_project/internal/forum/delivery/http"
	forumRepository "github.com/t1d333/vk_edu_db_project/internal/forum/repository/postgres"
	forumService "github.com/t1d333/vk_edu_db_project/internal/forum/service"

	threadDelivery "github.com/t1d333/vk_edu_db_project/internal/thread/delivery/http"
	threadRepository "github.com/t1d333/vk_edu_db_project/internal/thread/repository/postgres"
	threadService "github.com/t1d333/vk_edu_db_project/internal/thread/service"

	postDelivery "github.com/t1d333/vk_edu_db_project/internal/post/delivery/http"
	postRepository "github.com/t1d333/vk_edu_db_project/internal/post/repository/postgres"
	postService "github.com/t1d333/vk_edu_db_project/internal/post/service"

	serviceDelivery "github.com/t1d333/vk_edu_db_project/internal/service/delivery/http"
	serviceRepository "github.com/t1d333/vk_edu_db_project/internal/service/repository/postgres"
	serviceService "github.com/t1d333/vk_edu_db_project/internal/service/service"

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

	conf, _ := pgxpool.ParseConfig("postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?" + "pool_max_conns=100")
	conn, err := pgxpool.NewWithConfig(context.Background(), conf)
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

	threadRep := threadRepository.NewRepository(logger, conn)
	threadServ := threadService.NewService(logger, threadRep)
	threadDelivery.RegisterHandlers(router, logger, threadServ)

	postRep := postRepository.NewRepository(logger, conn)
	postServ := postService.NewService(logger, postRep)
	postDelivery.RegisterHandlers(router, logger, postServ)

	servRep := serviceRepository.NewRepository(logger, conn)
	servServ := serviceService.NewService(logger, servRep)
	serviceDelivery.RegisterHandlers(router, logger, servServ)

	logger.Info("Server starting on port: 5000")
	fasthttp.ListenAndServe(":5000", router.HandleRequest)
}
