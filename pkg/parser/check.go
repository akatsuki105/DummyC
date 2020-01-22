package parser

func (p *Parser) checkReDefinition(fn Function) (ok bool) {
	for _, prototype := range p.prototypeTable {
		if fn.Name == prototype.Name {
			// 既にプロトタイプ宣言されている
			return false
		}
	}

	for _, function := range p.functionTable {
		if fn.Name == function.Name {
			if fn.Argc != function.Argc {
				// 関数が既に定義されているが引数の数が合わないとき
				return false
			}
		}
	}

	return true
}

func (p *Parser) checkCorrectDefinition(fn Function) (ok bool) {
	// 関数がプロトタイプ宣言されていた場合に、プロトタイプ宣言と関数定義の引数の数が同一であることを確認する
	correctDeclared := true
	for _, prototype := range p.prototypeTable {
		if fn.Name == prototype.Name {
			if fn.Argc == prototype.Argc {
				correctDeclared = true
				break
			} else {
				correctDeclared = false
			}
		}
	}
	if !correctDeclared {
		return false
	}

	isReDefined := false
	for _, function := range p.functionTable {
		if fn.Name == function.Name {
			isReDefined = true
			break
		}
	}
	if isReDefined {
		// 同名の関数が既に定義されている
		return false
	}

	return true
}
