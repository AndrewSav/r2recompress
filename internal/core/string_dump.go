package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/AndrewSav/r2recompress/internal/model"
)

func StringDump(input, output string, options model.Options) error {
	inputFile, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}
	outputFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("error opening output file: %v", err)
	}

	info, err := inputFile.Stat()
	if err != nil {
		return fmt.Errorf("error calling stat on input file: %v", err)
	}

	fileSize := info.Size()

stream:
	for {
		pos, err := inputFile.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("error getting position in the input file: %v", err)
		}
		b := make([]byte, model.ChunkSize)
		szint, err := inputFile.Read(b)
		sz := int64(szint)
		if err == io.EOF || err == nil {
		byte:
			for i := int64(0); i < sz-4; i++ {
				r := bytes.NewBuffer(b[i : i+4])
				var len32 int32
				err = binary.Read(r, binary.LittleEndian, &len32)
				l := int64(len32)
				if err != nil {
					return fmt.Errorf("error reading input buffer: %v", err)
				}
				if l > 2 && pos+i+4+l > 0 && pos+i+4+l <= fileSize && l < model.ChunkSize {
					if i+4+l >= model.ChunkSize {
						if i == 0 {
							continue byte
						}
						_, err := inputFile.Seek(pos+i, io.SeekStart)
						if err != nil {
							return fmt.Errorf("error setting position in the input file to the start of a string: %v", err)
						}
						continue stream
					} else {
						if b[i+4+l-1] == 0 {
							complaint := true
							for j := i + 4; j < i+4+l-1; j++ {
								if b[j] >= 0x7f || b[j] <= 0x1f {
									complaint = false
									break
								}
							}
							if complaint {
								outputFile.Write(b[i+4 : i+3+l])
								outputFile.Write([]byte{0xa})
							}
						}
					}
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading input: %v", err)
		}
	}

	return nil
}
