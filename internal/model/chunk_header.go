package model

import (
	"fmt"
)

type ChunkHeader struct {
	HeaderTag         uint64
	ChunkSize         uint64
	CompressionType   byte
	CompressedSize1   uint64
	DecompressedSize1 uint64
	CompressedSize2   uint64
	DecompressedSize2 uint64
}

func (ch *ChunkHeader) Dump(chunkNumber int) {
	fmt.Printf("Chunk %d\n", chunkNumber)
	fmt.Printf("HeaderTag: %016x\n", ch.HeaderTag)
	fmt.Printf("ChunkSize: %016x\n", ch.ChunkSize)
	fmt.Printf("CompressionType: %x\n", ch.CompressionType)
	fmt.Printf("CompressionSize1: %016x\n", ch.CompressedSize1)
	fmt.Printf("DecompressionSize1: %016x\n", ch.DecompressedSize1)
	fmt.Printf("CompressionSize2: %016x\n", ch.CompressedSize2)
	fmt.Printf("DecompressionSize2: %016x\n", ch.DecompressedSize2)
}
