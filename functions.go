package main

import (
	"fmt"
	"strings"
)

// A Function definition.
type Function struct {
	Name       string // C name of the function
	GoName     string // Go name of the function with the API prefix stripped
	Parameters []Parameter
	Return     Type
	Overloads  []Overload
}

// An Overload describes an alternative signature for the same function.
type Overload struct {
	GoName       string // Go name of the original function
	OverloadName string // Go name of the overload
	Parameters   []Parameter
	Return       Type
}

func (o Overload) function() Function {
	return Function{
		GoName:     o.GoName,
		Parameters: o.Parameters,
		Return:     o.Return,
	}
}

// IsImplementedForSyscall reports whether the function is implemented for syscall or not.
func (o Overload) IsImplementedForSyscall() bool {
	return o.function().IsImplementedForSyscall()
}

// Syscall returns a syscall expression for Windows.
func (o Overload) Syscall() string {
	return o.function().Syscall()
}

// IsImplementedForSyscall reports whether the function is implemented for syscall or not.
func (f Function) IsImplementedForSyscall() bool {
	return len(f.Parameters) <= 18
}

// Syscall returns a syscall expression for Windows.
func (f Function) Syscall() string {
	var ps []string
	for _, p := range f.Parameters {
		ps = append(ps, p.Type.ConvertGoToUintptr(p.GoName()))
	}
	for len(ps) == 0 || len(ps)%3 != 0 {
		ps = append(ps, "0")
	}

	post := ""
	if len(ps) > 3 {
		post = fmt.Sprintf("%d", len(ps))
	}

	return fmt.Sprintf("syscall.Syscall%s(gp%s, %d, %s)", post, f.GoName, len(f.Parameters), strings.Join(ps, ", "))
}

// A Parameter to a Function.
type Parameter struct {
	Name string
	Type Type
}

// CName returns a C-safe parameter name.
func (p Parameter) CName() string {
	return renameIfReservedCWord(p.Name)
}

// GoName returns a Go-safe parameter name.
func (p Parameter) GoName() string {
	return renameIfReservedGoWord(p.Name)
}

func renameIfReservedCWord(word string) string {
	switch word {
	case "near", "far":
		return fmt.Sprintf("x%s", word)
	}
	return word
}

func renameIfReservedGoWord(word string) string {
	switch word {
	case "func", "type", "struct", "range", "map", "string":
		return fmt.Sprintf("x%s", word)
	}
	return word
}
