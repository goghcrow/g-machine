package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestCompile(t *testing.T) {
	fmt.Println(strings.Join(SliceMap(Parse(preludeDefs), func(t Term) string { return t.String() }), "\n"))
}
