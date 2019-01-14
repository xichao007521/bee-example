package controllers

import (
	"do-global.com/bee-example/error"
	"regexp"
	"strconv"
)

// 必须为一个数字，其值必须小于等于指定的最小值
func (t *BasicController) GetNumCheckMax(key string, max int) int {
	v := t.GetInt(key)
	if v > max {
		panic(myError.NewBizError(400, "param:"+key+" 的值必须小于等于 "+strconv.Itoa(max)))
	}
	return v
}

// 必须为一个数字，其值必须大于等于指定的最小值
func (t *BasicController) GetNumCheckMin(key string, min int) int {
	v := t.GetInt(key)
	if v < min {
		panic(myError.NewBizError(400, "param:"+key+" 的值必须大于等于 "+strconv.Itoa(min)))
	}
	return v
}

// 元素的个数必须在指定的范围
func (t *BasicController) GetStringsCheckSize(key string, min int, max int) []string {
	v := t.GetStringsNE(key)
	if len(v) < min || len(v) > max {
		panic(myError.NewBizError(400, "param: "+key+" 的值的个数必须大于等于 "+strconv.Itoa(min)+"且小于等于"+strconv.Itoa(max)))
	}
	return v
}

// 校验字符串长度,校验成功返回字符串
func (t *BasicController) GetStringCheckLength(key string, min int, max int) string {
	v := t.GetStringNE(key)
	if len(v) < min || len(v) > max {
		panic(myError.NewBizError(400, "param:"+key+" length must be in "+strconv.Itoa(min)+" - "+strconv.Itoa(max)))
	}
	return v
}

// 元素必须是电子邮箱地址
func (t *BasicController) GetEmail(key string) string {
	if m, _ := regexp.MatchString("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+", t.GetStringNE(key)); !m {
		panic(myError.NewBizError(400, "param:"+t.GetString(key)+" is not a email"))
	}
	return t.GetString(key)
}
