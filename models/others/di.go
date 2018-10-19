package others

import "time"


type DICache struct{
	Values []int   //value
	Time time.Time //feild
}
type DITitle struct{
	Params string  //title
}
type DI struct{
	DITitles DITitle
	DICaches DICache
}
