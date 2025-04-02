package utils

func IfArrayContains(ar []string, val string) bool {
	for _, v := range ar {
		if v == val {
			return true
		}
	}
	return false
}

func IfMapContains(m map[string]interface{}, val string) bool {
	_, ok := m[val]
	return ok
}
