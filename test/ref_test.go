package test

import (
	"fmt"
	"reflect"
	"testing"
)

const tagName = "Testing"

type Info struct {
	Name string `Testing:"-" bun:"aa"`
	Age  int    `Testing:"age,min=17,max=60" bun:"da"`
	Sex  string `Testing:"sex,required" bun:"ca"`
}

func TestRef(t *testing.T)  {
	info := Info{
		Name: "benben",
		Age:  23,
		Sex:  "male",
	}

	t1 := reflect.TypeOf(info)
	fmt.Println("Type:", t1.Name())
	fmt.Println("Kind:", t1.Kind())

	for i := 0; i < t1.NumField(); i++ {
		field := t1.Field(i) //获取结构体的每一个字段
		tag := field.Tag.Get(tagName)
		fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
	}
}
