package luachunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte
}

func (r *reader) readByte() byte {
	b := r.data[0]
	r.data = r.data[1:]
	return b
}

func (r *reader) readUint32() uint32 {
	val := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:]
	return val
}

func (r *reader) readUint64() uint64 {
	val := binary.LittleEndian.Uint64(r.data)
	r.data = r.data[8:]
	return val
}

func (r *reader) readLuaInteger() int64 {
	return int64(r.readUint64())
}

func (r *reader) readLuaNumber() float64 {
	return math.Float64frombits(r.readUint64())
}

func (r *reader) readString() string {
	size := uint(r.readByte())

	if size == 0 {
		return ""
	}

	if size == 0xFF {
		size = uint(r.readUint64())
	}

	bytes := r.readBytes(size - 1)
	return string(bytes)
}

func (r *reader) readBytes(size uint) []byte {
	bytes := r.data[:size]
	r.data = r.data[size:]
	return bytes
}

func (r *reader) checkHeader() {
	if string(r.readBytes(4)) != LuaSignature {
		panic("not a precompiled chunk")
	} else if r.readByte() != LuacVersion {
		panic("version mismatch")
	} else if r.readByte() != LuacFormat {
		panic("format mismatch")
	} else if string(r.readBytes(6)) != LuacData {
		panic("corrupted")
	} else if r.readByte() != CIntSize {
		panic("int size mismatch")
	} else if r.readByte() != CSizetSize {
		panic("size_t size mismatch")
	} else if r.readByte() != InstructionSize {
		panic("instruction size mismatch")
	} else if r.readByte() != LuaIntegerSize {
		panic("lua_Integer size mismatch")
	} else if r.readByte() != LuaNumberSize {
		panic("lua_Number size mismatch")
	} else if r.readLuaInteger() != LuacInt {
		panic("endianness mismatch")
	} else if r.readLuaNumber() != LuacNumber {
		panic("float format mismatch")
	}
}

func (r *reader) readProto(parentSource string) *Prototype {
	source := r.readString()

	if source == "" {
		source = parentSource
	}

	return &Prototype{
		Source:          source,
		LineDefined:     r.readUint32(),
		LastLineDefined: r.readUint32(),
		NumParams:       r.readByte(),
		IsVararg:        r.readByte(),
		MaxStackSize:    r.readByte(),
		Code:            r.readCode(),
		Constants:       r.readConstants(),
		Upvalues:        r.readUpvalues(),
		Protos:          r.readProtos(source),
		LineInfo:        r.readLineInfo(),
		LocVars:         r.readLocVars(),
		UpvalueNames:    r.readUpvalueNames(),
	}
}

func (r *reader) readProtos(source string) []*Prototype {
	protos := make([]*Prototype, r.readUint32())
	return protos
}

func (r *reader) readUpvalueNames() []string {
	names := make([]string, r.readUint32())
	for i := range names {
		names[i] = r.readString()
	}

	return names
}

func (r *reader) readCode() []uint32 {
	code := make([]uint32, r.readUint32())
	for i := range code {
		code[i] = r.readUint32()
	}

	return code
}

func (r *reader) readConstants() []interface{} {
	constants := make([]interface{}, r.readUint32())
	for i := range constants {
		constants[i] = r.readConstant()
	}

	return constants
}

func (r *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, r.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			InStack: r.readByte(),
			Idx:     r.readByte(),
		}
	}

	return upvalues
}

func (r *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, r.readUint32())
	for i := range lineInfo {
		lineInfo[i] = r.readUint32()
	}

	return lineInfo
}

func (r *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, r.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: r.readString(),
			StartPC: r.readUint32(),
			EndPC:   r.readUint32(),
		}
	}

	return locVars
}

func (r *reader) readConstant() interface{} {
	switch r.readByte() {
	case TagNil:
		return nil
	case TagBoolean:
		return r.readByte() != 0
	case TagInteger:
		return r.readLuaInteger()
	case TagNumber:
		return r.readLuaNumber()
	case TagShortStr:
		return r.readString()
	case TagLongStr:
		return r.readString()
	default:
		panic("corrupted")
	}
}
