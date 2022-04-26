#include <assert.h>
#include <stdio.h>

#include "cosim_iface.h"

#include "wbfbd/wbfbd.h"
#include "wbfbd/Main.h"
#define WBFBD_IFACE &iface


int main(int argc, char *argv[]) {
	assert(argc == 3);

	wbfbd_iface_t iface = cosim_iface_iface();

	cosim_iface_init(argv[1], argv[2], NULL);

	uint32_t id;
	wbfbd_read(Main_ID, &id);
	if (id != WBFBD_ID) {
		fprintf(stderr, "read wrong ID %x, expecting %x\n", id, WBFBD_ID);
		cosim_iface_end(1);
	}

	uint32_t timestamp;
	wbfbd_read(Main_TIMESTAMP, &timestamp);
	if (timestamp != WBFBD_TIMESTAMP) {
		fprintf(stderr, "read wrong TIMESTAMP %x, expecting %x\n", id, WBFBD_ID);
		cosim_iface_end(1);
	}

	cosim_iface_end(0);

	return 0;
}
