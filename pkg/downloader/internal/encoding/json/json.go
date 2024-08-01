package json

import jsoniter "github.com/json-iterator/go"

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func Stringify(a any) string {
	raw, err := JSON.MarshalToString(a)
	if err != nil {
		return ""
	}
	return raw
}

func Prettify(a any) string {
	bs, err := JSON.MarshalIndent(a, "", "    ")
	if err != nil {
		return ""
	}
	return string(bs)
}
