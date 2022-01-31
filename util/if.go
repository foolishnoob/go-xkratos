package util

func If(isTrue bool, yes interface{}, no interface{}) interface{} {
	if isTrue {
		return yes
	}
	return no
}
