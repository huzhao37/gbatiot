package device

import "time"

type SIM struct{
	SimId string //==motorid
	Count int32
	Time time.Time
}
