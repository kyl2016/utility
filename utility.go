package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

func PanicIfNotNil(err error) {
	if err == nil {
		return
	}
	log.Println(err)
	debug.PrintStack()
	panic(err)
}

var envVars map[string]string

func EnvironmentVariables() map[string]string {
	if envVars != nil {
		return envVars
	}
	lines := os.Environ()
	envVars = make(map[string]string, len(lines))
	for _, line := range lines {
		comps := strings.Split(line, "=")
		if len(comps) > 1 {
			envVars[comps[0]] = comps[1]
		}
	}
	return envVars
}

type StrMap = map[string]interface{}
type AnyMap = map[interface{}]interface{}

func AnyToAnyMap(value interface{}) AnyMap {
	if value == nil {
		return nil
	}
	switch val := value.(type) {
	case AnyMap:
		return val
	case StrMap:
		count := len(val)
		if count == 0 {
			return nil
		}
		m := make(AnyMap, count)
		for k, v := range val {
			m[k] = v
		}
		return m
	default:
		return nil
	}
}

func AnyToStrMap(value interface{}) StrMap {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case StrMap:
		return v
	case AnyMap:
		l := len(v)
		if l == 0 {
			return nil
		}
		m := make(StrMap, l)
		for k, v := range v {
			m[AnyToString(k)] = v
		}
		return m
	default:
		return nil
	}
}

func AnyToString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch val := value.(type) {
	case *string:
		if val == nil {
			return ""
		}
		return *val
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case error:
		return val.Error()
	default:
		return fmt.Sprint(value)
	}
}

func AnyToInt64(value interface{}) int64 {
	if value == nil {
		return 0
	}
	switch val := value.(type) {
	case int:
		return int64(val)
	case int8:
		return int64(val)
	case int16:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return val
	case uint:
		return int64(val)
	case uint8:
		return int64(val)
	case uint16:
		return int64(val)
	case uint32:
		return int64(val)
	case uint64:
		return int64(val)
	case *string:
		if val == nil {
			return 0
		}
		if i, err := StringToInt64(*val); err == nil {
			return i
		}
	case string:
		if i, err := StringToInt64(val); err == nil {
			return i
		}
	case float32:
		return int64(val)
	case float64:
		return int64(val)
	case bool:
		if val {
			return 1
		}
		return 0
	case json.Number:
		v, _ := val.Int64()
		return v
	}
	return 0
}

func AnyToFloat64(value interface{}) float64 {
	if value == nil {
		return 0
	}
	switch val := value.(type) {
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return val
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case *string:
		if val == nil {
			return 0
		}
		if v, err := strconv.ParseFloat(*val, 64); err == nil {
			return v
		}
	case string:
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			return v
		}
	case bool:
		if val {
			return 1
		}
		return 0
	case json.Number:
		v, _ := val.Float64()
		return v
	}
	return 0
}

func AnyToBool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch v := v.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case int8:
		return v != 0
	case int16:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case uint:
		return v != 0
	case uint8:
		return v != 0
	case uint16:
		return v != 0
	case uint32:
		return v != 0
	case uint64:
		return v != 0
	case float32:
		return v != 0
	case float64:
		return v != 0
	case string:
		if len(v) == 0 {
			return false
		}
		c := strings.ToLower(v[0:1])
		return c == "y" || c == "t" || c == "1"
	case *string:
		return v != nil && AnyToBool(*v)
	default:
		return false
	}
}

func AnyToInt(value interface{}) int {
	if value == nil {
		return 0
	}
	switch val := value.(type) {
	case int:
		return int(val)
	case int8:
		return int(val)
	case int16:
		return int(val)
	case int32:
		return int(val)
	case int64:
		return int(val)
	case uint:
		return int(val)
	case uint8:
		return int(val)
	case uint16:
		return int(val)
	case uint32:
		return int(val)
	case uint64:
		return int(val)
	case *string:
		v, err := strconv.Atoi(*val)
		if err != nil {
			return 0
		}
		return v
	case string:
		v, err := strconv.Atoi(val)
		if err != nil {
			return 0
		}
		return v
	case float32:
		return int(val)
	case float64:
		return int(val)
	case bool:
		if val {
			return 1
		} else {
			return 0
		}
	case json.Number:
		v, _ := val.Int64()
		return int(v)
	}
	return 0
}

func CheckWithRangeRandom(values interface{}) float64 {
	res, err := AnyToFloatWithRangeRandom(values)
	if err != nil {
		return AnyToFloat64(values)
	}
	return res
}

func CalculateRandom(minValue, maxValue int) int {
	if minValue >= maxValue {
		return minValue
	}
	return rand.Intn(maxValue-minValue) + minValue
}
func AnyToFloatWithRangeRandom(values interface{}) (float64, error) {

	interRange, ok := values.([]interface{})
	if ok && len(interRange) == 2 {
		return float64(CalculateRandom(AnyToInt(interRange[0]), AnyToInt(interRange[1]))), nil
	}

	inter2Range, ok := values.([2]interface{})
	if ok && len(inter2Range) == 2 {

		return float64(CalculateRandom(AnyToInt(inter2Range[0]), AnyToInt(inter2Range[1]))), nil
	}

	intRange, ok := values.([]int)
	if ok && len(intRange) == 2 {
		return float64(CalculateRandom(intRange[0], intRange[1])), nil
	}
	int2Range, ok := values.([2]int)
	if ok && len(int2Range) == 2 {

		return float64(CalculateRandom(int2Range[0], int2Range[1])), nil
	}
	floatRange, ok := values.([]float64)
	if ok && len(floatRange) == 2 {

		return float64(CalculateRandom(int(floatRange[0]), int(floatRange[1]))), nil
	}

	float2Range, ok := values.([2]float64)
	if ok && len(float2Range) == 2 {

		return float64(CalculateRandom(int(float2Range[0]), int(float2Range[1]))), nil
	}
	return 0, errors.New("infoMap value error")
}

