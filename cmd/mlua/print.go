package main

import (
	"strconv"
	"strings"

	"github.com/ocurr/Meia-Lua/ast"
)

// PrintLua returns a properly formatted text representation of the given Lua AST.
func PrintLua(node ast.Node) string {

	switch n := node.(type) {
	case ast.Chunk:
		return PrintLua(n.Block)
	case ast.Block:
		var sb strings.Builder
		for _, s := range n.StatLst {
			sb.Write([]byte(PrintLua(s) + "\n"))
		}
		return sb.String()
	case ast.Cond:
		cnd := PrintLua(n.Cnd)
		block := PrintLua(n.Block)
		elseifs := ""
		for _, c := range n.Elseifs {
			elseifs += "else" + PrintLua(c)
		}
		elseifs = strings.TrimSuffix(elseifs, "end")
		els := ""
		if len(n.Else.StatLst) != 0 {
			els = "else\n" + PrintLua(n.Else)
		}

		return "if " + cnd + " then\n" +
			block + elseifs + els + "end"
	case ast.While:
		cnd := PrintLua(n.Cnd)
		block := PrintLua(n.Block)

		return "while " + cnd + "\n" + " do\n" + block + "\n" + "end"
	case ast.For:
		def := PrintLua(n.Assign)
		cnd := PrintLua(n.Cnd)
		step := PrintLua(n.Step)
		block := PrintLua(n.Block)

		if step != "" {
			step = ", " + step
		}

		return "for " + def + ", " + cnd + step + " do\n" +
			block + "end"
	case ast.Def:
		id := PrintLua(n.Id)
		exp := PrintLua(n.Exp)
		scope := ""
		if n.Scope == ast.LOCAL {
			scope = "local "
		}
		return scope + id + " = " + exp
	case ast.DefLst:
		defs := ""
		exps := ""
		sep := ""
		if len(n.List) > 1 {
			sep = ", "
		}
		for _, d := range n.List {
			def := PrintLua(d)
			def = strings.TrimPrefix(def, "local ")
			split := strings.Split(def, " = ")
			defs += strings.TrimSpace(split[0]) + sep
			exps += strings.TrimSpace(split[1]) + sep
		}
		defs = strings.TrimRight(defs, ", ")
		exps = strings.TrimRight(exps, ", ")
		return defs + " = " + exps
	case ast.Id:
		return n.Id
	case ast.BinaryOp:
		lhs := PrintLua(n.Lhs)
		rhs := PrintLua(n.Rhs)
		return lhs + " " + n.Op + " " + rhs
	case ast.Int:
		return strconv.FormatInt(n.N, 10)
	case ast.Float:
		return strconv.FormatFloat(n.N, 'f', -1, 64)
	case ast.String:
		return n.S
	case ast.Bool:
		if n.True {
			return "true"
		}
		return "false"
	case ast.Nil:
		return "nil"
	default:
		return ""
	}
}
