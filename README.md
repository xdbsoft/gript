# gript - An expression evaluator in go

[![Godoc](https://godoc.org/github.com/xdbsoft/gript?status.png)](https://godoc.org/github.com/xdbsoft/gript)
[![Build Status](https://travis-ci.org/xdbsoft/gript.svg?branch=master)](https://travis-ci.org/xdbsoft/gript)
[![Coverage](http://gocover.io/_badge/github.com/xdbsoft/gript)](http://gocover.io/_badge/github.com/xdbsoft/gript)
[![Report](https://goreportcard.com/badge/github.com/xdbsoft/gript)](https://goreportcard.com/report/github.com/xdbsoft/gript)

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