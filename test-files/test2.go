package simple_test_2

import "fmt"

func PrintSum(a, b int) {
    fmt.Println("a + b =", a + b)
}

func LengthsSum(s1 string, s2 string) (n int) {
    n = len(s1) + len(s2)
    if n > 80 { panic("really panic!!!") }
    return
}
