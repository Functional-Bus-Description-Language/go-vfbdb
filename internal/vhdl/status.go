package vhdl

import (
	_ "embed"
	"fmt"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
)

func generateStatus(status *fbdl.Status, fmts *EntityFormatters) {
	if status.IsArray {
		generateStatusArray(status, fmts)
	} else {
		generateStatusSingle(status, fmts)
	}
}

func generateStatusArray(status *fbdl.Status, fmts *EntityFormatters) {
	switch status.Access.(type) {
	case fbdl.AccessArrayMultiple:
		generateStatusArrayMultiple(status, fmts)
	default:
		panic("not yet implemented")
	}
}

func generateStatusSingle(status *fbdl.Status, fmts *EntityFormatters) {
	if status.Name != "x_uuid_x" && status.Name != "x_timestamp_x" {
		s := fmt.Sprintf(";\n   %s_i : in std_logic_vector(%d - 1 downto 0)", status.Name, status.Width)
		fmts.EntityFunctionalPorts += s
	}

	switch status.Access.(type) {
	case fbdl.AccessSingleSingle:
		generateStatusSingleSingle(status, fmts)
	default:
		panic("unknown single access strategy")
	}
}

func generateStatusSingleSingle(status *fbdl.Status, fmts *EntityFormatters) {
	fbdlAccess := status.Access.(fbdl.AccessSingleSingle)
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
	access = fmt.Sprintf(access, status.Name, addr, mask.Upper, mask.Lower)
	fmts.StatusesAccess += access

	var routing string
	if status.Name == "x_uuid_x" || status.Name == "x_timestamp_x" {
		routing = fmt.Sprintf(
			"   registers(%d)(%d downto %d) <= %s;\n",
			addr, mask.Upper, mask.Lower, string(status.Default),
		)
	} else {
		routing = fmt.Sprintf(
			"   registers(%d)(%d downto %d) <= %s_i(%[2]d downto %[3]d);\n",
			addr, mask.Upper, mask.Lower, status.Name,
		)
	}

	fmts.StatusesRouting += routing
}

func generateStatusArrayMultiple(status *fbdl.Status, fmts *EntityFormatters) {
	fbdlAccess := status.Access.(fbdl.AccessArrayMultiple)

	port := fmt.Sprintf(";\n   %s_i : in t_slv_vector(%d downto 0)(%d downto 0)", status.Name, status.Count-1, status.Width-1)
	fmts.EntityFunctionalPorts += port

	itemsPerAccess := busWidth / fbdlAccess.ItemWidth

	var access string
	if fbdlAccess.ItemCount%itemsPerAccess == 0 {
		access = fmt.Sprintf(`
         %s : if %d <= internal_addr and internal_addr <= %d then
            internal_master_in.dat(%d downto 0) <= registers(internal_addr)(%[4]d downto 0);
            if internal_master_out.we = '0' then
               internal_master_in.ack <= '1';
               internal_master_in.err <= '0';
            end if;
         end if;
`,
			status.Name,
			fbdlAccess.StartAddr(), fbdlAccess.StartAddr()+fbdlAccess.Count()-1,
			fbdlAccess.ItemWidth*itemsPerAccess-1,
		)
	} else if fbdlAccess.ItemCount < itemsPerAccess {
		access = fmt.Sprintf(`
         %s : if internal_addr = %d then
            internal_master_in.dat(%d downto 0) <= registers(internal_addr)(%[3]d downto 0);
            if internal_master_out.we = '0' then
               internal_master_in.ack <= '1';
               internal_master_in.err <= '0';
            end if;
         end if;
`,
			status.Name, fbdlAccess.StartAddr(), fbdlAccess.ItemWidth*fbdlAccess.ItemCount-1,
		)
	} else {
		access = fmt.Sprintf(`
         %s : if %d <= internal_addr and internal_addr <= %d then
            if internal_addr = %[3]d then
               internal_master_in.dat(%d downto 0) <= registers(internal_addr)(%[4]d downto 0);
            else
               internal_master_in.dat(%d downto 0) <= registers(internal_addr)(%[5]d downto 0);
            end if;
            if internal_master_out.we = '0' then
               internal_master_in.ack <= '1';
               internal_master_in.err <= '0';
            end if;
         end if;
`,
			status.Name, fbdlAccess.StartAddr(), fbdlAccess.StartAddr()+fbdlAccess.Count()-1,
			fbdlAccess.ItemWidth*(fbdlAccess.ItemCount%itemsPerAccess)-1, fbdlAccess.ItemWidth*itemsPerAccess-1,
		)
	}

	fmts.StatusesAccess += access

	var routing string
	if fbdlAccess.ItemCount%itemsPerAccess == 0 {
		routing = fmt.Sprintf(`
   %[1]s_registers : for reg in 0 to %[2]d loop
      %[1]s_items : for item in 0 to  %[3]d loop
         registers(%[4]d + reg)(%[5]d * (item + 1) - 1 downto %[5]d * item) <= %[1]s_i(%[6]d * reg + item);
      end loop;
   end loop;
`,
			status.Name,
			fbdlAccess.Count()-1,
			itemsPerAccess-1,
			fbdlAccess.StartAddr(),
			fbdlAccess.ItemWidth,
			itemsPerAccess,
		)
	} else if fbdlAccess.ItemCount < itemsPerAccess {
		routing = fmt.Sprintf(`
   %[1]s_registers : for reg in 0 to %[2]d loop
      %[1]s_items : for item in 0 to %[3]d loop
         registers(%[4]d + reg)(%[5]d * (item + 1) - 1 downto %[5]d * item) <= %[1]s_i(item);
      end loop;
   end loop;
`,
			status.Name,
			fbdlAccess.Count()-1,
			fbdlAccess.ItemCount-1,
			fbdlAccess.StartAddr(),
			fbdlAccess.ItemWidth,
		)
	} else {
		routing = fmt.Sprintf(`
   %[1]s_registers : for reg in 0 to %[2]d loop
      %[1]s_last_register : if reg = %[2]d then
         %[1]s_tail_items : for item in 0 to %[3]d loop
            registers(%[4]d + reg)(%[5]d * (item + 1) - 1 downto %[5]d * item) <= %[1]s_i(%[6]d * reg + item);
         end loop;
      else
         %[1]s_head_items : for item in 0 to %[7]d loop
            registers(%[4]d + reg)(%[5]d * (item + 1) - 1 downto %[5]d * item) <= %[1]s_i(%[6]d * reg + item);
         end loop;
      end if;
   end loop;
`,
			status.Name,
			fbdlAccess.Count()-1,
			status.Count%itemsPerAccess-1,
			fbdlAccess.StartAddr(),
			fbdlAccess.ItemWidth,
			itemsPerAccess,
			itemsPerAccess-1,
		)
	}

	fmts.StatusesRouting += routing
}
