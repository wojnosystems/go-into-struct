# Into Struct

This is a parsing library that abstracts away walking through a struct and its values so that implementations of the Parser interface can focus on extracting each field into the struct.

# Features

* Destination checking: automatically checks the destination struct and fields to see if they're settable and returns ErrProgramming if they're not.
* Handles fall-backs. If your parser doesn't support a struct, it will fall-back to evaluating the structs individual fields. This allows you to support complex types.

# Example usage

See [github.com/wojnosystems/go-env](https://github.com/wojnosystems/go-env). This is an environment parsing library. It takes in the operating system environment variables and maps them to structs using rules defined by that library. The go-into-struct walks the structure, the go-env tells go-into-struct how to put items into the structure.

You need only implement the methods in Parser for your class and pass it to the Unmarshal method. When your methods are called, you'll be provided a structPath for each field and slice. This allows you to navigate and translate the struct path in Go to whatever your data source is. It's JSON-path-like. Each field is named after its Go Field name, which is whatever you named it in code. For example:

```go
package main
import into_struct "github.com/wojnosystems/go-into-struct"
type Struct2 struct {
  Struct2Value string
}
type Top struct {
  TopValue Struct2
  Slices   Struct2
}
func main() {
  var myTop Top
  into_struct.Unmarshall(&myTop, myParser{})
}
```

myParser's SetValue will receive the following structPaths:

* TopValue
* TopValue.Struct2Value
* Slices
* Slices[0].Struct2Value
* Slices[1].Struct2Value
* Slices[n].Struct2Value

where n is the number you returned from SliceLen.

# FAQ

## No map support?

Nope. Not out of the box. I build this library for type-checked parsing of values. You're free to implement a map-maker in your parser, however. This library simply provides a convenience for every settable value in the struct as well as slices.
