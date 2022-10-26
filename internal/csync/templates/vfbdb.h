#ifndef _VFBDB_VFBDB_H_
#define _VFBDB_VFBDB_H_

#include <stddef.h>
#include <stdint.h>

typedef struct {
	int (*read)(const {{.AddrType}} addr, {{.ReadType}} const data);
	int (*write)(const {{.AddrType}} addr, const {{.WriteType}} data);
	int (*readb)(const {{.AddrType}} addr, {{.ReadType}} buf, size_t count);
	int (*writeb)(const {{.AddrType}} addr, const {{.WriteType}} * buf, size_t count);
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
