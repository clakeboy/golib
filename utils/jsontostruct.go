package utils

import "fmt"

// JSON 生成 struct 结构内容
func GenerateStruct(key string, data M, indent string) {
	if indent == "" {
		fmt.Printf("%stype %s struct {\n", indent, UcFirst(key))
	} else {
		fmt.Print("struct {\n")
	}

	for k, v := range data {
		generateType(k, v, indent+"    ")
	}
	fmt.Print(indent, "}")
	if indent == "" {
		print("\n")
	}
}

func generateType(key string, obj interface{}, indent string) {
	print(indent, UcFirst(key), " ")
	switch obj.(type) {
	case M:
		GenerateStruct(key, obj.(M), indent+"    ")
		fmt.Printf(" `json:\"%s,omitempty\"`\n", key)
	case map[string]interface{}:
		GenerateStruct(key, obj.(map[string]interface{}), indent+"    ")
		fmt.Printf(" `json:\"%s,omitempty\"`\n", key)
	case string:
		fmt.Printf("string `json:\"%s,omitempty\"`\n", key)
	case float64:
		fmt.Printf("float64 `json:\"%s,omitempty\"`\n", key)
	case int64:
		fmt.Printf("int64 `json:\"%s,omitempty\"`\n", key)
	case bool:
		fmt.Printf("bool `json:\"%s,omitempty\"`\n", key)
	case []interface{}:
		// println("[]interface{}")
		switch obj.([]interface{})[0].(type) {
		case string:
			println("[]string")
		case float64:
			println("[]float64")
		case int64:
			println("[]int64")
		case map[string]interface{}:
			print("[]")
			GenerateStruct(key, obj.([]interface{})[0].(map[string]interface{}), indent+"    ")
			fmt.Printf(" `json:\"%s,omitempty\"`\n", key)
		}
	default:
		print("string")
		fmt.Printf(" `json:\"%s,omitempty\"`\n", key)
	}
}
