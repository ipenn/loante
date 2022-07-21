package test

import (
	"fmt"
	"loante/tools"
	"testing"
)

func TestStructToMapReflect(t *testing.T)  {
	type Student struct{
	   Name string `json:"name"`
	   Age int `json:"age"`
	}

	stuObj := Student{
		Name: "wanghw",
		Age: 22,
	}

	m2, _ := tools.StructToMapReflect(&stuObj,"json")
	for key, val := range m2{
		fmt.Printf("key: %s, val: %v, typeVal: %T \n", key, val, val)
	}
}
