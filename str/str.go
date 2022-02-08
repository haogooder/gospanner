// Package str
// @Title
// @Description
// @Author
// @Update
package str

import (
    "strings"
    "unicode"
)

// Capitalize
// 首字母大写
func Capitalize(s string) string {
    if s == "" {
        return ""
    }
    return strings.ToUpper(s[:1]) + s[1:]
}

// Count
// 返回substr在s中出现的次数
func Count(s, substr string) int {
    return strings.Count(s, substr)
}

func Compare(a, b string) int {
    return strings.Compare(a, b)
}

func Contains(s, substr string) bool {
    return strings.Contains(s, substr)
}

func ContainsAny(s, chars string) bool {
    return strings.ContainsAny(s, chars)
}

// Endswith
// 判断是否以suffix结尾
func Endswith(s, suffix string) bool {
    return strings.HasSuffix(s, suffix)
}

// Startswith
// 判断是否以prefix开头
func Startswith(s, prefix string) bool {
    return strings.HasPrefix(s, prefix)
}

func Find(s, substr string) int {
    return strings.Index(s, substr)
}

func FindLast(s, substr string) int {
    return strings.LastIndex(s, substr)
}

func Index(s, substr string) int {
    return strings.Index(s, substr)
}

func IndexLast(s, substr string) int {
    return strings.LastIndex(s, substr)
}

func IsAlphaNum(s string) bool {
    if len(s) == 0 {
        return false
    }
    for _, char := range s {
        if unicode.IsDigit(char) || unicode.IsLetter(char) {
            continue
        } else {
            return false
        }
    }
    return true
}

func IsAlpha(s string) bool {
    if len(s) == 0 {
        return false
    }
    for _, char := range s {
        if unicode.IsLetter(char) {
            continue
        } else {
            return false
        }
    }
    return true
}

func IsDigit(s string) bool {
    if len(s) == 0 {
        return false
    }
    for _, char := range s {
        if unicode.IsDigit(char) {
            continue
        } else {
            return false
        }
    }
    return true
}

func IsLower(s string) bool {
    if len(s) == 0 {
        return false
    }
    for _, char := range s {
        if unicode.IsLower(char) {
            continue
        } else {
            return false
        }
    }
    return true
}

func IsSpace(s string) bool {
    if len(s) == 0 {
        return false
    }
    for _, char := range s {
        if unicode.IsSpace(char) {
            continue
        } else {
            return false
        }
    }
    return true
}

// func IsTitle(s string) bool {
//     if len(s) == 0 {
//         return false
//     }
// TODO 待完成
// }

func IsUpper(s string) bool {
    if len(s) == 0 {
        return false
    }
    for _, char := range s {
        if unicode.IsUpper(char) {
            continue
        } else {
            return false
        }
    }
    return true
}

func Join(sep string, elems ...string) string {
    return strings.Join(elems, sep)
}

func Len(s string) int {
    return len(s)
}

func Lower(s string) string {
    return strings.ToLower(s)
}

func Upper(s string) string {
    return strings.ToUpper(s)
}

func Strip(s string, cutset string) string {
    return strings.Trim(s, cutset)
}

func StripPrefix(s string, cutset string) string {
    return strings.TrimPrefix(s, cutset)
}

func StripSuffix(s string, cutset string) string {
    return strings.TrimSuffix(s, cutset)
}

func Fields(s string) []string {
    return strings.Fields(s)
}

func Split(s, sep string) []string {
    return strings.Split(s, sep)
}

func SplitLines(s string, keepEnds bool) []string {
    if keepEnds {
        return strings.SplitAfter(s, "\n")
    } else {
        temp := strings.SplitAfter(s, "\n")
        for i, e := range temp {
            temp[i] = strings.TrimRight(e, "\n")
        }
        return temp
    }
}

func Max(s string) string {
    if s == "" {
        return ""
    }
    var max rune
    for i, char := range s {
        if i == 0 {
            max = char
            continue
        }
        if char > max {
            max = char
        }
    }
    return string(max)
}

func Min(s string) string {
    if s == "" {
        return ""
    }
    var min rune
    for i, char := range s {
        if i == 0 {
            min = char
            continue
        }
        if char < min {
            min = char
        }
    }
    return string(min)
}

func Replace(s, old, new string, max int) string {
    return strings.Replace(s, old, new, max)
}

func ReplaceAll(s, old, new string) string {
    return strings.ReplaceAll(s, old, new)
}

func Title(s string) string {
    return strings.Title(s)
}
