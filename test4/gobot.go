package main

import (
	"fmt"
	//"html"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	"net/http"
)

func main() {
	master := gobot.NewMaster()

	// start the api server
	a := api.NewAPI(master)
	//a.AddHandler(func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "Hello, %q \n", html.EscapeString(r.URL.Path))
	//})
	//a.Debug()
	//a.Get("/static", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "Hello, %q \n", html.EscapeString(r.URL.Path))
	//})
	//a.Post("/static", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "Hello, %q \n", html.EscapeString(r.URL.Path))
	//})
	a.AddHandler(api.AllowRequestsFrom("http://localhost:3030"))
	a.AddC3PIORoutes()
	//a.Start()
	a.StartWithoutDefaults()

	// define our commands
	master.AddCommand("custom_gobot_command",
		func(params map[string]interface{}) interface{} {
			return "This command is attached to the mcp!"
		})

	hello := master.AddRobot(gobot.NewRobot("hello"))

	hello.AddCommand("hi_there", func(params map[string]interface{}) interface{} {
		return fmt.Sprintf("This command is attached to the robot %v", hello.Name)
	})

	// server our static directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("listening on 3030 as well?")
	http.ListenAndServe(":3030", nil)

	// waits forever
	master.Start()
}
