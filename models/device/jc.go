package device

import "time"

type JC struct{
	MotorId string
	Voltage float32
	MotiveSpindleTemperature2 float32
	MotiveSpindleTemperature1 float32
	RackSpindleTemperature2 float32
	RackSpindleTemperature1 float32
	Current float32
	TotalPower float32
	Time time.Time
	}
