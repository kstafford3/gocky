package gocky

func contains(searchSpace []string, value string) bool {
	for _, searchSpaceValue := range searchSpace {
		if value == searchSpaceValue {
			return true
		}
	}
	return false
}