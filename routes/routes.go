package routes

import (
	"adiachenko/go-scaffold/internal/platform/jaeger"
	"adiachenko/go-scaffold/routes/handlers"

	"github.com/Vinelab/go-reporting"
	"github.com/Vinelab/go-reporting/sentry"
	"github.com/go-chi/chi"

	tracingMiddleware "github.com/Vinelab/tracing-go/middleware"
	chiMiddleware "github.com/go-chi/chi/middleware"
)

// Register holds the routes to be registered
// when the server starts listening
func Register() *chi.Mux {
	router := chi.NewRouter()

	router.Use(tracingMiddleware.NewTraceRequests(jaeger.Trace, []string{"application/json"}, []string{"/"}).Handler)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.StripSlashes)
	router.Use(sentry.LogResponseMiddleware)
	router.Use(chiMiddleware.Recoverer)
	router.Use(reporting.LogPanicMiddleware)

	router.Get("/", handlers.Welcome)

	return router
}
