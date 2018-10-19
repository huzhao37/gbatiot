package device

import "time"

type CY struct{
	MotorId string
	AccumulativeWeight float32
	TotalPower float32
	InstantWeight float32
	Current float32
	Unit   float32
	BootFlagBit bool
	ZeroCalibration  bool
	Time time.Time
}

