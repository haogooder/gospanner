package conv

import (
    "testing"
)

func TestStrToInt(t *testing.T) {
    t.Run("字符串转数字", func(t *testing.T) {
        s := "10"
        i := StrToInt(s)
        if i != 10 {
            t.Fatal("fail")
        }
    })
}

func TestIntToStr(t *testing.T) {
    t.Run("数字转字符串", func(t *testing.T) {
        i := 10
        s := IntToStr(i)
        if s != "10" {
            t.Fatal("fail")
        }
    })
}

func TestFloatToStr(t *testing.T) {
    t.Run("浮点数转字符串", func(t *testing.T) {
        f := 3.1415926
        s := FloatToStr(f)
        if s != "3.1415926" {
            t.Fatal("fail")
        }
    })

}

func TestStrToFloat(t *testing.T) {
    t.Run("字符串转浮点数", func(t *testing.T) {
        s := "3.1415926"
        f := StrToFloat(s)
        if f != 3.1415926 {
            t.Fatal("fail")
        }
    })
}
