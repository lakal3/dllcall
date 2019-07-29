
TEXT ·_if_fastcall(SB), 0, $65536-32
    MOVQ ptr+8(FP), DI
    MOVQ size+16(FP), SI
    MOVQ trap+0(FP), AX
    ADDQ $65504, SP // SP - 4 * 8
    CALL AX
    ADDQ $-65504, SP
    MOVQ AX, ret+24(FP)
    RET

