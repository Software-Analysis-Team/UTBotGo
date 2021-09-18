package main

import (
    "math/rand"
    "fmt"
    "strconv"
)

func GenerateInt () int {
    return int(rand.Uint64())
}

func GenerateString (length int) (res string) {
    n := rand.Intn(length)
    res = ""
    for i := 0; i < n; i++ {
        res += string(rand.Intn(128))
    }
    return
}

func GenerateData (type_ Type) Data {
    var str string
    switch type_ {
    case "int":
        str = strconv.Itoa(GenerateInt())
    case "string":
        str = fmt.Sprintf("%#v", GenerateString(42))
    default:
        str = "nil"
    }
    return Data(str)
}
