library work;
   context work.cosim_context;
   use work.cosim.all;

library vfbdb;
   use vfbdb.main_pkg.all;


entity tb_cosim is
   generic(
      G_SW_GW_FIFO_PATH : string;
      G_GW_SW_FIFO_PATH : string
   );
end entity;


architecture test of tb_cosim is

   signal clk : std_logic := '0';

   -- Wishbone interfaces
   signal uvvm_wb_if : t_wishbone_if (
      dat_o(31 downto 0),
      dat_i(31 downto 0),
      adr_o(31 downto 0)
   ) := init_wishbone_if_signals(32, 32);

   signal wb_ms: t_wishbone_master_out;
   signal wb_sm: t_wishbone_slave_out;

   -- Testbench specific signals
   signal add : add_t;
   signal add_stb : std_logic;

   signal result : result_t;
   signal result_stb : std_logic;

   signal buff : slv_vector(0 to to_integer(wb3.main_pkg.DEPTH) - 1)(40 downto 0);
   signal buff_write_ptr, buff_read_ptr : natural := 0;

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

      add_o     => add,
      add_stb_o => add_stb,

      result_i     => result,
      result_stb_o => result_stb
   );


   Write_Driver : process (clk) is
   begin
      if rising_edge(clk) then
         if add_stb = '1' then
            buff(buff_write_ptr) <= std_logic_vector(
               resize(unsigned(add.a), buff(0)'length) +
               resize(unsigned(add.b), buff(0)'length) +
               resize(unsigned(add.c), buff(0)'length)
            );
            buff_write_ptr <= buff_write_ptr + 1;
         end if;
      end if;
   end process;


   result.res <= buff(buff_read_ptr);


   Read_Driver : process (clk) is
   begin
      if rising_edge(clk) then
         if result_stb = '1' then
            buff_read_ptr <= (buff_read_ptr + 1) mod to_integer(wb3.main_pkg.DEPTH);
         end if;
      end if;
   end process;

end architecture;
