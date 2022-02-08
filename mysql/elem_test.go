// Package mysql
// @Title elem_test.go
// @Description 
// @Author haogooder
// @Update 2022/2/8
package mysql

import (
    "fmt"
    "testing"
)

func TestMakeElemType(t *testing.T) {

    data := "435"
    ref := MakeElemType(data)
    if !ref.IsString() {
        t.Errorf("call IsString failed")
    }
    idata, err := ref.ToInt()
    fmt.Println(idata, err)
    m := map[interface{}]interface{}{
        "k1":    "k1val",
        2:       "k2val",
        "k3":    "k3val",
        4:       "k4val",
        3.44444: 5.88888,
    }
    md := MakeElemType(m)
    sm, err := md.ToSlice()
    fmt.Println(err)
    fmt.Printf("%# v\n", sm)
    fmt.Println(ref.IsSimpleType(), md.IsComplexType())

}
