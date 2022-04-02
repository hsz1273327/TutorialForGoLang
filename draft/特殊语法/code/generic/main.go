package main

import (
	"fmt"
	"sync"
)

// helloworld
func Map[T1, T2 any](iter []T1, fn func(T1) T2) []T2 {
	wg := sync.WaitGroup{}
	result := make([]T2, len(iter))
	for i, e := range iter {
		wg.Add(1)
		go func(e T1, i int) {
			defer wg.Done()
			r := fn(e)
			result[i] = r
		}(e, i)
	}
	wg.Wait()
	return result
}

func Reduce[T any](iter []T, fn func(T, T) T) T {
	var first, last T
	for i, e := range iter {
		switch i {
		case 0:
			{
				first = e
			}
		default:
			{
				last = e
				result := fn(first, last)
				first = result
			}
		}
	}
	return first
}

// 泛型函数
type RealNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type ComplexNumber interface {
	~complex64 | ~complex128
}

type Number interface {
	RealNumber | ComplexNumber
}

func Sum[T Number](x ...T) T {
	return Reduce(x, func(first, last T) T { return first + last })
}

// 泛型自定义类型
type MyClass[T ComplexNumber] struct {
	A T
	B string
}

func NewMyClass[T1 ComplexNumber](a T1, b string) *MyClass[T1] {
	s := new(MyClass[T1])
	s.A = a
	s.B = b
	return s
}

func (s *MyClass[T]) ToString() string {
	return fmt.Sprintf("A is %f, B is %s", s.A, s.B)
}

// 泛型函数签名
type GFn[T RealNumber] func(x T) T

func (fn GFn[T]) Echo(x T) string {
	return fmt.Sprintf("echo %v", fn(x))
}

func Callback[T1 RealNumber](x T1, fn GFn[T1]) T1 {
	return fn(x)
}

// 泛型内置结构派生
type NumberSlice[T Number] []T

func (s NumberSlice[T]) Map(fn func(T) any) []any {
	return Map(s, fn)
}

// 泛型约束
type Callable[T RealNumber] interface {
	RealNumber
	A(T) error
}

type ACall int

func (a ACall) A(x int) error {
	fmt.Println("value is ", a)
	fmt.Println("get ", x)
	return nil
}

func CallCallable[T Callable[int]](a T) {
	a.A(8)
}

func main() {
	fmt.Println(Reduce(Map([]int{1, 2, 3, 4, 5, 6}, func(x int) int { return x + 10 }), func(x, y int) int { return x + y }))

	fmt.Println(Sum(Map([]complex64{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i, 5 + 5i, 6 + 6i}, func(x complex64) complex64 { return x + 10 })...))

	c := NewMyClass(1+10i, "test")
	fmt.Println(c.ToString())

	f := GFn[int](func(x int) int { return x * 2 })
	fmt.Println(f.Echo(12))
	fmt.Println(Callback(123, f))

	fmt.Println(Callback(1234.2, func(x float32) float32 { return x * 2 }))

	x := NumberSlice[int64]([]int64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1})
	fmt.Println(x.Map(func(x int64) any { return x*2 + 3 }))
	CallCallable(ACall(12))
}
