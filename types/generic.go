package types

var AndEmptyString = ""
var AndTrue = true
var AndFalse = false

//when you have a slice of string
// you want to remove a specific value from an index
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// when you want to find the index of the
func FindIndex(s []string, element string) int {
	for p, v := range s {
		if v == element {
			return p
		}
	}
	return -1
}
