package main

import (
	"fmt"
	"os"
	"sync"
	"vas/bootstrap"
	"vas/config"
	"vas/controllers"
	"vas/logger"
	"vas/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	exitCtx, shutDown, err := bootstrap.Startup()
	if err != nil {
		fmt.Println("Bootup error: " + err.Error())
		os.Exit(1)
	}

	curEnv := os.Getenv("GO_ENV")
	cfg := config.GetConfig()

	var nrApp *newrelic.Application

	if curEnv == "prod" {
		nrApp, err = newrelic.NewApplication(
			newrelic.ConfigAppName("template"),
			newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		)
		if err != nil {
			fmt.Println("newRelic error: " + err.Error())
			os.Exit(1)
		}
	}

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "template",
		AppName:       "template v0.0.0",
	})

	app.Use(recover.New())

	app.Use(helmet.New())

	app.Use(cors.New())

	app.Use(requestid.New(requestid.Config{
		ContextKey: "requestId",
	}))

	app.Use(compress.New())

	// Request logger middlewares
	app.Use(middlewares.LogRequestMiddleware())

	if curEnv == "prod" {
		// Newrelic metric middleware
		app.Use(middlewares.New(middlewares.Config{
			NewRelicApp: nrApp,
		}))
	}

	vas := app.Group("/vas")
	api := vas.Group("/api")
	// api version group
	v1 := api.Group("/v1")

	// Health check
	v1.Get("/health", func(c *fiber.Ctx) error {
		res := map[string]interface{}{
			"data":   "Server is up and running",
			"errors": nil,
			"status": 200,
		}
		if err := c.JSON(res); err != nil {
			return err
		}
		return nil
	})

	// ================= API Routes ================= //

	v1.Post("/sample", controllers.ValiadateSamplePayload, controllers.Sample)

	// ============================================== //

	// ============== Async Processes =============== //

	shutDownWg := &sync.WaitGroup{}
	shutDownWg.Add(1)
	go bootstrap.GracefulShutDown(*exitCtx, app, shutDownWg)
	// => Add your async processes here and increase number in shutDown waitGroup

	// ============================================== //

	err = app.Listen(":" + cfg.APP.PORT)
	if err != nil {
		logger.LogPanic(nil, "Fiber server error: "+err.Error(), nil)
		fmt.Println("Fiber server error: ", err)
		os.Exit(1)
	}

	// cleanUp
	bootstrap.CleanUp(*shutDown)
	// waiting for shutting down
	shutDownWg.Wait()
}
