// Package mysql
// @Title elem.go
// @Description 
// @Author haogooder
// @Update 2022/2/8
package mysql

import (
    "bytes"
    "encoding/binary"
    "encoding/json"
    "errors"
    "fmt"
    "math"
    "reflect"
    "strconv"
    "strings"
    "unsafe"
)

// 任意类型的数据转为结构体
func MakeElemType(d interface{}) ElemType {
    dval := reflect.ValueOf(d)
    for {
        if dval.Kind() == reflect.Ptr {
            dval = dval.Elem()
            continue
        }
        break
    }
    return ElemType{Data: d, RefVal: reflect.ValueOf(d)}
}

// 基本的元素类型
type ElemType struct {
    Data   interface{}   // 数据元素
    RefVal reflect.Value // 通过反射获取的value
}

// 转换为bool类型，如果是bool型则直接返回
func (et ElemType) ToBool() (bool, error) {
    return TryBestToBool(et.Data)
}

// 转换为int类型
func (et ElemType) ToInt() (int, error) {
    tmp, err := et.ToInt64()
    if err != nil {
        return 0, err
    }
    return int(tmp), nil
}

// 转换为int8类型
func (et ElemType) ToInt8() (int8, error) {
    tmp, err := et.ToInt64()
    if err != nil {
        return 0, err
    }
    if tmp > math.MaxInt8 || tmp < math.MinInt8 {
        return 0, fmt.Errorf("toInt8 failed, %v overflow [math.MinInt8,math.MaxInt8]", et.Data)
    }
    return int8(tmp), nil
}

// 转换为int16类型
func (et ElemType) ToInt16() (int16, error) {
    tmp, err := et.ToInt64()
    if err != nil {
        return 0, err
    }
    if tmp > math.MaxInt16 || tmp < math.MinInt16 {
        return 0, fmt.Errorf("toInt16 failed, %v overflow [math.MinInt16,math.MaxInt16]", et.Data)
    }
    return int16(tmp), nil
}

// 转换为int32类型
func (et ElemType) ToInt32() (int32, error) {
    tmp, err := et.ToInt64()
    if err != nil {
        return 0, err
    }
    if tmp > math.MaxInt32 || tmp < math.MinInt32 {
        return 0, fmt.Errorf("toInt32 failed, %v overflow [math.MinInt32,math.MaxInt32]", et.Data)
    }
    return int32(tmp), nil
}

// 转换为int64类型
func (et ElemType) ToInt64() (int64, error) {
    return TryBestToInt64(et.Data)
}

// 转换为uint类型
func (et ElemType) ToUint() (uint, error) {
    tmp, err := et.ToUint64()
    if err != nil {
        return 0, err
    }
    return uint(tmp), nil
}

// 转换为uint8类型
func (et ElemType) ToUint8() (uint8, error) {
    tmp, err := et.ToUint64()
    if err != nil {
        return 0, err
    }
    if tmp > math.MaxUint8 {
        return 0, fmt.Errorf("toUint8 failed, %v overflow", et.Data)
    }
    return uint8(tmp), nil
}

// 转换为uint16类型
func (et ElemType) ToUint16() (uint16, error) {
    tmp, err := et.ToUint64()
    if err != nil {
        return 0, err
    }
    if tmp > math.MaxUint16 {
        return 0, fmt.Errorf("toUint16 failed, %v overflow", et.Data)
    }
    return uint16(tmp), nil
}

// 转换为uint32类型
func (et ElemType) ToUint32() (uint32, error) {
    tmp, err := et.ToUint64()
    if err != nil {
        return 0, err
    }
    if tmp > math.MaxUint32 {
        return 0, fmt.Errorf("toUint32 failed, %v overflow", et.Data)
    }
    return uint32(tmp), nil
}

// 转换为uint64类型
func (et ElemType) ToUint64() (uint64, error) {
    return TryBestToUint64(et.Data)
}

// 转换为string类型
func (et ElemType) ToString() string {
    str, _ := TryBestToString(et.Data)
    return str
}

