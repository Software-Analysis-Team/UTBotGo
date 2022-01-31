package main

type Name string

type NameEncoder struct{}

var nameEncoder = NameEncoder{}

func (nameEncoder NameEncoder) Encode(name Name) Name {
	return "*" + name
}

func (nameEncoder NameEncoder) Decode(name Name) (Name, bool) {
	if name[0] == '*' {
		return name[1:], true
	}
	return "", false
}
