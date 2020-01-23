package generator

import (
	"../ast"
	"fmt"
	"llvm.org/llvm/bindings/go/llvm"
)

type CodeGen struct {
	function interface{}  // 現在コード生成中のFunction
	mod      *llvm.Module // 生成したModuleを格納
	builder  llvm.Builder // LLVM-IRを生成するIRBuilderクラス
}

func New() *CodeGen {
	cg := &CodeGen{}
	cg.builder = llvm.NewBuilder()
	return cg
}

func (cg *CodeGen) getModule() llvm.Module {
	if cg.mod != nil {
		return *cg.mod
	}
	return llvm.NewModule("null")
}

// generateTranslationUnit - モジュール生成メソッド
func (cg *CodeGen) generateTranslationUnit(tu *ast.TranslationUnit, name string) bool {
	mod := llvm.NewModule(name)

	// function declaration
	for i := 0; ; i++ {
		if i+1 >= len(tu.Prototypes) {
			break
		}
		proto := tu.Prototypes[i]
	}

	// function definition
	for i := 0; ; i++ {
		if i+1 >= len(tu.Functions) {
			break
		}
		function := tu.Functions[i]
	}

	return true
}

func (cg *CodeGen) generatePrototype(prototype *ast.Prototype, mod *llvm.Module) llvm.Value {

	// 既に定義済みか
	{
		function := mod.NamedFunction(prototype.GetName())
		if !function.IsNull() {
			if function.ParamsCount() == len(prototype.Parameters) && function.IsNil() {
				return function
			}
			msg := fmt.Sprintf("error::function %s is redefined", prototype.GetName())
			panic(msg)
		}
	}

	// まだ未定義のとき

	// create arg_types
	intTypes := []llvm.Type{}
	for i := 0; i < len(prototype.Parameters); i++ {
		intTypes[i] = llvm.Int32Type()
	}

	// create func type
	functionType := llvm.FunctionType(llvm.Int32Type(), intTypes, false)

	// create function
	function := llvm.AddFunction(*mod, prototype.GetName(), functionType)
	for i, _ := range function.Params() {
		function.Params()[i].SetName(prototype.Parameters[i].Name() + "Arg")
	}

	return function
}

func (cg *CodeGen) generateFunctionDefinition() {

}
