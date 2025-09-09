package di

import (
	"github.com/ayyoob-k-a/finora/configs"
	"github.com/ayyoob-k-a/finora/db"
	"github.com/ayyoob-k-a/finora/handler"
	"github.com/ayyoob-k-a/finora/repo"
	"github.com/ayyoob-k-a/finora/routes"
	"github.com/ayyoob-k-a/finora/server"
	"github.com/ayyoob-k-a/finora/usecase"
)

func InitDI(cfg configs.Config) error {
	// Initialize the database connection
	db, err := db.InitDB(cfg)
	if err != nil {
		return err
	}

	ginServer := server.InitRouter()

	// Here you can set up your server with the database connection
	repoInstance := repo.NewRepo(db)
	usecase := usecase.NewUsecase(repoInstance)
	handler := handler.NewHandler(usecase)
	routes.AuthRoutes(ginServer, handler)
	server.StartServer(ginServer)

	return nil

}
