package utils

// WidthToCTypeRead returns type that is sufficient to represent data
// of given width in the C language for read functions.
func WidthToCTypeRead(width int64) string {
	if width > 64 {
		return "uint8_t *"
	} else if width > 32 {
		return "uint64_t *"
	} else if width > 16 {
		return "uint32_t *"
	} else if width > 8 {
		return "uint8_t *"
	}
	return "uint8_t *"
}

// WidthToCTypeWrite returns type that is sufficient to represent data
// of given width in the C language for write functions.
func WidthToCTypeWrite(width int64) string {
	if width > 64 {
		return "uint8_t *"
	} else if width > 32 {
		return "uint64_t"
	} else if width > 16 {
		return "uint32_t"
	} else if width > 8 {
		return "uint8_t"
	}
	return "uint8_t"
}
