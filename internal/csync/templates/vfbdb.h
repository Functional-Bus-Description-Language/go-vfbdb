#ifndef _VFBDB_VFBDB_H_
#define _VFBDB_VFBDB_H_

#include <stdint.h>

extern const uint32_t VFBDB_ID;
extern const uint32_t VFBDB_TIMESTAMP;

typedef struct {
	int (*read)(const {{.AddrType}} addr, {{.ReadDataType}} const data);
	int (*write)(const {{.AddrType}} addr, const {{.WriteDataType}} data);
} vfbdb_iface_t;

#define vfbdb_read(elem, data) (vfbdb_ ## elem ## _read(VFBDB_IFACE, data))
#define vfbdb_write(elem, data) (vfbdb_ ## elem ## _write(VFBDB_IFACE, data))

#ifdef VFBDB_SHORT_MACROS
	#undef vfbdb_read
	#undef vfbdb_write
	#define read(elem, data) (vfbdb_ ## elem ## _read(VFBDB_IFACE, data))
	#define write(elem, data) (vfbdb_ ## elem ## _write(VFBDB_IFACE, data))
#endif

#endif // _VFBDB_VFBDB_H_
