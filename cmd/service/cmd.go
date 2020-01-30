package main

import (
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/z7zmey/golang-microservice-boilerplate/internal/endpoint"
	"github.com/z7zmey/golang-microservice-boilerplate/internal/repository"
	"github.com/z7zmey/golang-microservice-boilerplate/internal/service"
	"github.com/z7zmey/golang-microservice-boilerplate/internal/transport/rest/handler"
	"github.com/z7zmey/golang-microservice-boilerplate/internal/transport/rest/server/restapi"
	"github.com/z7zmey/golang-microservice-boilerplate/internal/transport/rest/server/restapi/operations"
)

var cfg Config

func init() {
	Cmd.Flags().AddFlagSet(cfg.Flags())
}

var Cmd = &cobra.Command{
	Use:   "serve",
	Short: "run server",
	RunE: func(cmd *cobra.Command, args []string) error {
		BindEnv(cmd)

		l := logrus.New()
		l.SetFormatter(new(logrus.JSONFormatter))
		l.Hooks.Add(newServiceLogHook("boilerplate"))

		// load embedded swagger file
		swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
		if err != nil {
			return errors.WithStack(err)
		}

		// create new service API
		api := operations.NewBoilerplateMicroserviceAPI(swaggerSpec)
		api.Logger = func(s string, i ...interface{}) {
			l.Log(logrus.InfoLevel, s)
		}
		server := restapi.NewServer(api)
		defer func() {
			if err := server.Shutdown(); err != nil {
				logrus.WithError(err).Error()
			}
		}()

		// set the port this service will be run on
		server.Port = cfg.Port
		server.Host = cfg.Host
		server.ReadTimeout = cfg.ReadTimeout
		server.WriteTimeout = cfg.WriteTimeout

		elasticsearchRepo := repository.NewCarRepo()
		searchService := service.NewSearchService(elasticsearchRepo)
		searchEndpoint := endpoint.NewBoilerplateEndpoint(searchService)
		apiHandler := handler.NewApiHandler(searchEndpoint)

		apiHandler.ConfigureHandlers(api)

		server.SetHandler(PanicRecovery(l, AccessLogger(l, api.Serve(middleware.PassthroughBuilder))))

		// serve API
		if err := server.Serve(); err != nil {
			return errors.WithStack(err)
		}
		return nil
	},
}

// BindEnv parse config values from environment variables
func BindEnv(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		envVar := strings.ToUpper(f.Name)

		if val := os.Getenv(envVar); val != "" {
			if err := cmd.Flags().Set(f.Name, val); err != nil {
				logrus.WithError(err).Error("failed to set flag")
			}
		}
	})
}

func PanicRecovery(l *logrus.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rerr := recover(); rerr != nil {
				defer func(rw http.ResponseWriter) {
					rw.WriteHeader(http.StatusInternalServerError)
				}(w)

				httprequest, err := httputil.DumpRequest(r, false)
				if err != nil {
					l.WithError(err).Errorf("failed to dump the request on panic %v", rerr)
					return
				}

				l.WithError(errors.New("panic recovery")).
					WithField("stack", string(debug.Stack())).
					WithField("request", string(httprequest)).
					Error("panic")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func AccessLogger(l *logrus.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		// todo add request status code
		l.WithFields(logrus.Fields{
			"method":       r.Method,
			"path":         r.URL,
			"took-seconds": time.Since(start).Seconds(),
		}).Info("REST request")
	})
}
