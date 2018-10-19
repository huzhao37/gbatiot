package device

import "time"

type MF struct{
	MotorId string
	InFrequency float32
	Frequency float32
	Time time.Time
}
