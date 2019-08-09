package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
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
