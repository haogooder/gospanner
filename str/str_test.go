package str

import (
    "fmt"
    "testing"
)

func TestCapitalize(t *testing.T) {
    t.Run("首字母大写正常测试", func(t *testing.T) {
        s := "abcdefg"
        s = Capitalize(s)
        if s != "Abcdefg" {
            t.Fatal("fail")
        }
    })
    t.Run("首字母大写空字符串测试", func(t *testing.T) {
        s := ""
        s = Capitalize(s)
        if s != "" {
            t.Fatal("fail")
        }
    })
}

func TestSplitLines(t *testing.T) {
    t.Run("按行分割保留换行符", func(t *testing.T) {
        s := "abc\n1\n\ndef\n"
        lines := SplitLines(s, true)
        if len(lines) != 5 {
            t.Fatal("fail")
        }
        if lines[0] != "abc\n" {
            t.Fatal("fail")
        }
    })
    t.Run("按行分割不保留换行符", func(t *testing.T) {
        s := "abc\n1\n\ndef\n"
        lines := SplitLines(s, false)
        if len(lines) != 5 {
            t.Fatal("fail")
        }
        if lines[0] != "abc" {
            fmt.Println(lines[0])
            t.Fatal("fail")
        }
    })
}
