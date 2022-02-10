package utils

func InSlice(a interface{}, list []interface{}) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Pos(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}
