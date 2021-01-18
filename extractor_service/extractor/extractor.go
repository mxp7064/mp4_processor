package exctractor

import (
	"encoding/binary"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// convert a slice form byte array to int
func convertByteSliceToInt(mp4Bytes []byte, startIndex int, endIndex int) int {
	return int(binary.BigEndian.Uint32(mp4Bytes[startIndex:endIndex]))
}

// convert a slice from byte array to string
func convertByteSliceToString(mp4Bytes []byte, startIndex int, endIndex int) string {
	return string(mp4Bytes[startIndex:endIndex])
}

// get initialization segment bytes from mp4 file from the provided file path
func getInitSegment(mp4Path string) ([]byte, error) {
	mp4Bytes, err := ioutil.ReadFile(mp4Path)
	if err != nil {
		return nil, err
	}

	ftypBoxSize := convertByteSliceToInt(mp4Bytes, 0, 4)
	ftypBoxType := convertByteSliceToString(mp4Bytes, 4, 8)
	if ftypBoxType != "ftyp" {
		return nil, errors.New("ftyp box not found")
	}

	ftypBoxBytes := mp4Bytes[0:ftypBoxSize]

	moovBoxBytes, err := getMoovBox(mp4Bytes, ftypBoxSize)
	if err != nil {
		return nil, err
	}

	return append(ftypBoxBytes, moovBoxBytes...), nil
}

/*
get moov box bytes from the provided mp4 byte array
starts searching for the moov box from boxStartIndex
*/
func getMoovBox(mp4Bytes []byte, boxStartIndex int) ([]byte, error) {
	for {

		if boxStartIndex+8 > len(mp4Bytes) {
			break
		}

		boxSize := convertByteSliceToInt(mp4Bytes, boxStartIndex, boxStartIndex+4)
		boxType := convertByteSliceToString(mp4Bytes, boxStartIndex+4, boxStartIndex+8)

		if boxType == "moov" {
			return mp4Bytes[boxStartIndex : boxStartIndex+boxSize], nil
		}

		boxStartIndex += int(boxSize) + boxStartIndex
	}

	return nil, errors.New("moov box not found")
}

type InitSegmentExtractor interface {
	ExtractInitSegment(mp4Path string) (string, error)
}

type InitSegmentExtractorService struct {
}

/*
InitSegmentExtractorService implements interface InitSegmentExtractor
extracts the initialization segment from mp4 file from the provided file path
writes the initialization segment to a new file and returns its path
*/
func (i InitSegmentExtractorService) ExtractInitSegment(mp4Path string) (string, error) {
	initSegmentBytes, err := getInitSegment(mp4Path)
	if err != nil {
		return "", errors.Wrap(err, "Extracting initital segment from file at path '"+mp4Path+"' failed")
	} else {
		initSegmentFilePath := filepath.Join(filepath.Dir(mp4Path), "init-segment-"+uuid.NewV4().String())
		err := ioutil.WriteFile(initSegmentFilePath, initSegmentBytes, 0644)
		if err != nil {
			return "", errors.Wrap(err, "Writing initialization segment to a file path '"+initSegmentFilePath+"' failed")
		}
		return initSegmentFilePath, nil
	}
}
