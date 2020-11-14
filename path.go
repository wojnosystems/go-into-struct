package into_struct

import (
	"strings"
)

const (
	pathPartSeparator = "."
)

type Path struct {
	parts []PathParter
}

// Parts is all of the components of the path. Do not edit these values
func (p Path) Parts() (parts []PathParter) {
	return p.parts
}

// String converts the path into a canoncial string representation with dots (.)
// separating components and square brackets ([]) indicating slice components.
func (p Path) String() string {
	stringParts := make([]string, len(p.parts))
	for i := range stringParts {
		stringParts[i] = p.parts[i].String()
	}
	return strings.Join(stringParts, pathPartSeparator)
}

// push a component to the end/top of the path
func (p *Path) push(part PathParter) {
	p.parts = append(p.parts, part)
}

// pop the end/top most element of the path
func (p *Path) pop() {
	if len(p.parts) > 0 {
		p.parts = p.parts[0 : len(p.parts)-1]
	}
}

// within performs an push and a pop
func (p *Path) within(part PathParter, callback func()) {
	p.push(part)
	callback()
	p.pop()
}

// Top gets the element most recently pushed to the stack or nil if empty
func (p Path) Top() PathParter {
	if len(p.parts) == 0 {
		return nil
	}
	return p.parts[len(p.parts)-1]
}

// ParentOfTop gets the element pushed just prior to the most recently pushed component
// to the stack or nil if one or no elements are in the stack
func (p Path) ParentOfTop() PathParter {
	if len(p.parts) < 2 {
		return nil
	}
	return p.parts[len(p.parts)-2]
}

// setTop replaces the top-most element with as.
func (p *Path) setTop(as PathParter) {
	p.parts[len(p.parts)-1] = as
}
