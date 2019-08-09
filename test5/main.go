package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	master := gobot.NewMaster()

	// start the api server
	a := api.NewAPI(master)
	a.AddHandler(AllowRequestsFrom())
	//a.AddC3PIORoutes()
	//a.StartWithoutDefaults()
	a.Start()

	// define our I/O
	r := raspi.NewAdaptor()
	//led := gpio.NewLedDriver(r, "7") // physical pin 7 => gpio 4
	led := gpio.NewLedDriver(r, "40")        // physical pin 40 => gpio 21
	hbridge1_a := gpio.NewLedDriver(r, "36") // physical pin 36 => gpio 16
	hbridge1_b := gpio.NewLedDriver(r, "38") // physical pin 38 => gpio 20
	pwm := SoftPwmInit(led, 10*time.Nanosecond, 100)

	// define our motors

	motors := []Motor{
		NewFakeMotor(), // you hack!
		NewMotor(pwm, hbridge1_a, hbridge1_b),
	}

	// define our commands
	car := master.AddRobot(gobot.NewRobot("car"))

	car.AddCommand("speed", func(params map[string]interface{}) interface{} {
		speed := parseInt(params, "speed")
		slider := parseInt(params, "slider")
		fmt.Println("got speed command", speed, slider)

		if speed < 0 {
			motors[slider].Back()
			motors[slider].Speed(-1 * speed)
		} else {
			motors[slider].Fwd()
			motors[slider].Speed(speed)
		}
		return handleSpeed(speed, slider)
	})

	// server our static directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Serving HTML on :3030...")
	http.ListenAndServe(":3030", nil)

	// waits forever
	master.Start()

	// cleanup
	watchForShutdown(pwm,
		hbridge1_a,
		hbridge1_b,
	)
}
