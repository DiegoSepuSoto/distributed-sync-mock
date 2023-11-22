package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

const closeAppTimeout = time.Second * 10
const rangeInterval = 10

var rangesForIDs = []string{"10-20", "21-30", "31-40", "41-50"}
var rangesInUse []string

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/range", assignIDsRange)

	quit := make(chan os.Signal, 1)
	go startServer(e, quit)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	gracefulShutdown(e)
}

func assignIDsRange(c echo.Context) error {
	rangeToUse := ""

	for _, rangeForIDs := range rangesForIDs {
		if !existsInSlice(rangeForIDs, rangesInUse) {
			rangeToUse = rangeForIDs
			rangesInUse = append(rangesInUse, rangeForIDs)

			break
		}
	}

	if rangeToUse == "" {
		rangeToUse = addRange()
	}

	log.Println("sending range: ", rangeToUse)

	return c.JSON(http.StatusOK, echo.Map{"rangeIDs": rangeToUse})
}

func existsInSlice(element string, slice []string) bool {
	for _, e := range slice {
		if element == e {
			return true
		}
	}

	return false
}

func addRange() string {
	lastRange := rangesForIDs[len(rangesForIDs)-1]

	re := regexp.MustCompile(`(\d+)-(\d+)`)

	matches := re.FindStringSubmatch(lastRange)

	topRangeValue, _ := strconv.Atoi(matches[2])

	nextTopRangeValue := topRangeValue + rangeInterval

	newRange := fmt.Sprintf("%d-%d", topRangeValue+1, nextTopRangeValue)

	rangesForIDs = append(rangesForIDs, newRange)
	rangesInUse = append(rangesInUse, newRange)

	return newRange
}

func startServer(e *echo.Echo, quit chan os.Signal) {
	log.Print("starting server")

	if err := e.Start(":" + os.Getenv("APP_PORT")); err != nil {
		log.Error(err.Error())
		close(quit)
	}
}

func gracefulShutdown(e *echo.Echo) {
	log.Print("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), closeAppTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
