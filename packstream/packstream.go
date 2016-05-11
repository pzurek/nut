package packstream

// PackStream marker declarations
const (
	TinyString  = byte(0x80) // 128
	TinyList    = byte(0x90) // 144
	TinyMap     = byte(0xA0) // 160
	TinyStruct  = byte(0xB0) // 176
	Null        = byte(0xC0) // 192
	Float64     = byte(0xC1) // 193
	False       = byte(0xC2) // 194
	True        = byte(0xC3) // 195
	Int8        = byte(0xC8) // 200
	Int16       = byte(0xC9) // 201
	Int32       = byte(0xCA) // 202
	Int64       = byte(0xCB) // 203
	Bytes8      = byte(0xCC) // 204
	Bytes16     = byte(0xCD) // 205
	Bytes32     = byte(0xCE) // 206
	String8     = byte(0xD0) // 208
	String16    = byte(0xD1) // 209
	String32    = byte(0xD2) // 210
	List8       = byte(0xD4) // 212
	List16      = byte(0xD5) // 213
	List32      = byte(0xD6) // 214
	ListStream  = byte(0xD7) // 215
	Map8        = byte(0xD8) // 216
	Map16       = byte(0xD9) // 217
	Map32       = byte(0xDA) // 218
	MapStream   = byte(0xDB) // 219
	Struct8     = byte(0xDC) // 220
	Struct16    = byte(0xDD) // 221
	EndOfStream = byte(0xDF) // 223

	MinTinyInt = -16
	MaxTinyInt = 127
)

// Type is used to define PackStream type enumerations
type Type int

// PackStream type enumerations
const (
	PSNull Type = iota
	PSBool
	PSInt
	PSFloat
	PSBytes
	PSString
	PSList
	PSMap
	PSStruct
)
