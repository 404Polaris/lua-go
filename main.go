package main

import (
	"fmt"
	"io/ioutil"
	"olua/luachunk"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])

		if err != nil {
			panic(err)
		}

		proto := luachunk.UnDump(data)
		list(proto)
	}
}

func list(proto *luachunk.Prototype) {
	printHeader(proto)
	printCode(proto)
	printDetail(proto)

	for _, p := range proto.Protos {
		list(p)
	}
}

func printHeader(proto *luachunk.Prototype) {
	funcType := "main"

	if proto.LineDefined > 0 {
		funcType = "function"
	}

	varArgFlag := ""

	if proto.IsVararg > 0 {
		varArgFlag = "+"
	}

	fmt.Printf("\n%s <%s:%d,%d> (%d instructions)\n\n", funcType, proto.Source, proto.LineDefined, proto.LastLineDefined, len(proto.Code))
	fmt.Printf("%d%s params, %d slots, %d upvalues,", proto.NumParams, varArgFlag, proto.MaxStackSize, len(proto.Upvalues))
	fmt.Printf("%d locals, %d constants, %d functions\n", len(proto.LocVars), len(proto.Constants), len(proto.Protos))
}

func printCode(proto *luachunk.Prototype) {
	for pc, c := range proto.Code {
		line := "_"
		if len(proto.LineInfo) > 0 {
			line = fmt.Sprintf("%d", proto.LineInfo[pc])
		}

		fmt.Printf("\t%d\t[%s]\t0x%08X\n", pc+1, line, c)
	}
}

func printDetail(proto *luachunk.Prototype) {
	fmt.Printf("constants (%d):\n", len(proto.Constants))
	for i, k := range proto.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantsToString(k))
	}
	fmt.Printf("locals (%d):\n", len(proto.LocVars))
	for i, locVar := range proto.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}
	fmt.Printf("upvalues (%d):\n", len(proto.Upvalues))
	for i, upvalue := range proto.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, upvalueName(proto, i), upvalue.InStack, upvalue.Idx)
	}
}

func upvalueName(proto *luachunk.Prototype, i int) interface{} {
	if len(proto.UpvalueNames) > 0 {
		return proto.UpvalueNames[i]
	}

	return "-"
}

func constantsToString(val interface{}) string {
	switch val.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", val)
	case float64:
		return fmt.Sprintf("%g", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case string:
		return fmt.Sprintf("%q", val)
	default:
		return "?"
	}
}
