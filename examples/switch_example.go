package main

// getDayName returns the name of the day based on a number
func getDayName(day int) int {
	switch day {
	case 0:
		return 1 // Sunday
	case 1:
		return 2 // Monday
	case 2:
		return 3 // Tuesday
	case 3:
		return 4 // Wednesday
	case 4:
		return 5 // Thursday
	case 5:
		return 6 // Friday
	case 6:
		return 7 // Saturday
	default:
		return 0 // Invalid
	}
}

func main() {
	result := getDayName(3)
	println(result)
}
