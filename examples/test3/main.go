package test4

func f(n int) (res string) {
    if n > 5 {
        return "n > 5, sorry"
    }
    for i := 0; i < n; i++ {
        res = res + "*"
    }
    return
}

func g(a byte, b rune, c uint32) (d int) {
    d = int(a)
    if d < int(b) {
        d = int(b)
    }
    if d < int(c) {
        d = int(c)
    }
    return
}
