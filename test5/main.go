package main

import (
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

	// define our I/O (laid out here in order of pins)
	r := raspi.NewAdaptor()

	// left front
	pwm1, hbridge1_a, hbridge1_b := motorIO(r,
		"40", // physical pin 40 => gpio 21 (pwm pin at the left)
		"36", // physical pin 36 => gpio 16
		"38", // physical pin 38 => gpio 20
	)

	// left back
	pwm2, hbridge2_a, hbridge2_b := motorIO(r,
		"24", // physical pin 24 => gpio 8 (pwm pin at the right)
		"32", // physical pin 32 => gpio 12
		"26", // physical pin 26 => gpio 7
	)

	// right back
	pwm3, hbridge3_a, hbridge3_b := motorIO(r,
		"22", // physical pin 22 => gpio 25 (pwm pin at the left)
		"18", // physical pin 18 => gpio 24
		"16", // physical pin 16 => gpio 23
	)

	// right front
	pwm4, hbridge4_a, hbridge4_b := motorIO(r,
		"11", // physical pin 11 => gpio 17 (pwm pin at the right)
		"13", // physical pin 13 => gpio 27
		"15", // physical pin 15 => gpio 22
	)

	// define our motors

	motors := []Motor{ // the order matters, because it does
		NewFakeMotor(), // you hack!
		NewMotor(pwm1, hbridge1_a, hbridge1_b, "front left"),
		NewMotor(pwm4, hbridge4_a, hbridge4_b, "front right"),
		NewMotor(pwm2, hbridge2_a, hbridge2_b, "back left"),
		NewMotor(pwm3, hbridge3_a, hbridge3_b, "back right"),
	}

	// define our commands
	car := master.AddRobot(gobot.NewRobot("car"))

	car.AddCommand("speed", func(params map[string]interface{}) interface{} {
		speed := parseInt(params, "speed")
		slider := parseInt(params, "slider")
		handleSpeed(motors, speed, slider)
		return "stuff"
	})

	car.AddCommand("status", func(params map[string]interface{}) interface{} {
		return "on"
	})

	// server our static directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Serving HTML on :3030...")
	http.ListenAndServe(":3030", nil)

	// waits forever
	master.Start()

	// cleanup
	watchForShutdown([]*pwmDriver{
		pwm1,
		pwm2,
		pwm3,
		pwm4,
	},
		hbridge1_a,
		hbridge1_b,
		hbridge2_a,
		hbridge2_b,
		hbridge3_a,
		hbridge3_b,
		hbridge4_a,
		hbridge4_b,
	)
}

func motorIO(r *raspi.Adaptor, pin1 string, pin2 string, pin3 string) (pwm *pwmDriver, hbridge1 *gpio.LedDriver, hbridge2 *gpio.LedDriver) {
	led := gpio.NewLedDriver(r, pin1)
	hbridge1 = gpio.NewLedDriver(r, pin2)
	hbridge2 = gpio.NewLedDriver(r, pin3)
	pwm = SoftPwmInit(led, 10*time.Nanosecond, 100)
	pwm.SetDutyCycle(0)
	pwm.Start() // if we don't start, we hang on a channel forever

	return
}
