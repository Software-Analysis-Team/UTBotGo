package main

type RuneScanner struct {
	InitialValue     interface{}
	MutationFunction func(rune, interface{}) interface{}
}

func ScanRunes(rs []rune, s RuneScanner) (result interface{}) {
	result = s.InitialValue
	for _, r := range rs {
		result = s.MutationFunction(r, result)
	}
	return
}

func ScanBytes(bs []byte, s RuneScanner) interface{} {
	return ScanRunes([]rune(string(bs)), s)
}
