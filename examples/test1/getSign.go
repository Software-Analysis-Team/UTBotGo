package main

func getSign(x int) int {
    if x == 0 {
        return 0
    }
    if x > 0 {
        return 1
    } else {
        return -1
    }
}
