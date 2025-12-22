//go:build ignore

package main

import (
	"fmt"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func sort2BlocksOf256Bits(
	pa, pb, tmp VecVirtual,
	VPMIN, VPMAX func(...Op)) {
	VPMIN(pa, pb, tmp)
	VPMAX(pa, pb, pb)
	VMOVDQU(tmp, pa)
}

func sort3BlocksOf256Bits(
	wordSizeBytes uint8,
	pa, pb, pc, i Register,
	MM func() VecVirtual,
	VPMIN, VPMAX func(...Op)) {
	va, vb, vc, tmp := MM(), MM(), MM(), MM()

	VMOVDQU(Mem{Base: pa, Index: i, Scale: wordSizeBytes}, va)
	VMOVDQU(Mem{Base: pb, Index: i, Scale: wordSizeBytes}, vb)
	VMOVDQU(Mem{Base: pc, Index: i, Scale: wordSizeBytes}, vc)

	sort2BlocksOf256Bits(vb, vc, tmp, VPMIN, VPMAX)
	sort2BlocksOf256Bits(va, vb, tmp, VPMIN, VPMAX)

	VMOVDQU(va, Mem{Base: pa, Index: i, Scale: wordSizeBytes})
	VMOVDQU(vb, Mem{Base: pb, Index: i, Scale: wordSizeBytes})
	VMOVDQU(vc, Mem{Base: pc, Index: i, Scale: wordSizeBytes})
}

func generateBlockSortUpward(
	nElements uint64, wordSizeBytes uint8,
	MM func() VecVirtual,
	VPMIN, VPMAX func(...Op)) {
	TEXT(
		fmt.Sprintf("blockSortUpwardBy%dElementsOf%dBytes", nElements, wordSizeBytes),
		NOSPLIT,
		fmt.Sprintf("func(T []uint%d, i, gap int) int", 8*wordSizeBytes))
	pa := Load(Param("T").Base(), GP64())
	TLen := Load(Param("T").Len(), GP64())
	i := Load(Param("i"), GP64())
	gap := Load(Param("gap"), GP64())

	pb, pc := GP64(), GP64()
	LEAQ(Mem{Base: pa, Index: gap, Scale: wordSizeBytes}, pb)
	LEAQ(Mem{Base: pb, Index: gap, Scale: wordSizeBytes}, pc)
	SUBQ(Imm(nElements), TLen)

	JMP(LabelRef("before_end_of_loop_upward"))
	Label("loop_upward")
	sort3BlocksOf256Bits(wordSizeBytes, pa, pb, pc, i, MM, VPMIN, VPMAX)
	ADDQ(Imm(nElements), i)
	Label("before_end_of_loop_upward")
	CMPQ(i, TLen)
	JLE(LabelRef("loop_upward"))

	Store(i, ReturnIndex(0))
	RET()
}

func generateBlockSortDownward(
	nElements uint64, wordSizeBytes uint8,
	MM func() VecVirtual,
	VPMIN, VPMAX func(...Op)) {
	TEXT(
		fmt.Sprintf("blockSortDownwardBy%dElementsOf%dBytes", nElements, wordSizeBytes),
		NOSPLIT,
		fmt.Sprintf("func(T []uint%d, i, gap int) int", 8*wordSizeBytes))
	pa := Load(Param("T").Base(), GP64())
	i := Load(Param("i"), GP64())
	gap := Load(Param("gap"), GP64())

	pb, pc := GP64(), GP64()
	LEAQ(Mem{Base: pa, Index: gap, Scale: wordSizeBytes}, pb)
	LEAQ(Mem{Base: pb, Index: gap, Scale: wordSizeBytes}, pc)

	ADDQ(Imm(nElements), i)
	JMP(LabelRef("before_end_of_loop_downward"))
	Label("loop_downward")
	sort3BlocksOf256Bits(wordSizeBytes, pa, pb, pc, i, MM, VPMIN, VPMAX)
	Label("before_end_of_loop_downward")
	SUBQ(Imm(nElements), i)
	JGE(LabelRef("loop_downward"))

	Store(i, ReturnIndex(0))
	RET()
}

func main() {
	// uint32, sort by 8 unaligned words
	generateBlockSortUpward(8, 4, YMM, VPMINUD, VPMAXUD)
	generateBlockSortDownward(8, 4, YMM, VPMINUD, VPMAXUD)
	// uint32, sort by 4 unaligned words
	generateBlockSortUpward(4, 4, XMM, VPMINUD, VPMAXUD)
	generateBlockSortDownward(4, 4, XMM, VPMINUD, VPMAXUD)
	// uint64, sort by 4 unaligned double words
	generateBlockSortUpward(4, 8, YMM, VPMINUQ, VPMAXUQ)
	generateBlockSortDownward(4, 8, YMM, VPMINUQ, VPMAXUQ)
	Generate()
}
