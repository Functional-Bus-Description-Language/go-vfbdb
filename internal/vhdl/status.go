package vhdl

import (
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func generateStatus(st *fbdl.Status, fmts *EntityFormatters) {
	if st.IsArray {
		generateStatusArray(st, fmts)
	} else {
		generateStatusSingle(st, fmts)
	}
}

func generateStatusArray(st *fbdl.Status, fmts *EntityFormatters) {
	switch st.Access.(type) {
	case fbdl.AccessArrayMultiple:
		generateStatusArrayMultiple(st, fmts)
	default:
		panic("not yet implemented")
	}
}

func generateStatusSingle(st *fbdl.Status, fmts *EntityFormatters) {
	if st.Name != "x_uuid_x" && st.Name != "x_timestamp_x" {
		s := fmt.Sprintf(";\n   %s_i : in std_logic_vector(%d downto 0)", st.Name, st.Width-1)
		fmts.EntityFunctionalPorts += s
	}

	switch st.Access.(type) {
	case fbdl.AccessSingleSingle:
		generateStatusSingleSingle(st, fmts)
	default:
		panic("unknown single access strategy")
	}
}

func generateStatusSingleSingle(st *fbdl.Status, fmts *EntityFormatters) {
	fbdlAccess := st.Access.(fbdl.AccessSingleSingle)
	addr := fbdlAccess.Addr
	mask := fbdlAccess.Mask

	access := `
         %s : if internal_addr = %d then
            internal_master_in.dat(%d downto %d) <= registers(internal_addr)(%[3]d downto %[4]d);

            if internal_master_out.we = '0' then
               internal_master_in.ack <= '1';
               internal_master_in.err <= '0';
            end if;
         end if;
`
	access = fmt.Sprintf(access, st.Name, addr, mask.Upper, mask.Lower)
	fmts.StatusesAccess += access

	var routing string
	if st.Name == "x_uuid_x" || st.Name == "x_timestamp_x" {
		routing = fmt.Sprintf(
			"   registers(%d)(%d downto %d) <= %s;\n",
			addr, mask.Upper, mask.Lower, string(st.Default),
		)
	} else {
		routing = fmt.Sprintf(
			"   registers(%d)(%d downto %d) <= %s_i;\n",
			addr, mask.Upper, mask.Lower, st.Name,
		)
	}

	fmts.StatusesRouting += routing
}

func generateStatusArrayMultiple(st *fbdl.Status, fmts *EntityFormatters) {
	fbdlAccess := st.Access.(fbdl.AccessArrayMultiple)

	port := fmt.Sprintf(";\n   %s_i : in t_slv_vector(%d downto 0)(%d downto 0)", st.Name, st.Count-1, st.Width-1)
	fmts.EntityFunctionalPorts += port

	itemsPerAccess := fbdlAccess.ItemsPerAccess

	var access string
	if fbdlAccess.ItemCount <= itemsPerAccess {
		access = fmt.Sprintf(`
         %[1]s : if internal_addr = %[2]d then
            internal_master_in.dat(%[3]d downto %[4]d) <= registers(internal_addr)(%[3]d downto %[4]d);
            if internal_master_out.we = '0' then
               internal_master_in.ack <= '1';
               internal_master_in.err <= '0';
            end if;
         end if;
`,
			st.Name,
			fbdlAccess.StartAddr(),
			fbdlAccess.EndBit(),
			fbdlAccess.StartBit,
		)
	} else if fbdlAccess.ItemCount%itemsPerAccess == 0 {
		access = fmt.Sprintf(`
         %[1]s : if %[2]d <= internal_addr and internal_addr <= %[3]d then
            internal_master_in.dat(%[4]d downto %[5]d) <= registers(internal_addr)(%[4]d downto %[5]d);
            if internal_master_out.we = '0' then
               internal_master_in.ack <= '1';
               internal_master_in.err <= '0';
            end if;
         end if;
`,
			st.Name,
			fbdlAccess.StartAddr(),
			fbdlAccess.StartAddr()+fbdlAccess.Count()-1,
			fbdlAccess.EndBit(),
			fbdlAccess.StartBit,
		)
	} else {
		access = fmt.Sprintf(`
         %[1]s : if %[2]d <= internal_addr and internal_addr <= %[3]d then
            if internal_addr = %[3]d then
               internal_master_in.dat(%[4]d downto %[5]d) <= registers(internal_addr)(%[4]d downto %[5]d);
            else
               internal_master_in.dat(%[6]d downto %[5]d) <= registers(internal_addr)(%[6]d downto %[5]d);
            end if;
            if internal_master_out.we = '0' then
               internal_master_in.ack <= '1';
               internal_master_in.err <= '0';
            end if;
         end if;
`,
			st.Name,
			fbdlAccess.StartAddr(),
			fbdlAccess.StartAddr()+fbdlAccess.Count()-1,
			fbdlAccess.EndBit(),
			fbdlAccess.StartBit,
			fbdlAccess.ItemWidth*itemsPerAccess-1,
		)
	}

	fmts.StatusesAccess += access

	var routing string
	if fbdlAccess.ItemCount <= itemsPerAccess {
		routing = fmt.Sprintf(`
   %[1]s_registers : for reg in 0 to 0 loop
      %[1]s_items : for item in 0 to %[2]d loop
         registers(%[3]d + reg)(%[4]d * (item + 1) - 1 + %[5]d downto %[4]d * item + %[5]d) <= %[1]s_i(item);
      end loop;
   end loop;
`,
			st.Name,
			fbdlAccess.ItemCount-1,
			fbdlAccess.StartAddr(),
			fbdlAccess.ItemWidth,
			fbdlAccess.StartBit,
		)
	} else if fbdlAccess.ItemCount%itemsPerAccess == 0 {
		routing = fmt.Sprintf(`
   %[1]s_registers : for reg in 0 to %[2]d loop
      %[1]s_items : for item in 0 to  %[3]d loop
         registers(%[4]d + reg)(%[5]d * (item + 1) - 1 + %[6]d downto %[5]d * item + %[6]d) <= %[1]s_i(%[7]d * reg + item);
      end loop;
   end loop;
`,
			st.Name,
			fbdlAccess.Count()-1,
			itemsPerAccess-1,
			fbdlAccess.StartAddr(),
			fbdlAccess.ItemWidth,
			fbdlAccess.StartBit,
			itemsPerAccess,
		)
	} else {
		routing = fmt.Sprintf(`
   %[1]s_registers : for reg in 0 to %[2]d loop
      %[1]s_last_register : if reg = %[2]d then
         %[1]s_tail_items : for item in 0 to %[3]d loop
            registers(%[4]d + reg)(%[5]d * (item + 1) - 1 + %[6]d downto %[5]d * item + %[6]d) <= %[1]s_i(%[7]d * reg + item);
         end loop;
      else
         %[1]s_head_items : for item in 0 to %[8]d loop
            registers(%[4]d + reg)(%[5]d * (item + 1) - 1 + %[6]d downto %[5]d * item + %[6]d) <= %[1]s_i(%[7]d * reg + item);
         end loop;
      end if;
   end loop;
`,
			st.Name,
			fbdlAccess.Count()-1,
			st.Count%itemsPerAccess-1,
			fbdlAccess.StartAddr(),
			fbdlAccess.ItemWidth,
			fbdlAccess.StartBit,
			itemsPerAccess,
			itemsPerAccess-1,
		)
	}

	fmts.StatusesRouting += routing
}
