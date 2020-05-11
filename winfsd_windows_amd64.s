#include "go_asm.h"
#include "textflag.h"

// createFileWTrampoline adds FILE_SHARE_DELETE to CreateFileW's dwShareMode
// parameter and then jumps to createFileWAddr, which holds the address of
// kernel32!CreateFileW.
TEXT ·createFileWTrampoline(SB),NOSPLIT|NOFRAME,$0
	ORQ	$4, R8		// 4 is FILE_SHARE_DELETE; R8 holds dwShareMode.
	MOVQ	·createFileWAddr(SB), AX
	JMP	AX
