package test4

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
