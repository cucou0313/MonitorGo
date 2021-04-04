/*
Project: Monitor
Author: Guo Kaikuo
Create time: 2021-04-01 16:09
IDE: GoLand
*/

package utils

import "github.com/tidwall/gjson"

/**
 * @Description:使用gjson获取string类型，不存在则返回默认值
 * @param str json string
 * @param key  key
 * @param default_res   不存在或值为空时返回默认值
 * @return string
 */
func JsonGetString(str *string, key string, default_res string) string {
	res := gjson.Get(*str, key)
	if res.Exists() {
		if res.String() == "" {
			return default_res
		}
		return res.String()
	}
	return default_res
}

/**
 * @Description:使用gjson获取int类型，不存在则返回默认值
 * @param str json int
 * @param key  key
 * @param default_res   不存在或值为空时返回默认值
 * @return int64
 */
func JsonGetInt(str *string, key string, default_res int64) int64 {
	res := gjson.Get(*str, key)
	if res.Exists() {
		return res.Int()
	}
	return default_res
}

func JsonGetBool(str *string, key string, default_res bool) bool {
	res := gjson.Get(*str, key)
	if res.Exists() {
		return res.Bool()
	}
	return default_res
}

func JsonGetArray(str *string, key string, default_res []gjson.Result) []gjson.Result {
	res := gjson.Get(*str, key)
	if res.Exists() {
		return res.Array()
	}
	return default_res
}
