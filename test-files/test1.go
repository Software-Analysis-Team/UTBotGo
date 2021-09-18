package test1

func f (x, y int, z string) int {
    if len(z) > 10 { panic(z) }
    return x + 2 * y - len(z)
}

func g (a string) {
    if len(a) > 15 { panic(a) }
}
