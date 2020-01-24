package generator

import (
	"../ast"
	"fmt"
	"llvm.org/llvm/bindings/go/llvm"
)

type CodeGen struct {
	curFunc *llvm.Value  // 現在コード生成中のFunction
	mod     *llvm.Module // 生成したModuleを格納
	builder llvm.Builder // LLVM-IRを生成するIRBuilderクラス
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
	for _, proto := range tu.Prototypes {
		cg.generatePrototype(&proto, &mod)
	}

	// function definition
	for _, function := range tu.Functions {
		cg.generateFunctionDefinition(&function, &mod)
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
	for i := range function.Params() {
		function.Params()[i].SetName(prototype.Parameters[i].Name() + "Arg")
	}

	return function
}

func (cg *CodeGen) generateFunctionDefinition(functionLiteral *ast.FunctionLiteral, mod *llvm.Module) llvm.Value {
	function := cg.generatePrototype(&functionLiteral.Prototype, mod)
	cg.curFunc = &function

	bblock := llvm.AddBasicBlock(function, "entry")
	cg.builder.SetInsertPoint(bblock, bblock.FirstInstruction())

	// TODO: Functionのボディを生成

	return function
}

func (cg *CodeGen) generateFunctionStatement(functionStmt *ast.FunctionStatement) *llvm.Value {

	// insert variable declarations
	var vdecl *ast.DeclarationStatement
	var v llvm.Value

	for i := 0; ; i++ {
		if i+1 >= len(functionStmt.Declarations) {
			break
		}

		vdecl = &functionStmt.Declarations[i]
		v = *cg.generateVariableDeclaration(vdecl)
	}

	// insert expr statement
	for _, stmt := range functionStmt.Statements {
		v = cg.generateStatement(stmt)
	}

	return &v
}

func (cg *CodeGen) generateVariableDeclaration(vdecl *ast.DeclarationStatement) *llvm.Value {

	// create alloca
	alloca := cg.builder.CreateAlloca(llvm.Int32Type(), vdecl.Name.Name())

	// if args alloca
	if vdecl.GetDeclType() == ast.Param {
		// store args
		cg.builder.CreateStore(*cg.curFunc, alloca)
	}

	return &alloca
}

func (cg *CodeGen) generateStatement(stmt ast.Statement) llvm.Value {
	if infixStmt, ok := stmt.(*ast.InfixExpression); ok {
		return cg.generateInfixExpression(infixStmt)
	}

	if callStmt, ok := stmt.(*ast.CallExpression); ok {
		return cg.generateCallExpression(callStmt)
	}

	if returnStmt, ok := stmt.(*ast.ReturnStatement); ok {
		return cg.generateReturnStatement(returnStmt)
	}

	panic("generateStatement")
}

func (cg *CodeGen) generateInfixExpression(infixStmt *ast.InfixExpression) llvm.Value {
	lhs := infixStmt.Left
	rhs := infixStmt.Right
	var lhsValue, rhsValue llvm.Value

	if infixStmt.Operator == "=" {
		// lhs is variable
		lhsVariable := lhs.(*ast.Identifier)
		lhsValue = cg.builder.CreateLoad(*cg.curFunc, lhsVariable.Name())
	} else {
		// lhs ?
		// Binary ?
		if infixStmt, ok := lhs.(*ast.InfixExpression); ok {
			lhsValue = cg.generateInfixExpression(infixStmt)
		} else if identStmt, ok := lhs.(*ast.Identifier); ok {
			lhsValue = cg.generateIdentifier(identStmt)
		} else if numberStmt, ok := lhs.(*ast.Number); ok {
			lhsValue = cg.generateNumber(numberStmt.Val())
		}
	}

	// create rhs value
	if infixStmt, ok := rhs.(*ast.InfixExpression); ok {
		rhsValue = cg.generateInfixExpression(infixStmt)
	} else if callStmt, ok := rhs.(*ast.CallExpression); ok {
		rhsValue = cg.generateCallExpression(callStmt)
	} else if identStmt, ok := rhs.(*ast.Identifier); ok {
		rhsValue = cg.generateIdentifier(identStmt)
	} else if numberStmt, ok := rhs.(*ast.Number); ok {
		rhsValue = cg.generateNumber(numberStmt.Val())
	}

	// execute op
	switch infixStmt.Operator {
	case "=":
		return cg.builder.CreateStore(rhsValue, lhsValue)
	case "+":
		return cg.builder.CreateAdd(lhsValue, rhsValue, "add_tmp")
	case "-":
		return cg.builder.CreateSub(lhsValue, rhsValue, "sub_tmp")
	case "*":
		return cg.builder.CreateMul(lhsValue, rhsValue, "mul_tmp")
	case "/":
		return cg.builder.CreateSDiv(lhsValue, rhsValue, "div_tmp")
	default:
		panic("invalid operator")
	}
}

func (cg *CodeGen) generateCallExpression(callExpression *ast.CallExpression) llvm.Value {
	var argSlice []llvm.Value
	var argValue llvm.Value

	// 各引数について
	for _, arg := range callExpression.Arguments {
		if callExpr, ok := arg.(*ast.CallExpression); ok {
			argValue = cg.generateCallExpression(callExpr)
		} else if infixStmt, ok := arg.(*ast.InfixExpression); ok {
			argValue = cg.generateInfixExpression(infixStmt)
			if infixStmt.Operator == "=" {
				ident := infixStmt.Left.(*ast.Identifier)
				argValue = cg.builder.CreateLoad(cg.builder.CreateLoad(*cg.curFunc, ident.Name()), "arg_val")
			}
		} else if ident, ok := arg.(*ast.Identifier); ok {
			argValue = cg.generateIdentifier(ident)
		} else if numberStmt, ok := arg.(*ast.Number); ok {
			argValue = cg.generateNumber(numberStmt.Val())
		}
		argSlice = append(argSlice, argValue)
	}
	return cg.builder.CreateCall(cg.mod.NamedFunction(callExpression.GetCallee()), argSlice, "call_tmp")
}

func (cg *CodeGen) generateReturnStatement(retStmt *ast.ReturnStatement) llvm.Value {
	expr := retStmt.ReturnValue
	var retValue llvm.Value

	if infixStmt, ok := expr.(*ast.InfixExpression); ok {
		retValue = cg.generateInfixExpression(infixStmt)
	} else if ident, ok := expr.(*ast.Identifier); ok {
		retValue = cg.generateIdentifier(ident)
	} else if numberStmt, ok := expr.(*ast.Number); ok {
		retValue = cg.generateNumber(numberStmt.Val())
	}

	return cg.builder.CreateRet(retValue)
}

func (cg *CodeGen) generateIdentifier(ident *ast.Identifier) llvm.Value {
	return cg.builder.CreateLoad(*cg.curFunc, "var_tmp")
}

func (cg *CodeGen) generateNumber(value int) llvm.Value {
	return llvm.ConstInt(llvm.Int32Type(), uint64(value), false)
}
