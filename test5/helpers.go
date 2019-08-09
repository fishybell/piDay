package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/drivers/gpio"
)

// AllowRequestsFrom returns handler to verify that requests come from allowedOrigins
// stolen, and fixed from gobot
func AllowRequestsFrom(allowedOrigins ...string) http.HandlerFunc {
	c := &api.CORS{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "POST", "OPTIONS"}, // added OPTIONS, doesn't seem to matter
		AllowHeaders: []string{"Origin", "Content-Type"},
		ContentType:  "application/json; charset=utf-8",
	}

	return func(w http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin) // lol, because we got rid of the check, this lets any IP
		w.Header().Set("Access-Control-Allow-Headers", c.AllowedHeaders())
		w.Header().Set("Access-Control-Allow-Methods", c.AllowedMethods())
		w.Header().Set("Content-Type", c.ContentType)
	}
}

// for some reason we sometimes end up with the params as strings
func parseInt(something map[string]interface{}, param string) int64 {
	inter, ok := something[param]
	if !ok {
		fmt.Println("got no param in params")
		return 0
	}
	switch val := inter.(type) {
	case string:
		retval, _ := strconv.ParseInt(val, 10, 32) // if this errors, we'll just have no speed, or no slider, which is fine
		return retval
	case float64:
		return int64(val)
	case int:
		return int64(val)
	case int64:
		return val
	default:
		fmt.Printf("got something odd %T:%+v\n", inter, inter)
		return 0
	}
}

func watchForShutdown(pwms []*pwmDriver, hbridges ...*gpio.LedDriver) {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel
	fmt.Println("Finishing...")
	for _, hbridge := range hbridges {
		hbridge.Off()
	}
	for _, pwm := range pwms {
		pwm.Stop()
	}
	fmt.Println("Finished")
	os.Exit(0)
}
