package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/junaidshaikh-js/workout-api/internal/api"
	"github.com/junaidshaikh-js/workout-api/internal/store"
	"github.com/junaidshaikh-js/workout-api/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()

	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")

	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)

	userStore := store.NewPostgresUserStore(pgDB)
	userHandler := api.NewUserHandler(userStore, logger)

	tokenStore := store.NewPostgresTokenStore(pgDB)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		DB:             pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Health: OK")
}
