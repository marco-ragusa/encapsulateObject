package main

import (
	"fmt"
	"reflect"
)

// Invoke - []results, err := invoke(AnyStructInterface, MethodName, Params...)
func Invoke(any interface{}, name string, args ...interface{}) ([]reflect.Value, error) {
	method := reflect.ValueOf(any).MethodByName(name)
	methodType := method.Type()
	numIn := methodType.NumIn()

	// check the number of args are minus of parameters
	if len(args) < numIn {
		return nil, fmt.Errorf("method %s must have minimum %d params. Have %d", name, numIn, len(args))
	}

	// check the method is variadic if the args are greater of parameters
	if  len(args) > numIn && !methodType.IsVariadic() {
		return nil, fmt.Errorf("method %s must have %d params. Have %d", name, numIn, len(args))
	}

	// build the parameters of the method
	in := make([]reflect.Value, len(args))
	for i := 0; i < len(args); i++ {
		// get method parameter type for implicit conversion
		var inType reflect.Type
		// for variadic (it's always in last position) get type of the last parameter
		if i >= numIn-1 && methodType.IsVariadic()  {
			inType = methodType.In(numIn - 1).Elem()
		} else {
			inType = methodType.In(i)
		}
		// check arg value is valid
		argValue := reflect.ValueOf(args[i])
		if !argValue.IsValid() {
			return nil, fmt.Errorf("method %s. Param[%d] must be %s. Have %s", name, i, inType, argValue.String())
		}
		// convert the arg to method parameter type
		argType := argValue.Type()
		if argType.ConvertibleTo(inType) {
			in[i] = argValue.Convert(inType)
		} else {
			return nil, fmt.Errorf("method %s. Param[%d] must be %s. Have %s", name, i, inType, argType)
		}
	}

	// return an array of return's method
	return method.Call(in), nil
}


// Triangle object structure
type Triangle struct {
	Height float64
	Width  float64
}

// NewTriangle constructor - return *Triangle object
func NewTriangle(h float64, w float64) *Triangle {
	t := &Triangle{}
	t.Height = h
	t.Width = w
	return t
}

// TriangleArea calculate the area of triangle
func (t *Triangle) TriangleArea(variadic ...string) float64 {
	// use variadic param only for test
	for _, s := range variadic {
		fmt.Printf("%s ",s)
	}
	fmt.Println()

	return (t.Height*t.Width)/2
}

// Container object structure
type Container struct {
	generic interface{}
}

// NewContainer constructor - return *Triangle object
func NewContainer(generic interface{}) *Container {
	c := &Container{}
	c.generic = generic
	return c
}

// ContainerArea call method of Triangle and return the result
func (c *Container) ContainerArea() float64 {
	r, _ := Invoke(c.generic,"TriangleArea","test1","test2","test3")
	return r[0].Float()
}

func main(){
	// create a triangle
	t := NewTriangle(5,5)
	// encapsulate object
	c := NewContainer(t)
	// cal method of encapsulated object through container method
	fmt.Println(c.ContainerArea())
}