// 转换为float类型
func (et ElemType) ToFloat32() (float32, error) {
    f64, err := et.ToFloat64()
    if err != nil {
        return 0, err
    }
    if f64 > math.MaxFloat32 {
        return 0, fmt.Errorf("toFloat32 failed, %v overflow", et.Data)
    }
    return float32(f64), nil
}

// 转换为float64类型
func (et ElemType) ToFloat64() (float64, error) {
    return TryBestToFloat(et.Data)
}

// 转换为slice类型
// 原始数据若为array/slice，则直接返回
// 原始数据为map时只返回[]value
// 原始数据若为数字、字符串等简单类型则将其放到slice中返回，即强制转为slice
// 否则，报错
func (et ElemType) ToSlice() ([]ElemType, error) {
    switch et.RefVal.Kind() {
    case reflect.Slice, reflect.Array: // in为slice类型
        vlen := et.RefVal.Len()
        ret := make([]ElemType, vlen)
        for i := 0; i < vlen; i++ {
            ret[i] = MakeElemType(et.RefVal.Index(i).Interface())
        }
        return ret, nil
    case reflect.Map: // in为map类型，取map的value，要不要取map的key呢？
        var ret []ElemType
        ks := et.RefVal.MapKeys()
        for _, k := range ks {
            kiface := et.RefVal.MapIndex(k).Interface()
            ret = append(ret, MakeElemType(kiface))
        }
        return ret, nil
    case reflect.String: // 字符串类型
        tmp := []byte(et.RefVal.String())
        var ret []ElemType
        for _, t := range tmp {
            ret = append(ret, MakeElemType(t))
        }
        return ret, nil
    default: // 其他的类型一律强制转为slice
        return []ElemType{et}, nil
    }
}

// 转换为map类型
// 如果原始数据是map则直接返回
// 如果是json字符串则尝试去解析
// 否则，报错
func (et ElemType) ToMap() (map[ElemType]ElemType, error) {
    ret := make(map[ElemType]ElemType)
    if et.RefVal.Kind() == reflect.Map {
        ks := et.RefVal.MapKeys()
        for _, k := range ks {
            kiface := MakeElemType(k.Interface())
            viface := et.RefVal.MapIndex(k).Interface()
            ret[kiface] = MakeElemType(viface)
        }
        return ret, nil
    }
    if et.RefVal.Kind() == reflect.String {
        str := et.RefVal.String()
        var vmap interface{}
        err := json.Unmarshal([]byte(str), &vmap)
        if err != nil {
            return ret, err
        }
        inRefVal := reflect.ValueOf(vmap)
        if inRefVal.Kind() == reflect.Map {
            ks := inRefVal.MapKeys()
            for _, k := range ks {
                kiface := MakeElemType(k.Interface())
                viface := inRefVal.MapIndex(k).Interface()
                ret[kiface] = MakeElemType(viface)
            }
            return ret, nil
        }
    }
    return ret, fmt.Errorf("cannot convert %v to map", et.Data)
}

// 提取原始数据的长度，只有string/slice/map/array/chan
func (et ElemType) Len() (int, error) {
    switch et.RefVal.Kind() {
    case reflect.String, reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
        return et.RefVal.Len(), nil
    default:
        return 0, fmt.Errorf("invalid type for len %v", et.Data)
    }
}

// 判断原始数据的类型是否为int
func (et ElemType) IsInt() bool {
    return et.Kind() == reflect.Int
}

// 判断原始数据的类型是否为int8
func (et ElemType) IsInt8() bool {
    return et.Kind() == reflect.Int8
}

// 判断原始数据的类型是否为int16
func (et ElemType) IsInt16() bool {
    return et.Kind() == reflect.Int16
}

// 判断原始数据的类型是否为int32
func (et ElemType) IsInt32() bool {
    return et.Kind() == reflect.Int32
}

// 判断原始数据的类型是否为int64
func (et ElemType) IsInt64() bool {
    return et.Kind() == reflect.Int64
}

// 判断原始数据的类型是否为uint
func (et ElemType) IsUint() bool {
    return et.Kind() == reflect.Uint
}

