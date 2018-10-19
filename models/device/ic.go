package device

import "time"

type IC struct{
	MotorId string
	SpindleTemperature2 float32
	SpindleTemperature1 float32
	Current2 float32
	Current float32
	Time time.Time
}
