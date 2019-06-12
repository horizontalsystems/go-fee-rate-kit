package feeratekit

func min(x, y int64) int64 {
    if x < y {
        return x
    }
    return y
}

func max(x, y int64) int64 {
    if x > y {
        return x
    }
    return y
}