// 判断原始数据的类型是否为uint8
func (et ElemType) IsUint8() bool {
    return et.Kind() == reflect.Uint8
}

// 判断原始数据的类型是否为uint16
func (et ElemType) IsUint16() bool {
    return et.Kind() == reflect.Uint16
}

// 判断原始数据的类型是否为uint32
func (et ElemType) IsUint32() bool {
    return et.Kind() == reflect.Uint32
}

// 判断原始数据的类型是否为uint64
func (et ElemType) IsUint64() bool {
    return et.Kind() == reflect.Uint64
}

// 判断原始数据的类型是否为float32
func (et ElemType) IsFloat32() bool {
    return et.Kind() == reflect.Float32
}

// 判断原始数据的类型是否为float64
func (et ElemType) IsFloat64() bool {
    return et.Kind() == reflect.Float64
}

// 判断原始数据的类型是否为string
func (et ElemType) IsString() bool {
    return et.Kind() == reflect.String
}

// 判断原始数据的类型是否为slice
func (et ElemType) IsSlice() bool {
    return et.Kind() == reflect.Slice
}

// 判断原始数据的类型是否为map
func (et ElemType) IsMap() bool {
    return et.Kind() == reflect.Map
}

// 判断原始数据的类型是否为array
func (et ElemType) IsArray() bool {
    return et.Kind() == reflect.Array
}

// 判断原始数据的类型是否为chan
func (et ElemType) IsChan() bool {
    return et.Kind() == reflect.Chan
}

// 判断原始数据的类型是否为bool
func (et ElemType) IsBool() bool {
    return et.Kind() == reflect.Bool
}

// 是否为字符切片
func (et ElemType) IsByteSlice() bool {
    return reflect.TypeOf(et.Data).String() == "[]uint8"
}

// 是否为简单类型：int/uint/string/bool/float....
func (et ElemType) IsSimpleType() bool {
    return et.IsInt() || et.IsInt8() || et.IsInt16() || et.IsInt32() || et.IsInt64() ||
        et.IsUint() || et.IsUint8() || et.IsUint16() || et.IsUint32() || et.IsUint64() ||
        et.IsString() || et.IsFloat32() || et.IsFloat64() || et.IsBool()
}

// 是否为复合类型：slice/array/map/chan
func (et ElemType) IsComplexType() bool {
    return et.IsSlice() || et.IsMap() || et.IsArray() || et.IsChan()
}

// 原始数据的类型
func (et ElemType) Kind() reflect.Kind {
    return et.RefVal.Kind()
}

// 获取原始数据
func (et ElemType) RawData() interface{} {
    return et.Data
}

type Basickind int

const (
    // 转换时最大值
    MaxInt64Float  = float64(math.MaxInt64)
    MinInt64Float  = float64(math.MinInt64)
    MaxUint64Float = float64(math.MaxUint64)
    // 基本类型归纳，类型转换时用得着
    InvalidKind Basickind = iota
    BoolKind
    ComplexKind
    IntKind
    FloatKind
    StringKind
    UintKind
    PtrKind
    ContainerKind
    FuncKind
)

var (
    // 一些错误信息
    ErrorOverflowMaxInt64  = errors.New("this value overflow math.MaxInt64")
    ErrorOverflowMaxUint64 = errors.New("this value overflow math.MaxUint64")
    ErrorLessThanMinInt64  = errors.New("this value less than math.MinInt64")
    ErrorLessThanZero      = errors.New("this value less than zero")
    ErrorBadComparisonType = errors.New("invalid type for comparison")
    ErrorBadComparison     = errors.New("incompatible types for comparison")
    ErrorNoComparison      = errors.New("missing argument for comparison")
    ErrorInvalidInputType  = errors.New("invalid input type")
)

