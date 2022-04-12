package main

import (
	"context"
	"strconv"

	"shortener/config"
	"shortener/internal/controllers"
	"shortener/internal/repositories"
	"shortener/internal/usecases"

	"github.com/fasthttp/router"
	"github.com/mehdiazizii/fastcontroller"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func init() {
	migrateCli.PersistentFlags().UintP("version", "v", 0, "migrate to version")
	dbCli.AddCommand(migrateCli)
	rootCli.AddCommand(serveCli, dbCli)
}

func main() {
	if err := rootCli.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func route(l *logrus.Logger, cfg fastcontroller.Config, db *repositories.DbSession) *router.Router {
	controller := fastcontroller.NewController(l, cfg)

	linkRepo := repositories.NewLink(db)
	linkService := usecases.NewLink(linkRepo)
	linkController := controllers.NewLink(controller, linkService)

	r := router.New()
	r.POST("/api/link", controller.Handle(linkController.Post))
	r.GET("/{link}", controller.Handle(linkController.Get))

	return r
}

func serve(ctx context.Context) {
	log := logrus.New()
	cfg := config.Config()
	db, err := repositories.NewSession(cfg.DbSession)
	if err != nil {
		logrus.Panic(err)
	}

	r := route(log, cfg, db)
	server := fasthttp.Server{Handler: r.Handler}

	go func() {
		<-ctx.Done()
		_ = server.Shutdown()
	}()

	if err := server.ListenAndServe(":" + strconv.Itoa(cfg.HTTPPort)); err != nil {
		logrus.Fatal(err)
	}
}
