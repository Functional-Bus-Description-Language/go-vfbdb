#ifndef _WBFBD_WBFBD_H_
#define _WBFBD_WBFBD_H_

#include <stdint.h>

typedef {{.AddrType}} wbfbd_addr_t;

extern const uint32_t WBFBD_ID;
extern const uint32_t WBFBD_TIMESTAMP;

typedef struct {
	int (*read)(const wbfbd_addr_t addr, {{.ReadDataType}} const data);
	int (*write)(const wbfbd_addr_t addr, const {{.WriteDataType}} data);
} wbfbd_iface_t;

#define wbfbd_read(elem, data) (wbfbd_ ## elem ## _read(WBFBD_IFACE, data))
#define wbfbd_write(elem, data) (wbfbd_ ## elem ## _write(WBFBD_IFACE, data))

#ifdef WBFBD_SHORT_MACROS
	#undef wbfbd_read
	#undef wbfbd_write
	#define read(elem, data) (wbfbd_ ## elem ## _read(WBFBD_IFACE, data))
	#define write(elem, data) (wbfbd_ ## elem ## _write(WBFBD_IFACE, data))
#endif

#endif // _WBFBD_WBFBD_H_