// 转换成特定类型，便于判断
func GetBasicKind(v reflect.Value) (Basickind, error) {
    switch v.Kind() {
    case reflect.Bool:
        return BoolKind, nil
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return IntKind, nil
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
        return UintKind, nil
    case reflect.Float32, reflect.Float64:
        return FloatKind, nil
    case reflect.Complex64, reflect.Complex128:
        return ComplexKind, nil
    case reflect.String:
        return StringKind, nil
    case reflect.Ptr:
        return PtrKind, nil
    case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
        return ContainerKind, nil
    case reflect.Func:
        return FuncKind, nil
    }
    return InvalidKind, ErrorInvalidInputType
}

// int64转byte
func Int64ToBytes(i int64) []byte {
    var buf = make([]byte, 8)
    binary.BigEndian.PutUint64(buf, uint64(i))
    return buf
}

// bytes转int64
func BytesToInt64(buf []byte) int64 {
    return int64(binary.BigEndian.Uint64(buf))
}

// float64转byte
func Float64ToByte(float float64) []byte {
    bits := math.Float64bits(float)
    bs := make([]byte, 8)
    binary.BigEndian.PutUint64(bs, bits)
    return bs
}

// bytes转float
func ByteToFloat64(bytes []byte) float64 {
    bits := binary.BigEndian.Uint64(bytes)
    return math.Float64frombits(bits)
}

// 字符串转为字节切片
func StrToByte(s string) []byte {
    x := (*[2]uintptr)(unsafe.Pointer(&s))
    h := [3]uintptr{x[0], x[1], x[1]}
    return *(*[]byte)(unsafe.Pointer(&h))
}

// 字节切片转为字符串
func ByteToStr(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}

// 尽最大努力将给定的类型转换为uint64
// 如"45.67"->45、"98.4abc3"->98、"34.87"->34
func TryBestToUint64(value interface{}) (uint64, error) {
    ret, err := tryBestConvertAnyTypeToInt(value, true)
    if err != nil {
        return 0, err
    }
    val := reflect.ValueOf(ret)
    if val.Kind() != reflect.Uint64 {
        return 0, ErrorInvalidInputType
    }
    return val.Uint(), nil
}

// 尽最大努力将给定的类型转换为int64
// 如"45.67"->45、"98.4abc3"->98、"34.87"->34
func TryBestToInt64(value interface{}) (int64, error) {
    ret, err := tryBestConvertAnyTypeToInt(value, false)
    if err != nil {
        return 0, err
    }
    val := reflect.ValueOf(ret)
    if val.Kind() != reflect.Int64 {
        return 0, ErrorInvalidInputType
    }
    return val.Int(), nil
}

// 从左边开始提取数据及小数点
func getFloatStrFromLeft(val string) string {
    val = strings.TrimSpace(val)
    valBytes := StrToByte(val)
    buf := bytes.Buffer{}
    for _, b := range valBytes {
        if b >= 48 && b <= 57 || b == 46 {
            buf.WriteByte(b)
            continue
        }
        break
    }
    return buf.String()
}

// 尽最大努力将任意类型转为int64或uint64
func tryBestConvertAnyTypeToInt(value interface{}, isUnsigned bool) (interface{}, error) {
    val := reflect.ValueOf(value)
    basicKind, err := GetBasicKind(val)
    if err != nil {
        return 0, err
    }
    switch basicKind {
    case IntKind:
        v := val.Int()
        if isUnsigned {
            if v >= 0 {
                return uint64(v), nil
            }
            return 0, ErrorLessThanZero
        }
        return v, nil
    case UintKind:
        v := val.Uint()
        if isUnsigned {
            return v, nil
        }
        if v > math.MaxInt64 {
            return 0, ErrorOverflowMaxInt64
        }
        return int(v), nil
    case StringKind: // 取连续的最长的数字或小数点
        floatStr := getFloatStrFromLeft(val.String())
        if len(floatStr) <= 0 {
            if isUnsigned {
                return uint64(0), nil
            }
            return int64(0), nil
        }
        // 先转成float，因为将"45.33"直接转为int/uint时会报错
        f, err := strconv.ParseFloat(floatStr, 10)
        if err != nil {
            return 0, err
        }
        return tryBestConvertAnyTypeToInt(f, isUnsigned)
        // float特殊处理，会有科学记数法表示形式
    case FloatKind:
        f := val.Float()
        if isUnsigned {
            if f > MaxUint64Float {
                return 0, ErrorOverflowMaxUint64
            }
            if f < 0 {
                return 0, ErrorLessThanZero
            }
            return uint64(f), nil
        }
        if f > MaxInt64Float {
            return 0, ErrorOverflowMaxInt64
        }
        if f < MinInt64Float {
            return 0, ErrorLessThanMinInt64
        }
        return int64(f), nil
    case BoolKind:
        b := val.Bool()
        tmp := 0
        if b {
            tmp = 1
        }
        if isUnsigned {
            return uint64(tmp), nil
        }
        return int64(tmp), nil
        // 指针类型递归调用，直到取本值为止
    case PtrKind:
        if val.IsNil() {
            if isUnsigned {
                return uint64(0), nil
            }
            return int64(0), nil
        }
        return tryBestConvertAnyTypeToInt(val.Elem().Interface(), isUnsigned)
    default:
        return 0, ErrorInvalidInputType
    }
}

