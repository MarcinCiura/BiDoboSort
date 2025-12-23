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

func initializeVariables(
	wordSizeBytes uint8) (
	pa Register, pb Register, pc Register, end Register) {
	pa = Load(Param("T").Base(), GP64())
	pb, pc = GP64(), GP64()
	gap := Load(Param("gap"), GP64())
	LEAQ(Mem{Base: pa, Index: gap, Scale: wordSizeBytes}, pb)
	LEAQ(Mem{Base: pb, Index: gap, Scale: wordSizeBytes}, pc)
	end = Load(Param("end"), GP64())
	return
}

func generateBlockSortUpward(
	nElements uint64, wordSizeBytes uint8,
	MM func() VecVirtual,
	VPMIN, VPMAX func(...Op)) {
	TEXT(
		fmt.Sprintf("blockSortUpwardBy%dElementsOf%dBytes", nElements, wordSizeBytes),
		NOSPLIT,
		fmt.Sprintf("func(T []uint%d, end, gap int) int", 8*wordSizeBytes))
	pa, pb, pc, end := initializeVariables(wordSizeBytes)
	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("before_end_of_loop"))

	Label("loop")
	sort3BlocksOf256Bits(wordSizeBytes, pa, pb, pc, i, MM, VPMIN, VPMAX)
	ADDQ(Imm(nElements), i)

	Label("before_end_of_loop")
	CMPQ(i, end)
	JL(LabelRef("loop"))

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
		fmt.Sprintf("func(T []uint%d, end, gap int) int", 8*wordSizeBytes))
	pa, pb, pc, i := initializeVariables(wordSizeBytes)
	JMP(LabelRef("before_end_of_loop"))

	Label("loop")
	sort3BlocksOf256Bits(wordSizeBytes, pa, pb, pc, i, MM, VPMIN, VPMAX)

	Label("before_end_of_loop")
	SUBQ(Imm(nElements), i)
	JGE(LabelRef("loop"))

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
