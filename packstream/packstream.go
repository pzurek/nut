package packstream

// PackStream marker declarations
const (
	TinyString  = 0x80 // 128
	TinyList    = 0x90 // 144
	TinyMap     = 0xA0 // 160
	TinyStruct  = 0xB0 // 176
	Null        = 0xC0 // 192
	Float64     = 0xC1 // 193
	False       = 0xC2 // 194
	True        = 0xC3 // 195
	Int8        = 0xC8 // 200
	Int16       = 0xC9 // 201
	Int32       = 0xCA // 202
	Int64       = 0xCB // 203
	Bytes8      = 0xCC // 204
	Bytes16     = 0xCD // 205
	Bytes32     = 0xCE // 206
	String8     = 0xD0 // 208
	String16    = 0xD1 // 209
	String32    = 0xD2 // 210
	List8       = 0xD4 // 212
	List16      = 0xD5 // 213
	List32      = 0xD6 // 214
	ListStream  = 0xD7 // 215
	Map8        = 0xD8 // 216
	Map16       = 0xD9 // 217
	Map32       = 0xDA // 218
	MapStream   = 0xDB // 219
	Struct8     = 0xDC // 220
	Struct16    = 0xDD // 221
	EndOfStream = 0xDF // 223

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
