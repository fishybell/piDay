package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

type pwmDriver struct {
	led    *gpio.LedDriver
	talker chan int
	tick   time.Duration
	cycle  int
	ticker *time.Ticker
}

func SoftPwmInit(led *gpio.LedDriver, tick time.Duration, numTicks int) *pwmDriver {
	talker := make(chan int, 1)
	return &pwmDriver{led, talker, tick, numTicks, nil}
}

func (p *pwmDriver) Start() {
	if p.ticker != nil {
		return
	}

	go func() {
		at := 0
		duty := 0

		p.ticker = gobot.Every(p.tick, func() {
			if duty < 0 { // won't read any new duty cycles either
				p.led.Off()
				return
			}

			at++
			if at > duty {
				p.led.Off()
			}
			if at == p.cycle {
				p.led.On()
				at = 0
			}

			select {
			case duty = <-p.talker:
			default:
			}
		})
	}()
}

func (p *pwmDriver) Stop() {
	if p.ticker == nil {
		return
	}

	p.talker <- -1 // makes sure the output is off

	time.Sleep(10 * p.tick) // make sure we have a chance for our function to see this and cancel
	p.ticker.Stop()
	p.ticker = nil
}

func (p *pwmDriver) SetDutyCycle(numTicks int) {
	if numTicks < 0 {
		panic("can't have a negative duty cycle")
	}
	p.talker <- numTicks
}

func main() {
	r := raspi.NewAdaptor()
	//led := gpio.NewLedDriver(r, "7") // physical pin 7 => gpio 4
	led := gpio.NewLedDriver(r, "40")        // physical pin 40 => gpio 21
	hbridge1_a := gpio.NewLedDriver(r, "36") // physical pin 36 => gpio 16
	hbridge1_b := gpio.NewLedDriver(r, "38") // physical pin 38 => gpio 20
	pwm := SoftPwmInit(led, 10*time.Nanosecond, 100)

	work := func() {
		led.On()
		pwm.SetDutyCycle(0)
		pwm.Start()

		// change speed constantly
		duty := 0

		gobot.Every(50*time.Millisecond, func() {
			duty++
			//fmt.Println("new duty:", duty)
			if duty == 100 {
				duty = 0
			}
			pwm.SetDutyCycle(duty)
		})

		// reverse directions every once in a while
		hbridge1_a.On()
		hbridge1_b.Off()
		gobot.Every(50*time.Second, func() {
			hbridge1_a.Toggle()
			hbridge1_b.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{led},
		work,
	)
	robot.AutoRun = false

	robot.Start()

	watchForShutdown(pwm,
		hbridge1_a,
		hbridge1_b,
	)
}

func watchForShutdown(pwm *pwmDriver, hbridges ...*gpio.LedDriver) {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel
	fmt.Println("Finishing...")
	for _, hbridge := range hbridges {
		hbridge.Off()
	}
	pwm.Stop()
	fmt.Println("Finished")
	os.Exit(0)
}
