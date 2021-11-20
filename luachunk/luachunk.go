package luachunk

const (
	LuaSignature    = "\x1bLua"
	LuacVersion     = 0x53
	LuacFormat      = 0
	LuacData        = "\x19\x93\r\n\x1a\n"
	CIntSize        = 4
	CSizetSize      = 8
	InstructionSize = 4
	LuaIntegerSize  = 8
	LuaNumberSize   = 8
	LuacInt         = 0x5678
	LuacNumber      = 370.5
)

const (
	TagNil      = 0x00
	TagBoolean  = 0x01
	TagInteger  = 0x13
	TagNumber   = 0x03
	TagShortStr = 0x04
	TagLongStr  = 0x14
)

type Upvalue struct {
	InStack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

type Prototype struct {
	Source          string
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVararg        byte
	MaxStackSize    byte
	Code            []uint32
	Constants       []interface{}
	Upvalues        []Upvalue
	Protos          []*Prototype
	LineInfo        []uint32
	LocVars         []LocVar
	UpvalueNames    []string
}

type header struct {
	signature       [4]byte
	version         byte
	format          byte
	luacData        [6]byte
	cintSize        byte
	sizetSize       byte
	instructionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte
	luacInt         int64
	luacNumber      int64
}

type luaChunk struct {
	header        header
	sizeUpValues  byte
	mainFuncProto *Prototype
}

func UnDump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()
	reader.readByte()
	return reader.readProto("")
}
