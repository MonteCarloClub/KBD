package util

import jsoniter "github.com/json-iterator/go"

func ToString(v interface{}) string {
	str, err := jsoniter.MarshalToString(v)
	if err != nil {
		return ""
	}
	return str
}
