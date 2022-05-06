library work;
   context work.cosim_context;
   use work.cosim.all;


entity tb_cosim is
   generic(
      G_SW_GW_FIFO_PATH : string;
      G_GW_SW_FIFO_PATH : string
   );
end entity;


architecture test of tb_cosim is

   signal clk : std_logic := '0';

   signal cfg0, cfg1, cfg2 : std_logic_vector(31 downto 0);

   -- Wishbone interfaces.
   signal uvvm_wb_if : t_wishbone_if (
      dat_o(31 downto 0),
      dat_i(31 downto 0),
      adr_o(31 downto 0)
   ) := init_wishbone_if_signals(32, 32);

   signal wb_ms, wb_blk0_ms, wb_blk1_ms, wb_blk2_ms : t_wishbone_master_out;
   signal wb_sm, wb_blk0_sm, wb_blk1_sm, wb_blk2_sm : t_wishbone_slave_out;

begin

   clk <= not clk after C_CLK_PERIOD / 2;


   wb_ms.cyc <= uvvm_wb_if.cyc_o;
   wb_ms.stb <= uvvm_wb_if.stb_o;
   wb_ms.adr <= uvvm_wb_if.adr_o;
   wb_ms.sel <= (others => '0');
   wb_ms.we  <= uvvm_wb_if.we_o;
   wb_ms.dat <= uvvm_wb_if.dat_o;

   uvvm_wb_if.dat_i <= wb_sm.dat;
   uvvm_wb_if.ack_i <= wb_sm.ack;

   cosim_interface(G_SW_GW_FIFO_PATH, G_GW_SW_FIFO_PATH, clk, uvvm_wb_if, C_WB_BFM_CONFIG);


   vfbdb_main : entity vfbdb.Main
   port map (
      clk_i => clk,
      rst_i => '0',
      slave_i(0) => wb_ms,
      slave_o(0) => wb_sm,
      blk0_master_o(0) => wb_blk0_ms,
      blk0_master_i(0) => wb_blk0_sm,
      blk1_master_o(0) => wb_blk1_ms,
      blk1_master_i(0) => wb_blk1_sm
   );


   vfbdb_blk0 : entity vfbdb.Blk0
   port map (
      clk_i => clk,
      rst_i => '0',
      slave_i(0) => wb_blk0_ms,
      slave_o(0) => wb_blk0_sm,
      cfg_o => cfg0,
      st_i => cfg0
   );


   vfbdb_blk1 : entity vfbdb.Blk1
   port map (
      clk_i => clk,
      rst_i => '0',
      slave_i(0) => wb_blk1_ms,
      slave_o(0) => wb_blk1_sm,
      blk2_master_o(0) => wb_blk2_ms,
      blk2_master_i(0) => wb_blk2_sm,
      cfg_o => cfg1,
      st_i => cfg1
   );


   vfbdb_blk2 : entity vfbdb.Blk2
   port map (
      clk_i => clk,
      rst_i => '0',
      slave_i(0) => wb_blk2_ms,
      slave_o(0) => wb_blk2_sm,
      cfg_o => cfg2,
      st_i => cfg2
   );

end architecture;
