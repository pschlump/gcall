
Calling a contract with a bytes32 parameter
	
	1.	You can pass 0x000... - a hex, or 0000... - a hex string, 64 chars long, or "abc" - a string 32 chars lon - where each char is converted to 0..255 value.
	2.  The data is left justified into a `[32]byte` array - zero filled.  Any extra is discarded without warning or error.

