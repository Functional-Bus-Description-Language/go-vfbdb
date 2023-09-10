package c

type Type interface {
	String() string
	Typ() string
	Depointer() Type
}

type ByteArray struct{}

func (ba ByteArray) String() string  { return "uint8_t *" }
func (ba ByteArray) Typ() string     { return "ByteArray" }
func (ba ByteArray) Depointer() Type { panic("unimplemented") }

type Uint8 struct{}

func (u8 Uint8) String() string  { return "uint8_t" }
func (u8 Uint8) Typ() string     { return "Uint8" }
func (u8 Uint8) Depointer() Type { return u8 }

type Uint8Ptr struct{}

func (u8 Uint8Ptr) String() string  { return "uint8_t *" }
func (u8 Uint8Ptr) Typ() string     { return "Uint8Ptr" }
func (u8 Uint8Ptr) Depointer() Type { return Uint8{} }

type Uint16 struct{}

func (u16 Uint16) String() string  { return "uint16_t" }
func (u16 Uint16) Typ() string     { return "Uint16" }
func (u16 Uint16) Depointer() Type { return u16 }

type Uint16Ptr struct{}

func (u16 Uint16Ptr) String() string  { return "uint16_t *" }
func (u16 Uint16Ptr) Typ() string     { return "Uint16Ptr" }
func (u16 Uint16Ptr) Depointer() Type { return Uint16{} }

type Uint32 struct{}

func (u32 Uint32) String() string  { return "uint32_t" }
func (u32 Uint32) Typ() string     { return "Uint32" }
func (u32 Uint32) Depointer() Type { return u32 }

type Uint32Ptr struct{}

func (u32 Uint32Ptr) String() string  { return "uint32_t *" }
func (u32 Uint32Ptr) Typ() string     { return "Uint32Ptr" }
func (u32 Uint32Ptr) Depointer() Type { return Uint32{} }

type Uint64 struct{}

func (u64 Uint64) String() string  { return "uint64_t" }
func (u64 Uint64) Typ() string     { return "Uint64" }
func (u64 Uint64) Depointer() Type { return u64 }

type Uint64Ptr struct{}

func (u64 Uint64Ptr) String() string  { return "uint64_t *" }
func (u64 Uint64Ptr) Typ() string     { return "Uint64Ptr" }
func (u64 Uint64Ptr) Depointer() Type { return Uint64{} }
