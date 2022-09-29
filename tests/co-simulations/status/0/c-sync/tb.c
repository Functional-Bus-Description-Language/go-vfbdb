#include <assert.h>
#include <stdio.h>

#include "cosim_iface.h"

#include "vfbdb/vfbdb.h"
#include "vfbdb/Main.h"
#define VFBDB_IFACE &iface


int main(int argc, char *argv[]) {
	assert(argc == 3);

	vfbdb_iface_t iface = cosim_iface_iface();

	cosim_iface_init(argv[1], argv[2], NULL);

	uint32_t id;
	vfbdb_read(Main_ID, &id);
	if (id != VFBDB_ID) {
		fprintf(stderr, "read wrong ID %x, expecting %x\n", id, VFBDB_ID);
		cosim_iface_end(1);
	}

	uint32_t timestamp;
	vfbdb_read(Main_TIMESTAMP, &timestamp);
	if (timestamp != VFBDB_TIMESTAMP) {
		fprintf(stderr, "read wrong TIMESTAMP %x, expecting %x\n", id, VFBDB_TIMESTAMP);
		cosim_iface_end(1);
	}

	cosim_iface_end(0);

	return 0;
}
