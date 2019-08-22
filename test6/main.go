package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	//speeds1 := speedS()
	//speeds2 := speedMap3()
	//blended := blendMaps(speeds1, speeds2, 0.4, 1.6)
	//speeds2AndSome := speedMap3AndSome()
	//blendedAndSome := blendMaps(speeds1, speeds2AndSome, 0.4, 1.6)
	////tightened := tighten(blended, 5, 5)
	//for i := -100; i < 101; i++ {
	//	log.Printf("%d : %d v %d v %d v %d", i, speeds1[i], speeds2[i], blended[i])
	//}
	////return

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

	// right front
	pwm4, hbridge4_a, hbridge4_b := motorIO(r,
		"22", // physical pin 22 => gpio 25 (pwm pin at the left)
		"18", // physical pin 18 => gpio 24
		"16", // physical pin 16 => gpio 23
	)

	// right back
	pwm3, hbridge3_a, hbridge3_b := motorIO(r,
		"11", // physical pin 11 => gpio 17 (pwm pin at the right)
		"13", // physical pin 13 => gpio 27
		"15", // physical pin 15 => gpio 22
	)

	// define our motors

	motors := []Motor{ // the order matters, because it does
		NewFakeMotor(), // you hack!
		NewMappedMotor(pwm1, speedMap3(), hbridge1_a, hbridge1_b, "front left"),
		NewMappedMotor(pwm4, speedMap3AndSome(), hbridge4_a, hbridge4_b, "front right"), // some idiot poured hot glue on this motor, give it a fighting chance
		NewMappedMotor(pwm2, speedMap3(), hbridge2_a, hbridge2_b, "back left"),
		NewMappedMotor(pwm3, speedMap3(), hbridge3_a, hbridge3_b, "back right"),
	}

	//motors := []Motor{ // the order matters, because it does
	//	NewFakeMotor(), // you hack!
	//	NewMappedMotor(pwm1, blended, hbridge1_a, hbridge1_b, "front left"),
	//	NewMappedMotor(pwm4, blendedAndSome, hbridge4_a, hbridge4_b, "front right"), // some idiot poured hot glue on this motor, give it a fighting chance
	//	NewMappedMotor(pwm2, blended, hbridge2_a, hbridge2_b, "back left"),
	//	NewMappedMotor(pwm3, blended, hbridge3_a, hbridge3_b, "back right"),
	//}

	// define our commands
	car := master.AddRobot(gobot.NewRobot("car"))

	car.AddCommand("speed", func(params map[string]interface{}) interface{} {
		speed := parseInt(params, "speed")
		slider := parseInt(params, "slider")
		handleSpeed(motors, speed, slider)
		return toggleStatus(motors)
	})

	car.AddCommand("status", func(params map[string]interface{}) interface{} {
		return toggleStatus(motors)
	})

	car.AddCommand("fail", func(params map[string]interface{}) interface{} {
		fail()
		return "on"
	})

	car.AddCommand("succeed", func(params map[string]interface{}) interface{} {
		succeed()
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
	pwm = SoftPwmInit(led, 25*time.Nanosecond, 100)
	pwm.SetDutyCycle(0)
	pwm.Start() // if we don't start, we hang on a channel forever

	return
}

// map of x^3 for a curvier feel
func speedMap3() map[int]int {
	myMap := make(map[int]int, 201) // -100 to 100
	for i := 0; i < 201; i++ {
		key := i - 100
		myMap[key] = (key * key * key) / 10000
	}

	return myMap
}

// map of x^3 for a curvier feel, but with a little extra under the hood
func speedMap3AndSome() map[int]int {
	myMap := make(map[int]int, 201) // -100 to 100
	for i := 0; i < 201; i++ {
		key := i - 100
		myMap[key] = (key * key * key) / 9500
	}

	return myMap
}

// map of a sigmoid "s-curve"
func speedS() map[int]int {
	myMap := make(map[int]int, 201) // -100 to 100
	for i := 0; i < 201; i++ {
		key := float64(i - 100)
		partA := 1.0 / (1.0 + math.Pow(math.E, (-1.0*(key/20.0))))
		partB := (math.Pow(math.E, (key / 20.0))) / (1.0 + math.Pow(math.E, (key/20.0)))
		partC := ((partA + partB) * 100.0) - 100.0
		myMap[int(key)] = int(partC)
	}

	return myMap
}

func blendMaps(map1, map2 map[int]int, weight1, weight2 float64) map[int]int {
	myMap := make(map[int]int, 201) // -100 to 100
	for i := 0; i < 201; i++ {
		myMap[i] = int(((float64(map1[i]) * weight1) + (float64(map2[i]) * weight2)) / 2.0)
	}

	return myMap
}

func toggleStatus(motors []Motor) string {
	toggles := struct {
		FL bool `json:"fl"`
		FR bool `json:"fr"`
		BL bool `json:"bl"`
		BR bool `json:"br"`
	}{
		motors[1].Enabled(),
		motors[2].Enabled(),
		motors[3].Enabled(),
		motors[4].Enabled(),
	}
	bytes, _ := json.Marshal(toggles)
	return string(bytes)
}

//func tighten(speedMap map[int]int, dropFirst, dropLast int) map[int]int {
//	myMap := make(map[int]int, 201) // -100 to 100
//	for i := 0; i < 201; i++ {
//		key := float64(i - 100)
//		if key < 0.0 && key > -50 {
//			// -50 to -1
//		} else if key < 0.0 {
//			// -100 to -51
//		} else if key > 50.0 {
//			// 51 to 100
//		} else if key > 0.0 {
//			// 1 to 51
//		} // f you zero, you stay where you are
//	}
//
//	return myMap
//}
