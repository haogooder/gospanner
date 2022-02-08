// Package uid
// @Title fakeid.go
// @Description 
// @Author haogooder
// @Update 2022/2/8
package uid

import (
    "fmt"
    "math/rand"
    "time"
)

var (
    LocalHostIP         = ""
    LocalHostIpArr      []string
    LocalHostIpTraceId  = ""
    SequenceIDGenerator SnowFlakeIdGenerator
    preTraceID          = ""
    ScreenWidth         int
    ScreenHeight        int
)

// 生成一个假的traceId
func FakeTraceId() (traceId string) {
    for {
        reRandSeed()
        traceId = fmt.Sprintf("%x%s%x", time.Now().UnixNano(), LocalHostIpTraceId, rand.Int63())
        if preTraceID != traceId {
            preTraceID = traceId
            break
        }
    }
    return traceId
}

// 重新设置随机数种子
func reRandSeed() {
    genId, err := SequenceIDGenerator.NextId()
    if err != nil {
        rand.Seed(time.Now().UnixNano())
    } else {
        rand.Seed(genId)
    }
}
