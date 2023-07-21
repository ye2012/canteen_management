package utils

import (
	"reflect"
)

func GetSpecifiedFieldsValueWithSpecialField(v interface{}, specialTag string, specifiedTags ...string) (interface{}, []interface{}) {
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	refType := r.Type()
	ret := make([]interface{}, 0)
	var specialValue interface{}

loopI:
	for i := 0; i < r.NumField(); i++ {
		tag := refType.Field(i).Tag.Get("json")
		if tag == specialTag {
			specialValue = r.Field(i).Interface()
		}
		for _, t := range specifiedTags {
			if t == tag {
				value := r.Field(i).Interface()
				ret = append(ret, value)
				continue loopI
			}
		}
	}
	return specialValue, ret
}

func GetFieldsValue(v interface{}, skipTags ...string) []interface{} {
	skipTags = append(skipTags, "-")
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	refType := r.Type()
	ret := make([]interface{}, 0)

loopI:
	for i := 0; i < r.NumField(); i++ {
		tag := refType.Field(i).Tag.Get("json")
		for _, t := range skipTags {
			if t == tag {
				continue loopI
			}
		}

		value := r.Field(i).Interface()
		ret = append(ret, value)
	}
	return ret
}

func GetSpecifiedFieldsTag(v interface{}, key string, specifiedTags ...string) []string {
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	ret := make([]string, 0)
	refType := r.Type()

loopI:
	for i := 0; i < r.NumField(); i++ {
		tag := refType.Field(i).Tag.Get(key)
		for _, t := range specifiedTags {
			if t == tag {
				ret = append(ret, tag)
				continue loopI
			}
		}
	}
	return ret
}

func GetFieldsTagByKey(v interface{}, key string, skipTags ...string) []string {
	skipTags = append(skipTags, "-")
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	ret := make([]string, 0)
	refType := r.Type()

loopI:
	for i := 0; i < r.NumField(); i++ {
		tag := refType.Field(i).Tag.Get(key)
		for _, t := range skipTags {
			if t == tag {
				continue loopI
			}
		}

		ret = append(ret, tag)
	}
	return ret
}

// 传入对象指针
func GetFieldsAddr(v interface{}, skipTags ...string) []interface{} {
	skipTags = append(skipTags, "-")
	r := reflect.ValueOf(v).Elem()
	ret := make([]interface{}, 0)
	refType := r.Type()

loopI:
	for i := 0; i < r.NumField(); i++ {
		tag := refType.Field(i).Tag.Get("json")
		for _, t := range skipTags {
			if t == tag {
				continue loopI
			}
		}

		addr := r.Field(i).Addr().Interface()
		ret = append(ret, addr)
	}
	return ret
}

func GetFieldsTagValueMap(v interface{}) map[string]interface{} {
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	refType := r.Type()
	ret := make(map[string]interface{})
	for i := 0; i < r.NumField(); i++ {
		tag := refType.Field(i).Tag.Get("json")
		if tag == "-" {
			continue
		}
		value := r.Field(i).Interface()
		ret[tag] = value
	}
	return ret
}

func GetMapValue(paramsMap map[string]string, key, defaultValue string) string {
	if value, ok := paramsMap[key]; ok {
		return value
	}
	return defaultValue
}
