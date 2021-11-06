package adb

import "strings"

func GetProp(device string, key string) string {
	r := ""
	callback := func(out string) {
		v := strings.Split(out, ":")
		for i, p := range v {
			v[i] = strings.Trim(p, "[]")
			v[i] = strings.TrimSpace(v[i])
			r = v[i]
		}
	}
	Run(device, callback, "getprop "+key)
	return r
}
