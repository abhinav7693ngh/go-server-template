package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"vas/config"
	"vas/logger"

	"github.com/gofiber/fiber/v2"
)

func GracefulShutDown(exitCtx context.Context, app *fiber.App, shutDownWg *sync.WaitGroup) {
	defer shutDownWg.Done()
	_ = <-exitCtx.Done()
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()
}

func CleanUp(shutdown context.CancelFunc) {
	// Do db, redis etc disconnection here
	shutdown()
}

func Startup() (*context.Context, *context.CancelFunc, error) {
	errLogger := logger.LoggerInit()
	if errLogger != nil {
		return nil, nil, errors.New("Error initializing logger: " + errLogger.Error())
	}

	errLoadConfig := config.ConfigInit()
	if errLoadConfig != nil {
		return nil, nil, errors.New("Error loading config " + errLoadConfig.Error())
	}

	// errDb := database.Connect()
	// if errDb != nil {
	// 	return nil, nil, errors.New("Database connection error: " + errDb.Error())
	// }

	// Exit context for graceful shutdown
	exitCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	return &exitCtx, &cancel, nil
}
