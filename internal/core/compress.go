package core

import (
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

func Compress(input, output string, options model.Options) error {
	inputFile, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}
	info, err := inputFile.Stat()
	if err != nil {
		return fmt.Errorf("error calling stat on input file: %v", err)
	}

	if info.Size()+8 > math.MaxUint32 {
		return fmt.Errorf("error input file is too large. size %d, should be less than %d", info.Size(), math.MaxUint32-8)
	}

	fh := model.FileHeader{}
	fh.DecompressedSize = uint32(info.Size() + 8)
	fh.Unknown = model.Unknown

	h := crc32.NewIEEE()

	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, &fh.DecompressedSize)
	binary.Write(&b, binary.LittleEndian, &fh.Unknown)

	_, err = io.Copy(h, &b)
	if err != nil {
		return fmt.Errorf("error hashing preamble: %v", err)
	}

	inputFile.Seek(4, io.SeekStart)
	_, err = io.Copy(h, inputFile)
	if err != nil {
		return fmt.Errorf("error hashing input file: %v", err)
	}

	fh.Crc32 = h.Sum32()

	inputFile.Seek(0, io.SeekStart)

	outputFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("error opening output file: %v", err)
	}

	err = binary.Write(outputFile, binary.LittleEndian, fh)
	if err != nil {
		return fmt.Errorf("could not write file header: %v", err)
	}

	if *options.Verbose {
		fh.Dump()
		fmt.Println()
	}

	for chunksCount := 1; ; chunksCount++ {
		ch := model.ChunkHeader{}
		ch.HeaderTag = model.ARCHIVE_V2_HEADER_TAG
		ch.ChunkSize = model.ChunkSize
		ch.CompressionType = model.CompressionType

		err := binary.Write(outputFile, binary.LittleEndian, ch)
		if err != nil {
			return fmt.Errorf("could not write initial chunk header: %v", err)
		}

		chunkBodyStart, err := outputFile.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("error getting chunk body start position in the output file: %v", err)
		}

		w := zlib.NewWriter(outputFile)
		sz, err := io.CopyN(w, inputFile, model.ChunkSize)
		w.Close()

		if err == io.EOF || err == nil {
			pos, err := outputFile.Seek(0, io.SeekCurrent)
			if err != nil {
				return fmt.Errorf("error getting chunk body end position in the output file: %v", err)
			}
			compressedSize := pos - chunkBodyStart
			ch.CompressedSize1 = uint64(compressedSize)
			ch.CompressedSize2 = uint64(compressedSize)
			ch.DecompressedSize1 = uint64(sz)
			ch.DecompressedSize2 = uint64(sz)
			_, err = outputFile.Seek(chunkBodyStart-int64(binary.Size(ch)), io.SeekStart)
			if err != nil {
				return fmt.Errorf("error moving output file position to write chunk header: %v", err)
			}
			err = binary.Write(outputFile, binary.LittleEndian, ch)
			if err != nil {
				return fmt.Errorf("could not write final chunk header: %v", err)
			}
			_, err = outputFile.Seek(pos, io.SeekStart)
			if err != nil {
				return fmt.Errorf("error moving output file position after writing chunk header: %v", err)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error compressing: %v", err)
		}
		if *options.Verbose {
			ch.Dump(chunksCount)
			fmt.Println()
		}

	}

	return nil
}
