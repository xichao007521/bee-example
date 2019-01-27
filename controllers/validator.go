package controllers

import (
	"do-global.com/bee-example/error"
	"regexp"
	"strconv"
)

// 获取字符串【非空】
func (t *BasicController) GetStringNE(key string) string {
	v := t.GetString(key)
	if v == "" {
		panic(myError.NewBizError(400, "param:"+key+" must not be empty"))
	}
	return v
}

// 获取字符串数组【非空】
func (t *BasicController) GetStringsNE(key string) []string {
	v := t.GetStrings(key)
	if len(v) == 0 {
		panic(myError.NewBizError(400, "param:"+key+" must not be empty array"))
	}
	return v
}

// 获取 int 类型值【非空】
func (t *BasicController) GetIntNE(key string) int {
	v, err := strconv.Atoi(t.GetStringNE(key))
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a int value"))
	}
	return v
}

// 获取 int8 类型值【非空】
func (t *BasicController) GetInt8NE(key string) int8 {
	i64, err := strconv.ParseInt(t.GetStringNE(key), 10, 8)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a int8 value"))
	}
	return int8(i64)
}

// 获取 uint8 类型值【非空】
func (t *BasicController) GetUint8NE(key string) uint8 {
	u64, err := strconv.ParseUint(t.GetStringNE(key), 10, 8)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a uint8 value"))
	}
	return uint8(u64)
}

// 获取 int16 类型值【非空】
func (t *BasicController) GetInt16NE(key string) int16 {
	i64, err := strconv.ParseInt(t.GetStringNE(key), 10, 16)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a int16 value"))
	}
	return int16(i64)
}

// 获取 uint16 类型值【非空】
func (t *BasicController) GetUint16NE(key string) uint16 {
	u64, err := strconv.ParseUint(t.GetStringNE(key), 10, 16)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a uint16 value"))
	}
	return uint16(u64)
}

// 获取 int32 类型值【非空】
func (t *BasicController) GetInt32NE(key string) int32 {
	i64, err := strconv.ParseInt(t.GetStringNE(key), 10, 32)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a int32 value"))
	}
	return int32(i64)
}

// 获取 uint32 类型值【非空】
func (t *BasicController) GetUint32NE(key string) uint32 {
	u64, err := strconv.ParseUint(t.GetStringNE(key), 10, 32)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a uint32 value"))
	}
	return uint32(u64)
}

// 获取 int64 类型值【非空】
func (t *BasicController) GetInt64NE(key string) int64 {
	i64, err := strconv.ParseInt(t.GetStringNE(key), 10, 64)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a int64 value"))
	}
	return int64(i64)
}

// 获取 uint64 类型值【非空】
func (t *BasicController) GetUint64NE(key string) uint64 {
	u64, err := strconv.ParseUint(t.GetStringNE(key), 10, 64)
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a uint64 value"))
	}
	return uint64(u64)
}

// 获取 bool 类型值【非空】
func (t *BasicController) GetBoolNE(key string) bool {
	v, err := strconv.ParseBool(t.GetStringNE(key))
	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a bool value"))
	}
	return v
}

// 获取 float64 类型值【非空】
func (t *BasicController) GetFloatNE(key string) float64 {
	v, err := strconv.ParseFloat(t.GetStringNE(key), 64)

	if err != nil {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a float value"))
	}
	return v
}

// 获取一个 int 类型的数字 ，其值必须小于等于指定的最小值 【非空】
func (t *BasicController) GetIntNECheckMax(key string, max int) int {
	v := t.GetIntNE(key)
	if v > max {
		panic(myError.NewBizError(400, "param:"+key+" 的值必须小于等于 "+strconv.Itoa(max)))
	}
	return v
}

// 获取一个 int 类型的数字，其值必须大于等于指定的最小值【非空】
func (t *BasicController) GetIntNECheckMin(key string, min int) int {
	v := t.GetIntNE(key)
	if v < min {
		panic(myError.NewBizError(400, "param:"+key+" 的值必须大于等于 "+strconv.Itoa(min)))
	}
	return v
}

// 获取一个 string 类型的数组 元素的个数必须在指定的范围【非空】
func (t *BasicController) GetStringsNECheckSize(key string, min int, max int) []string {
	v := t.GetStringsNE(key)
	if len(v) < min || len(v) > max {
		panic(myError.NewBizError(400, "param: "+key+" 的值的个数必须大于等于 "+strconv.Itoa(min)+"且小于等于"+strconv.Itoa(max)))
	}
	return v
}

// 获取一个string字符串，字符串的长度必须在指定的范围内【非空】
func (t *BasicController) GetStringNECheckLength(key string, min int, max int) string {
	v := t.GetStringNE(key)
	if len(v) < min || len(v) > max {
		panic(myError.NewBizError(400, "param:"+key+" length must be in "+strconv.Itoa(min)+" - "+strconv.Itoa(max)))
	}
	return v
}

// 获取一个string字符串，字符串必须是电子邮箱地址【非空】
func (t *BasicController) GetStringNECheckEmail(key string) string {
	if m, _ := regexp.MatchString("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+", t.GetStringNE(key)); !m {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a email"))
	}
	return t.GetString(key)
}
