#include <assert.h>
#include <stdbool.h>
#include <stdlib.h>
#include <stdio.h>
#include <stddef.h>
#include <string.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <errno.h>

#include "cosim_iface.h"

#define RESPONSE_BUFFER_SIZE 64

static char *write_fifo_path;
static FILE *write_fifo;
static char *read_fifo_path;
static FILE *read_fifo;
static delay_function_t delay_function;

// Statistic variables.
static uint32_t write_count;
static uint32_t read_count;
static uint32_t rmw_count;

static void remove_fifos(void) {
	printf("Removing FIFOs\n");
	// Ignore errors, file probably doesn't exist yet.
	remove(write_fifo_path);
	write_fifo = NULL;
	remove(read_fifo_path);
	read_fifo = NULL;
}

static void make_fifos(void) {
	remove_fifos();
	printf("Making FIFOs\n");
	if (mkfifo(write_fifo_path, 0755)) {
		perror("making write FIFO");
		exit(EXIT_FAILURE);
	}
	if (mkfifo(read_fifo_path, 0755)) {
		perror("making write FIFO");
		exit(EXIT_FAILURE);
	}
}

static void cosim_iface_wait(uint32_t time_ns) {
	printf("wait for %u ns", time_ns);

	int count;
	count = fprintf(write_fifo, "T%.8x\n", time_ns);
	if (count != 10) {
		fprintf(stderr, "wait: fprintf: not all characters written\n");
		exit(EXIT_FAILURE);
	}
	if (fflush(write_fifo)) {
		perror("flushing write FIFO");
		exit(EXIT_FAILURE);
	}

	char response[RESPONSE_BUFFER_SIZE];
	if (fgets(response, RESPONSE_BUFFER_SIZE, read_fifo) == NULL) {
		perror("reading read FIFO");
		exit(EXIT_FAILURE);
	}
	if (strncmp(response, "ACK\n", 4) != 0) {
		fprintf(stderr, "wrong status returned: %s", response);
		exit(EXIT_FAILURE);
	}
}

static int cosim_iface_write(const uint8_t addr, const uint32_t data) {
	if (delay_function) {
		cosim_iface_wait(delay_function());
	}

	printf("write: addr %#.8x, data %u (%#.8x)\n", addr, data, data);
	int count;
	count = fprintf(write_fifo, "W%.8x,%.8x\n", addr, data);
	if (count != 19) {
		fprintf(stderr, "write: fprintf: not all characters written\n");
		exit(EXIT_FAILURE);
	}
	if (fflush(write_fifo)) {
		perror("flushing write FIFO");
		exit(EXIT_FAILURE);
	}

	char response[RESPONSE_BUFFER_SIZE];
	if (fgets(response, RESPONSE_BUFFER_SIZE, read_fifo) == NULL) {
		perror("reading read FIFO");
		exit(EXIT_FAILURE);
	}
	if (strncmp(response, "ACK\n", 4) != 0) {
		fprintf(stderr, "wrong status returned: %s", response);
		exit(EXIT_FAILURE);
	}

	write_count++;

	return 0;
}


static int cosim_iface_writeb(const uint8_t addr, const uint32_t * buf, size_t count) {
	fprintf(stderr, "cosim iface: cosim_iface_writeb unimplemented");
	exit(EXIT_FAILURE);
}

static uint32_t bin_to_uint32(const char * const s) {
	uint32_t u32 = 0;

	for (int i = 0; i < 32; i ++) {
		assert(s[i] == '0' || s[i] == '1');
		if (s[i] == '1') {
			u32 |= 1 << (31 - i);
		}
	}

	return u32;
}

static int cosim_iface_read(const uint8_t addr, uint32_t * const data) {
	if (delay_function) {
		cosim_iface_wait(delay_function());
	}

	int count;
	count = fprintf(write_fifo, "R%.8x\n", addr);
	if (count != 10) {
		fprintf(stderr, "read: fprintf: not all characters written\n");
		exit(EXIT_FAILURE);
	}
	if (fflush(write_fifo)) {
		perror("flushing write FIFO");
		exit(EXIT_FAILURE);
	}

	char response[RESPONSE_BUFFER_SIZE];
	if (fgets(response, RESPONSE_BUFFER_SIZE, read_fifo) == NULL) {
		perror("reading read FIFO");
		exit(EXIT_FAILURE);
	}
	if (strncmp(response, "ERR\n", 4) == 0) {
		fprintf(stderr, "error status returned");
		exit(EXIT_FAILURE);
	}

	uint32_t aux = bin_to_uint32(response);

	printf("read: data %u (%#.8x)\n", aux, aux);

	*data = aux;

	read_count++;

	return 0;
}


static int cosim_iface_readb(const uint8_t addr, uint32_t * buf, size_t count) {
	fprintf(stderr, "cosim iface: cosim_iface_readb unimplemented");
	exit(EXIT_FAILURE);
}

static void cosim_iface_atexit(void) {
	static bool atexit = false;
	if (atexit || write_fifo == NULL || read_fifo == NULL) {
		return;
	}
	atexit = true;
	cosim_iface_end(1);
	atexit = false;
}

void cosim_iface_init(char *wr_fifo_path, char *rd_fifo_path, delay_function_t delay_func) {
	write_fifo_path = wr_fifo_path;
	read_fifo_path = rd_fifo_path;
	delay_function = delay_func;

	make_fifos();

	write_fifo = fopen(write_fifo_path, "w");
	if (write_fifo == NULL) {
		perror("opening write FIFO");
		exit(EXIT_FAILURE);
	}
	read_fifo = fopen(read_fifo_path, "r");
	if (read_fifo == NULL) {
		perror("opening read FIFO");
		exit(EXIT_FAILURE);
	}

	if (atexit(cosim_iface_atexit)) {
		fprintf(stderr, "cannot register atexit() function");
		exit(EXIT_FAILURE);
	}
}

vfbdb_iface_t cosim_iface_iface(void) {
	vfbdb_iface_t iface = {
		read: cosim_iface_read,
		write: cosim_iface_write,
		readb: cosim_iface_readb,
		writeb: cosim_iface_writeb
	};
	return iface;
}

static void cosim_iface_print_stats(void) {
	printf("cosim iface: transactions statistics:\n"
		"  Write Count: %d\n"
		"  Read Count:  %d\n"
		"  RMW Count:   %d\n",
		write_count, read_count, rmw_count
	);
}

void cosim_iface_end(int status) {
	printf("cosim iface: ending with status %d\n", status);

	int count;
	count = fprintf(write_fifo, "E%.8x\n", status);
	if (count != 10) {
		fprintf(stderr, "write: fprintf: not all characters written\n");
		exit(EXIT_FAILURE);
	}
	if (fflush(write_fifo)) {
		perror("flushing write FIFO");
		exit(EXIT_FAILURE);
	}

	char response[RESPONSE_BUFFER_SIZE];
	if (fgets(response, RESPONSE_BUFFER_SIZE, read_fifo) == NULL) {
		perror("reading read FIFO");
		exit(EXIT_FAILURE);
	}
	if (strncmp(response, "ACK\n", 4) != 0) {
		fprintf(stderr, "error status returned");
		exit(EXIT_FAILURE);
	}

	remove_fifos();
	cosim_iface_print_stats();
}
