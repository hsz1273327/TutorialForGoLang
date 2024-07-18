package main

import (
	"fmt"
	"reflect"
	"unicode"
)

type Btype struct {
	C int    `mytags:"c"`
	d string `mytags:"d"`
}

func (a *Btype) CallB(name string) {
	fmt.Println("called", name)
}

type Atype struct {
	A int    `mytags:"a"`
	b string `mytags:"b"`
	*Btype
}

func (a *Atype) CallA(name string) {
	fmt.Println("called", name)
}
func main() {
	test := Atype{
		A: 10,
		b: "abc",
		Btype: &Btype{
			C: 11,
			d: "def",
		},
	}
	vp := reflect.ValueOf(&test)
	// 首先获得的是指针类型
	tpk := vp.Kind()
	fmt.Println(tpk) //ptr
	// 经过转换获得的类型为结构体
	v := vp.Elem()
	tk := v.Kind()
	fmt.Println(tk) //struct
	// 获取指针对应的结构体类型名
	t := v.Type()
	fmt.Println(t.Name()) //Atype
	// 查看结构体的字段结构
	fieldmap := map[string]reflect.StructField{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Println("---field---")
		fmt.Println(field.Name)      // A b Btype
		fmt.Println(field.Type)      //int string *main.Btype
		fmt.Println(field.Tag)       //mytags:"a"  mytags:"b"
		fmt.Println(field.Anonymous) //false false true
		fmt.Println(field.Index)     //[0] [1] [2]
		fmt.Println(field.Offset)    //0 8 24
		fieldmap[field.Name] = field
	}
	// 查看结构体的方法集合
	methodmap := map[string]reflect.Method{}
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Println("---method---")
		fmt.Println(method.Name)  //CallB
		fmt.Println(method.Type)  //func(main.Atype, string)
		fmt.Println(method.Index) //0
		methodmap[method.Name] = method
	}

	// 创建结构体实例
	//创建指定类型的零值实例的指针的Value
	fmt.Println("---new p instance---")
	newvp := reflect.New(t)
	fmt.Println(newvp.Interface().(*Atype))

	fmt.Println("---new instance---")
	newv := reflect.Zero(t)
	fmt.Println(newv.Interface().(Atype))

	// 查看结构体实例的字段并重置字段为0值
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		fmt.Println("---field---")
		fmt.Println(field.Name)      // A b Btype
		fmt.Println(field.Type)      //int string *main.Btype
		fmt.Println(field.Tag)       //mytags:"a"  mytags:"b"
		fmt.Println(field.Anonymous) //false false true
		fmt.Println(field.Index)     //[0] [1] [2]
		fmt.Println(field.Offset)    //0 8 24
		// 首字母小写表示是私有字段,Value只能针对公开字段获取和赋值
		if unicode.IsUpper([]rune(field.Name)[0]) {
			fmt.Println("---field value---")
			fmt.Println(fieldValue.Interface())
			fmt.Println("---field value reset to zero---")
			zeroValue := reflect.Zero(fieldValue.Type())
			fieldValue.Set(zeroValue)
		}
	}
	fmt.Println(test)

	// 调用实例的方法
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		methodValue := v.Method(i)
		fmt.Println("---method---")
		fmt.Println(method.Name)  //CallB
		fmt.Println(method.Type)  //func(main.Atype, string)
		fmt.Println(method.Index) //0
		fmt.Println("---call method---")
		methodValue.Call([]reflect.Value{
			reflect.ValueOf("hello"),
		})
	}
}
