package cmd

var trueFlags []string = []string{"true", "t", "1"}

func ContainsStr(arr []string, target string) bool {
	for _, i := range arr {
		if i == target {
			return true
		}
	}

	return false
}

func ParseBool(v string) bool {
	if ContainsStr(trueFlags, v) {
		return true
	} else {
		return false
	}
}