// 尽最大努力转换为字符串
func TryBestToString(value interface{}) (string, error) {
    val := reflect.ValueOf(value)
    basicKind, err := GetBasicKind(val)
    if err != nil {
        return "", err
    }
    switch basicKind {
    case IntKind:
        return strconv.FormatInt(val.Int(), 10), nil
    case UintKind:
        return strconv.FormatUint(val.Uint(), 10), nil
    case StringKind:
        return val.String(), nil
    case FloatKind:
        return strconv.FormatFloat(val.Float(), 'f', -1, 64), nil
    case BoolKind:
        return strconv.FormatBool(val.Bool()), nil
    case PtrKind:
        if val.IsNil() {
            return "nil", nil
        }
        return TryBestToString(val.Elem().Interface())
    case ContainerKind:
        result, err := json.Marshal(value)
        if err != nil {
            return "", err
        }
        return string(result), err
    default:
        return val.String(), nil
    }
}

// 尽最大努力转换为float64
func TryBestToFloat(value interface{}) (float64, error) {
    val := reflect.ValueOf(value)
    basicKind, err := GetBasicKind(val)
    if err != nil {
        return 0, err
    }
    switch basicKind {
    case IntKind:
        return float64(val.Int()), nil
    case UintKind:
        return float64(val.Uint()), nil
    case StringKind:
        floatStr := getFloatStrFromLeft(val.String())
        if len(floatStr) <= 0 {
            return 0, nil
        }
        return strconv.ParseFloat(floatStr, 10)
    case FloatKind:
        return val.Float(), nil
    case BoolKind:
        if val.Bool() {
            return 1, nil
        }
        return 0, nil
    case PtrKind:
        if val.IsNil() {
            return 0, nil
        }
        return TryBestToFloat(val.Elem().Interface())
    default:
        return 0, ErrorInvalidInputType
    }
}

// 尽最大努力转为bool类型
func TryBestToBool(value interface{}) (bool, error) {
    val := reflect.ValueOf(value)
    basicKind, err := GetBasicKind(val)
    if err != nil {
        return false, err
    }
    switch basicKind {
    case FloatKind:
        return val.Float() != 0, nil
    case IntKind:
        return val.Int() != 0, nil
    case UintKind:
        return val.Uint() != 0, nil
    case StringKind:
        v := strings.TrimSpace(val.String())
        if len(v) > 0 {
            return true, nil
        }
        return false, nil
    case BoolKind:
        return val.Bool(), nil
    case PtrKind:
        if val.IsNil() {
            return false, nil
        }
        return TryBestToBool(val.Elem().Interface())
    case FuncKind:
        return !val.IsNil(), nil
    }

    // 对于Array, Chan, Map, Slice长度>0即可
    switch val.Kind() {
    case reflect.Array, reflect.Chan, reflect.Slice, reflect.Map:
        return val.Len() != 0, nil
    }
    return false, ErrorInvalidInputType
}
