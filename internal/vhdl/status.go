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

}

func generateStatusSingle(status *fbdl.Status, fmts *EntityFormatters) {
	if status.Name != "x_uuid_x" && status.Name != "x_timestamp_x" {
		s := fmt.Sprintf(";\n      %s_i : in std_logic_vector(%d - 1 downto 0)", status.Name, status.Width)
		fmts.EntityFunctionalPorts = s
	}

	fbdlAccess := status.Access.(*fbdl.AccessSingle)
	strategy := fbdlAccess.Strategy
	if strategy == "Single" {
		generateStatusSingleSingle(status, fmts)
	}
}

func generateStatusSingleSingle(status *fbdl.Status, fmts *EntityFormatters) {
	fbdlAccess := status.Access.(*fbdl.AccessSingle)
	addr := fbdlAccess.Address
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
