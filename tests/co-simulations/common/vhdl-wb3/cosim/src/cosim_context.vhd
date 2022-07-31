context cosim_context is
   library ieee;
      use ieee.std_logic_1164.all;
      use ieee.numeric_std.all;

   library vfbdb;
      use vfbdb.wb3;
      use vfbdb.wb3.all;

   library general_cores;
      use general_cores.wishbone_pkg.all;

   library ltypes;
      use ltypes.types.all;

   library uvvm_util;
      context uvvm_util.uvvm_util_context;

   library bitvis_vip_wishbone;
      use bitvis_vip_wishbone.wishbone_bfm_pkg.all;
end context;
