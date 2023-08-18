package core

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"math"
	"os"

	"github.com/AndrewSav/r2recompress/internal/model"
)

type chunkHeaderData struct {
	model.ChunkHeader
	chunkNumber               int
	previousDecompressedSize1 uint64
	previousDecompressedSize2 uint64
}

type warningContext struct {
	noWarnings          bool
	singleWarning       bool
	firstWarningPrinted bool
}

func printWarning(message string, context *warningContext) {
	if !(context.noWarnings || context.singleWarning && context.firstWarningPrinted) {
		fmt.Printf("%s. Recompression is unlikely to produce correct result!\n", message)
	}
	context.firstWarningPrinted = true
}

func fileHeaderChecks(fh model.FileHeader, totalDecompressedSize uint32, crc uint32, context *warningContext) {
	if crc != fh.Crc32 {
		message := fmt.Sprintf("wrong crc of decompressed data have %08x, want %08x", crc, fh.Crc32)
		printWarning(message, context)
	}
	if totalDecompressedSize+8 != fh.DecompressedSize {
		message := fmt.Sprintf("wrong total decompressed size have %016x, want %016x", totalDecompressedSize+8, fh.DecompressedSize)
		printWarning(message, context)
	}
	if fh.Unknown != model.Unknown {
		message := fmt.Sprintf("wrong file header unknown field value. have %d, want %d", fh.Unknown, model.Unknown)
		printWarning(message, context)
	}
}

func chunkHeaderChecks(chunkNumber int, ch chunkHeaderData, decompressedSize int64, context *warningContext) {
	if ch.HeaderTag != model.ARCHIVE_V2_HEADER_TAG {
		message := fmt.Sprintf("wrong chunk %d header HeaderTag field value. have %016x, want %016x", chunkNumber, ch.HeaderTag, model.ARCHIVE_V2_HEADER_TAG)
		printWarning(message, context)
	}
	if ch.ChunkSize != model.ChunkSize {
		message := fmt.Sprintf("wrong chunk %d header ChunkSize field value. have %016x, want %016x", chunkNumber, ch.ChunkSize, model.ChunkSize)
		printWarning(message, context)
	}
	if ch.CompressionType != model.CompressionType {
		message := fmt.Sprintf("wrong chunk %d header CompressionType field value. have %x, want %x", chunkNumber, ch.CompressionType, model.CompressionType)
		printWarning(message, context)
	}
	if ch.CompressedSize1 != ch.CompressedSize2 {
		message := fmt.Sprintf("wrong chunk %d header values. CompressedSize1 (%016x) should be equal to CompressedSize2 (%016x), but they are not", chunkNumber, ch.CompressedSize1, ch.CompressedSize2)
		printWarning(message, context)
	}
	if ch.DecompressedSize1 != ch.DecompressedSize2 {
		message := fmt.Sprintf("wrong chunk %d header values. DecompressedSize1 (%016x) should be equal to DecompressedSize2 (%016x), but they are not", chunkNumber, ch.DecompressedSize1, ch.DecompressedSize2)
		printWarning(message, context)
	}
	if chunkNumber > 1 && ch.previousDecompressedSize1 != model.ChunkSize {
		message := fmt.Sprintf("wrong chunk %d header DecompressedSize1 field value. have %016x, want %016x", chunkNumber-1, ch.previousDecompressedSize1, model.ChunkSize)
		printWarning(message, context)
	}
	if chunkNumber > 1 && ch.previousDecompressedSize2 != model.ChunkSize {
		message := fmt.Sprintf("wrong chunk %d header DecompressedSize2 field value. have %016x, want %016x", chunkNumber-1, ch.previousDecompressedSize2, model.ChunkSize)
		printWarning(message, context)
	}
	if decompressedSize > math.MaxUint32 {
		message := fmt.Sprintf("chunk %d decompressed size more than the threshold have %d, want < %d", chunkNumber, decompressedSize, ch.DecompressedSize1)
		printWarning(message, context)
	}
	if ch.DecompressedSize1 != uint64(decompressedSize) {
		message := fmt.Sprintf("wrong chunk %d decompressed size have %016x, want %016x", chunkNumber, decompressedSize, ch.DecompressedSize1)
		printWarning(message, context)
	}

}

func Decompress(input, output string, options model.Options) error {

	inputFile, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}

	inputReader := bufio.NewReader(inputFile)
	var fh model.FileHeader
	err = binary.Read(inputReader, binary.LittleEndian, &fh)
	if err != nil {
		return fmt.Errorf("error reading file header: %v", err)
	}

	if *options.Verbose {
		fh.Dump()
		fmt.Println()
	}

	outputFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("error opening output file: %v", err)
	}

	hash := crc32.NewIEEE()
	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, &fh.DecompressedSize)
	binary.Write(&b, binary.LittleEndian, &fh.Unknown)
	_, err = io.Copy(hash, &b)
	if err != nil {
		return fmt.Errorf("error hashing preamble: %v", err)
	}

	chunksCount := 0
	totalDecompressedSize := uint32(0)
	var previousDecompressedSize1, previousDecompressedSize2 uint64
	context := warningContext{
		noWarnings:    *options.NoWarnings,
		singleWarning: *options.SingleWarning,
	}

	for {
		chunksCount++
		var ch chunkHeaderData

		err = binary.Read(inputReader, binary.LittleEndian, &ch.ChunkHeader)
		if err != nil {
			return fmt.Errorf("error reading chunk header: %v", err)
		}
		if *options.Verbose {
			ch.Dump(chunksCount)
		}
		ch.previousDecompressedSize1, ch.previousDecompressedSize2 = previousDecompressedSize1, previousDecompressedSize2
		previousDecompressedSize1, previousDecompressedSize2 = ch.DecompressedSize1, ch.DecompressedSize2
		r, err := zlib.NewReader(inputReader)
		if err != nil {
			return fmt.Errorf("error creating zlib reader: %v", err)
		}
		decompressedSize := int64(0)
		if chunksCount == 1 {
			sz, err := io.CopyN(outputFile, r, 4)
			if err != nil {
				return fmt.Errorf("error decompressing preamble: %v", err)
			}
			decompressedSize += sz
		}
		mw := io.MultiWriter(outputFile, hash)
		sz, err := io.Copy(mw, r)
		if err != nil {
			return fmt.Errorf("error decompressing: %v", err)
		}
		decompressedSize += sz
		if *options.Verbose {
			fmt.Printf("Real decompressed size: %016x (%d)\n", decompressedSize, decompressedSize)
		}
		totalDecompressedSize += uint32(decompressedSize)
		if *options.Verbose {
			fmt.Println()
		}
		chunkHeaderChecks(chunksCount, ch, decompressedSize, &context)
		if _, err := inputReader.Peek(1); err == io.EOF {
			break
		}
	}

	crc := hash.Sum32()
	if *options.Verbose {
		fmt.Printf("Crc32 of decompressed data: %08x\n", crc)
	}
	fileHeaderChecks(fh, totalDecompressedSize, crc, &context)
	return nil
}
