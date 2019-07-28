
TEXT ·_fibon_if_fc_alloc(SB), 0, $65536-8
    MOVQ $0, ret+0(FP)
    RET

TEXT ·_fibon_if_fastcall(SB), 0, $65536-32
    MOVQ ptr+8(FP), CX
    MOVQ size+16(FP), DX
    MOVQ trap+0(FP), AX
    MOVQ SP, BX
    ADDQ $65472, SP // SP - 4 * 8
    ANDQ $~15, SP
    CALL AX
    MOVQ BX, SP
    MOVQ AX, ret+24(FP)
    RET

