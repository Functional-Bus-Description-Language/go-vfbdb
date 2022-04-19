#ifndef _WBFBD_WBFBD_H_
#define _WBFBD_WBFBD_H_

#include <stdint.h>

const uint32_t WBFBD_ID = {{.ID}};
const uint32_t WBFBD_TIMESTAMP = {{.TIMESTAMP}};

struct wbfbd_iface_t {
	int (*read)(const {{.AddrType}} addr, {{.ReadDataType}} const data);
	int (*write)(const {{.AddrType}} addr, const {{.WriteDataType}} data);
};

#define wbfbd_read(elem, data) (wbfbd_ ## elem ## _read(WBFBD_IFACE, data))
#define wbfbd_write(elem, data) (wbfbd_ ## elem ## _write(WBFBD_IFACE, data))

#ifdef WBFBD_SHORT_MACROS
	#undef wbfbd_read
	#undef wbfbd_write
	#define read(elem, data) (wbfbd_ ## elem ## _read(WBFBD_IFACE, data))
	#define write(elem, data) (wbfbd_ ## elem ## _write(WBFBD_IFACE, data))
#endif

#endif // _WBFBD_WBFBD_H_
