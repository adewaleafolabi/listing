package main

import (
	"fmt"
	"github.com/adewaleafolabi/listing/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/tinrab/retry"
	"log"
	"net/http"
	"time"
)

type Config struct {
	DBName     string `envconfig:"POSTGRES_DB"`
	DBUser     string `envconfig:"POSTGRES_USER"`
	DBPassword string `envconfig:"POSTGRES_PASSWORD"`
	DBHost     string `envconfig:"POSTGRES_HOST"`
	DBPort     string `envconfig:"POSTGRES_PORT"`
}

func setUpRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/api/listing/properties", PropertyRoutes())
	})

	return r
}
func main() {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	retry.ForeverSleep(2*time.Second, func(i int) error {
		url := fmt.Sprintf(`postgres://%s:%s@%s:%s/%s?sslmode=disable`, cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
		if repo, err := db.NewPostgres(url); err != nil {
			logrus.Error(err)
			return err
		} else {
			db.SetRepository(repo)
			return nil
		}
	})

	router := setUpRouter()

	walkFn := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logrus.Printf("%s %s\n", method, route)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(router, walkFn); err != nil {
		logrus.Panic(err)
	}

	logrus.Fatal(http.ListenAndServe(":8000", router))
}
