# gript - An expression evaluator in go

[![Godoc](https://godoc.org/github.com/xdbsoft/gript?status.png)](https://godoc.org/github.com/xdbsoft/gript)

## How-to

	package main

	import (
        "fmt"
		"github.com/xdbsoft/gript"
	)
		
	func main() {
        result, err := Eval(" abc > 3+1   ||	(ab < 4-2 && ab > 6%2) || d < 0", map[string]interface{}{"abc": 1, "d": 1})

        //result will contain the boolean true

        ...

	}