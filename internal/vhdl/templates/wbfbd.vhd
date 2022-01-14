library ieee;
   use ieee.std_logic_1164.all;
   use ieee.numeric_std.all;


package wbfbd is

type t_slv_vector is array (natural range <>) of std_logic_vector;

subtype int64 is signed(63 downto 0);
type int64_vector is array (natural range <>) of int64;
-- int converts int64 to integer type.
-- One can use to_integer() directly, but int() is much shorter.
-- Multiple to_integer() in the same line significatly extend the line.
function int(i64 : int64) return integer;

-- Packages constants
{{.PkgsConsts}}
end package;

package body wbfbd is

function int(i64 : int64) return integer is
begin
   return to_integer(i64);
end function;

end package body;
