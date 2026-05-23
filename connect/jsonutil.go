package connect

import (
	"strconv"

	"github.com/valyala/fastjson"
)

func jsonInt(v *fastjson.Value) int {
	if v == nil {
		return 0
	}
	switch v.Type() {
	case fastjson.TypeNumber:
		n, _ := v.Int64()
		return int(n)
	case fastjson.TypeString:
		n, err := strconv.ParseInt(string(v.GetStringBytes()), 10, 64)
		if err != nil {
			return 0
		}
		return int(n)
	default:
		return 0
	}
}

func jsonInt64(v *fastjson.Value) int64 {
	if v == nil {
		return 0
	}
	switch v.Type() {
	case fastjson.TypeNumber:
		n, _ := v.Int64()
		return n
	case fastjson.TypeString:
		n, err := strconv.ParseInt(string(v.GetStringBytes()), 10, 64)
		if err != nil {
			return 0
		}
		return n
	default:
		return 0
	}
}
