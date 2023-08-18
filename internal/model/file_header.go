package model

import (
	"fmt"
)

type FileHeader struct {
	Crc32            uint32
	DecompressedSize uint32
	Unknown          uint32
}

func (fh *FileHeader) Dump() {
	fmt.Printf("File header\n")
	fmt.Printf("Crc32: %08x\n", fh.Crc32)
	fmt.Printf("DecompressedSize: %08x\n", fh.DecompressedSize)
	fmt.Printf("Unknown: %08x\n", fh.Unknown)
}
