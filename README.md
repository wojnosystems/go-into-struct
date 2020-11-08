# Into Struct

This is a parsing library that abstracts away walking through a struct and its values so that implementations of the Parser interface can focus on extracting each field into the struct.

# Features

* Destination checking: automatically checks the destination struct and fields to see if they're settable and returns ErrProgramming if they're not.
* Handles fall-backs. If your parser doesn't support a struct, it will fall-back to evaluating the structs individual fields. This allows you to support complex types.

# Example usage

See [github.com/wojnosystems/go-env](https://github.com/wojnosystems/go-env). This is an environment parsing library. It takes in the operating system environment variables and maps them to structs using rules defined by that library. The go-into-struct walks the structure, the go-env tells go-into-struct how to put items into the structure.
