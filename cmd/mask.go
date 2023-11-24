package cmd

import "fmt"

func Mask(raw string) (string) {
	var retVal string
	switch {
		case len(raw) < 8:
			retVal = "********"
		case len(raw) < 16:
			retVal = fmt.Sprintf("%s********%s", raw[0:1], raw[len(raw)-1:])
		default:
			retVal = fmt.Sprintf("%s********%s", raw[0:4], raw[len(raw)-4:])	
	}
	return retVal
}