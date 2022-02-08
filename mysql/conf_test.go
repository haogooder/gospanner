// Package mysql
// @Title conf_test.go
// @Description 
// @Author haogooder
// @Update 2022/2/8
package mysql

import (
    "fmt"
    "testing"
)

func TestFormatCond(t *testing.T) {

    delim := "and"
    cond := make(map[string]interface{})
    cond["id:in"] = "444,666,888"
    cond["uid:in"] = 444
    cond["name:like"] = "liuyongshuai"
    cond["title:rlike"] = []string{"sina", "baidu"} // 非法
    cond["tid:lt"] = 400
    cond["tags:find"] = "google"
    cond["category:find"] = []interface{}{444, "didi"} // 非法
    cond["cnum"] = "99999"
    cond["praiseNum"] = []interface{}{"aaaa", 999} // 非法
    sqlCond, param := FormatCond(cond, delim)
    /**
      sqlCond "`id` IN (?,?,?) AND `tid` < ? AND FIND_IN_SET(?,`tags`) AND `cnum` = ? AND `uid` IN (?) AND `name` LIKE ?"
      param {}{"444","666","888",int(400),"google","99999",int(444),"%liuyongshuai%",}
    */
    fmt.Printf("sqlCond %# v\n", sqlCond)
    fmt.Printf("param %# v\n", param)

}
