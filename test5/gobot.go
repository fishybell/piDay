package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
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

func main() {
	master := gobot.NewMaster()

	// start the api server
	a := api.NewAPI(master)
	a.AddHandler(AllowRequestsFrom())
	//a.AddC3PIORoutes()
	//a.StartWithoutDefaults()
	a.Start()

	// define our commands
	master.AddCommand("custom_gobot_command",
		func(params map[string]interface{}) interface{} {
			return "This command is attached to the mcp!"
		})

	hello := master.AddRobot(gobot.NewRobot("hello"))

	hello.AddCommand("hi_there", func(params map[string]interface{}) interface{} {
		return fmt.Sprintf("This command is attached to the robot %v", hello.Name)
	})

	hello.AddCommand("speed", func(params map[string]interface{}) interface{} {
		speed := parseInt(params, "speed")
		slider := parseInt(params, "slider")
		return fmt.Sprintf("received speed %d for slider %d from params %+v", speed, slider, params)
	})

	// server our static directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Serving HTML on :3030...")
	http.ListenAndServe(":3030", nil)

	// waits forever
	master.Start()
}

// for some reason we always end up with the params as strings
func parseInt(something map[string]interface{}, param string) int64 {
	inter, ok := something[param]
	if !ok {
		return 0
	}
	valString, ok := inter.(string)
	if !ok {
		return 0
	}

	retval, _ := strconv.ParseInt(valString, 10, 32) // if this errors, we'll just have no speed, or no slider, which is fine

	return retval
}
