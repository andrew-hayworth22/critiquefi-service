package tools

import "strings"

// SetBuilder constructs an SQL SET for a PATCH request
type SetBuilder struct {
	cols []string
	args []any
}

// NewSetBuilder constructs a SetBuilder
func NewSetBuilder() *SetBuilder {
	return &SetBuilder{}
}

// Set adds a column assignment to the builder
func (b *SetBuilder) Set(col string, arg any) {
	b.cols = append(b.cols, col)
	b.args = append(b.args, arg)
}

// SetIf adds a column assignment if the condition is true
func (b *SetBuilder) SetIf(condition bool, col string, arg any) {
	if condition {
		b.Set(col, arg)
	}
}

// Empty checks if there are no column assignments in the set builder
func (b *SetBuilder) Empty() bool {
	return len(b.cols) == 0
}

// Args returns the arguments
func (b *SetBuilder) Args() []any {
	return b.args
}

// BuildSet constructs the SET clause (ex: "col=$1, col2=$2")
func (b *SetBuilder) BuildSet() string {
	if b.Empty() {
		return ""
	}

	out := make([]string, len(b.cols))
	for i, col := range b.cols {
		out[i] = col + "=$" + string(rune(i+1))
	}
	return strings.Join(out, ", ")
}
