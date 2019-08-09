package main

import (
	"fmt"

	"gobot.io/x/gobot/drivers/gpio"
)

type Motor interface {
	Fwd()
	Back()
	Speed(int64)
}

func NewFakeMotor() Motor {
	return &nullMotor{}
}

func NewMotor(speedController *pwmDriver, hbridge1 *gpio.LedDriver, hbridge2 *gpio.LedDriver) Motor {
	return &carMotor{
		speedController,
		hbridge1,
		hbridge2,
	}
}

type carMotor struct {
	speedController *pwmDriver
	hbridge1        *gpio.LedDriver
	hbridge2        *gpio.LedDriver
}

func (car *carMotor) Fwd() {
	fmt.Println("going forward")
}

func (car *carMotor) Back() {
	fmt.Println("going backward")
}

func (car *carMotor) Speed(speed int64) {
	fmt.Printf("going %d\n", speed)
}

// a do-nothing "motor", used as a hack to allow for no index checking
type nullMotor struct {
}

func (*nullMotor) Fwd() {
}

func (*nullMotor) Back() {
}

func (*nullMotor) Speed(int64) {
}

// does nothing, just prints back to client debugging purposes
func handleSpeed(speed int64, slider int64) string {
	return fmt.Sprintf("received speed %d for slider %d", speed, slider)
}
