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

   signal add : add_out_t;
   signal result : std_logic_vector(31 downto 0) := (others => '0');

   -- Wishbone interfaces.
   signal uvvm_wb_if : t_wishbone_if (
      dat_o(31 downto 0),
      dat_i(31 downto 0),
      adr_o(31 downto 0)
   ) := init_wishbone_if_signals(32, 32);

   signal wb_ms: t_wishbone_master_out;
   signal wb_sm: t_wishbone_slave_out;

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
      add_o    => add,
      result_i => result
   );


   adder : process (clk) is
   begin
      if rising_edge(clk) then
         if add.call = '1' then
            result <= std_logic_vector(
               resize(unsigned(add.terms(0)), result'length) +
               resize(unsigned(add.terms(1)), result'length) +
               resize(unsigned(add.terms(2)), result'length) +
               resize(unsigned(add.terms(3)), result'length) +
               resize(unsigned(add.terms(4)), result'length) +
               resize(unsigned(add.terms(5)), result'length) +
               resize(unsigned(add.terms(6)), result'length) +
               resize(unsigned(add.terms(7)), result'length) +
               resize(unsigned(add.terms(8)), result'length) +
               resize(unsigned(add.terms(9)), result'length)
            );
         end if;
      end if;
   end process;

end architecture;
