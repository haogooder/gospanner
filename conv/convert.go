package conv

import (
    "strconv"
)

func StrToInt(s string) int {
    i, err := strconv.Atoi(s)
    if err != nil {
        panic(err.Error())
    }
    return i
}

func IntToStr(i int) string {
    return strconv.Itoa(i)
}

func StrToInt64(s string, base int) int64 {
    i, err := strconv.ParseInt(s, base, 64)
    if err != nil {
        panic(err.Error())
    }
    return i
}

func Int64ToStr(i int64, base int) string {
    return strconv.FormatInt(i, base)
}

func StrToBytes(s string) []byte {
    return []byte(s)
}

func BytesToStr(bytes []byte) string {
    return string(bytes)
}

func StrToFloat(s string) float64 {
    f, err := strconv.ParseFloat(s, 64)
    if err != nil {
        panic(err)
    }
    return f
}

func FloatToStr(f float64) string {
    s := strconv.FormatFloat(f, 'f', -1, 64)
    return s
}