func AnyArrayToStrMap(mapInterface []interface{}) StrMap {
	if len(mapInterface)/2 < 1 {
		return nil
	}
	elementMap := make(StrMap)
	for i := 0; i < len(mapInterface)/2; i += 1 {
		key := AnyToString(mapInterface[i*2])
		elementMap[key] = mapInterface[i*2+1]
	}
	return elementMap
}

func AnyToStringArray(any interface{}) []string {
	if any == nil {
		return nil
	}
	switch v := any.(type) {
	case []string:
		return v
	case []interface{}:
		return AnyArrayToStringArray(v)
	default:
		return nil
	}
}

func AnyArrayToStringArray(arrInterface []interface{}) []string {
	elementArray := make([]string, len(arrInterface))
	for i, v := range arrInterface {
		elementArray[i] = AnyToString(v)
	}
	return elementArray
}

func StringToInt64(value string) (int64, error) {
	if index := strings.Index(value, "."); index > 0 {
		value = value[:index]
	}
	return strconv.ParseInt(value, 10, 64)
}

// BytesToString 按string的底层结构，转换[]byte
func BytesToString(b []byte) string {
	if b == nil {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes 按[]byte的底层结构，转换字符串，len与cap皆为字符串的len
func StringToBytes(s string) []byte {
	return StringPToBytes(&s)
}

func StringPToBytes(s *string) []byte {
	if s == nil {
		return nil
	}
	// 获取s的起始地址开始后的两个 uintptr 指针
	x := (*[2]uintptr)(unsafe.Pointer(s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func FindInSyncMap(m *sync.Map, keys ...interface{}) interface{} {
	return FindInSyncMapWithKeys(m, keys)
}

func FindInSyncMapWithKeys(m *sync.Map, keys []interface{}) interface{} {
	if m == nil {
		return nil
	}
	l := len(keys)
	if l == 0 {
		return nil
	}
	v0, ok := m.Load(keys[0])
	if !ok || l == 1 {
		return v0
	}
	switch v := v0.(type) {
	case StrMap:
		return FindInStrMapWithKeys(v, keys[1:])
	case AnyMap:
		return FindInAnyMapWithKeys(v, keys[1:])
	default:
		return nil
	}
}

func FindInAnyMap(m AnyMap, keys ...interface{}) interface{} {
	return FindInAnyMapWithKeys(m, keys)
}

func FindInAnyMapWithKeys(m AnyMap, keys []interface{}) interface{} {
	if m == nil {
		return nil
	}
	l := len(keys)
	if l == 0 {
		return nil
	}
	value := m[keys[0]]
	if l == 1 {
		return value
	}
	switch v := value.(type) {
	case AnyMap:
		return FindInAnyMapWithKeys(v, keys[1:])
	case StrMap:
		return FindInStrMapWithKeys(v, keys[1:])
	default:
		return nil
	}
}

func FindInStrMap(m StrMap, keys ...interface{}) interface{} {
	return FindInStrMapWithKeys(m, keys)
}

func FindInStrMapWithKeys(m StrMap, keys []interface{}) interface{} {
	if m == nil {
		return nil
	}
	l := len(keys)
	if l == 0 {
		return nil
	}
	value := m[AnyToString(keys[0])]
	if l == 1 {
		return value
	}
	switch v := value.(type) {
	case AnyMap:
		return FindInAnyMapWithKeys(v, keys[1:])
	case StrMap:
		return FindInStrMapWithKeys(v, keys[1:])
	default:
		return nil
	}
}

//flatten map ,e.g
// map A
// {
// 	"foo":{
// 		"bar":1
// 	}
// }
// map B = FlattenMap("",".",A)
// {
// 	"foo.bar":1
// }

func FlattenMap(rootKey, delimiter string, originData StrMap) StrMap {
	result := make(StrMap)
	for key, value := range originData {
		var tmpKey string
		if rootKey == "" {
			tmpKey = key
		} else {
			tmpKey = rootKey + delimiter + key
		}
		if reflect.ValueOf(value).Kind() == reflect.Map {
			v := AnyToStrMap(value)
			tmpMap := FlattenMap(tmpKey, delimiter, v)
			for k, v := range tmpMap {
				result[k] = v
			}
		} else {
			result[tmpKey] = value
		}
	}
	return result
}

func CanConvertToFloat32Loselessly(v float64) bool {
	absV := math.Abs(v)
	if absV < math.MaxFloat32 && absV > math.SmallestNonzeroFloat32 {
		return true
	}
	return false
}

func CanConvertToInt64Loselessly(v float64) bool {
	return v == math.Trunc(v)
}

func CanConvertToInt32Loselessly(v float64) bool {
	return v == math.Trunc(v) && v < math.MaxInt32 && v > math.MinInt32
}

// StringToChunks split a string into string slices with element's size <= chunkSize
// Examples:
// StringToChunks("abcd", 1) => []string{"a", "b", "c", "d"}
// StringToChunks("abcd", 2) => []string{"ab", "cd"}
// StringToChunks("abcd", 3) => []string{"abc", "d"}
// stringToChunks("abcd", 4) => []string{"abcd"}
// stringToChunks("abcd", 5) => []string{"abcd"}
func StringToChunks(s string, chunkSize int) []string {
	var chunks []string
	strLength := len(s)
	index := 0
	for index < strLength {
		endIndex := IntMin(index+chunkSize, strLength)
		chunk := s[index:endIndex]
		chunks = append(chunks, chunk)
		index = endIndex
	}
	return chunks
}

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
