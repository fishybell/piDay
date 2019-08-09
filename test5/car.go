package main

import (
	"fmt"

	"gobot.io/x/gobot/drivers/gpio"
)

type Motor interface {
	Fwd()
	Back()
	Speed(int64)
	Name() string
}

func NewFakeMotor() Motor {
	return &nullMotor{}
}

func NewMotor(speedController *pwmDriver, hbridge1 *gpio.LedDriver, hbridge2 *gpio.LedDriver, name string) Motor {
	return &carMotor{
		speedController,
		hbridge1,
		hbridge2,
		name,
	}
}

type carMotor struct {
	speedController *pwmDriver
	hbridge1        *gpio.LedDriver
	hbridge2        *gpio.LedDriver
	name            string
}

func (car *carMotor) Fwd() {
	fmt.Println(car.name, "going forward")
	car.hbridge1.On()
	car.hbridge2.Off()
}

func (car *carMotor) Back() {
	fmt.Println(car.name, "going backward")
	car.hbridge1.Off()
	car.hbridge2.On()
}

func (car *carMotor) Speed(speed int64) {
	if speed == 0 {
		fmt.Println(car.name, "stopping")
		car.hbridge1.Off()
		car.hbridge2.Off()
	} else {
		fmt.Printf("%s going %d\n", car.name, speed)
		car.speedController.SetDutyCycle(int(speed))
	}
}

func (car *carMotor) Name() string {
	return car.name
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

func (*nullMotor) Name() string {
	return "you should never see this"
}

func handleSpeed(motors []Motor, speed int64, slider int64) string {
	if slider <= 0 {
		return "I'm afraid I can't do that"
	}

	if speed < 0 {
		motors[slider].Back()
		motors[slider].Speed(-1 * speed)
		motors[slider+2].Back()
		motors[slider+2].Speed(-1 * speed)
	} else {
		motors[slider].Fwd()
		motors[slider].Speed(speed)
		motors[slider+2].Fwd()
		motors[slider+2].Speed(speed)
	}
	return fmt.Sprintf("received speed %d for slider %d", speed, slider)
}
