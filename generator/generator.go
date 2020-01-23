package generator

import (
	"llvm.org/llvm/bindings/go/llvm"
)

type CodeGen struct {
	builder llvm.Builder
}

func New() *CodeGen {
	cg := &CodeGen{}
	cg.builder = llvm.NewBuilder()
	return cg
}
