package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func extractBytesToIntSegment(mp4Bytes []byte, startIndex int, endIndex int) int {
	return int(binary.BigEndian.Uint32(mp4Bytes[startIndex:endIndex]))
}

func extractBytesToStringSegment(mp4Bytes []byte, startIndex int, endIndex int) string {
	return string(mp4Bytes[startIndex:endIndex])
}

func getInitSegmentPath(mp4Path string) (string, error) {
	initSegmentBytes, err := extractInitSegment(mp4Path)
	if err != nil {
		return "", errors.Wrap(err, "Extracting initital segment from file at path '"+mp4Path+"' failed")
	} else {
		initSegmentFilePath := filepath.Join(filepath.Dir(mp4Path), uuid.NewV4().String())
		err := ioutil.WriteFile(initSegmentFilePath, initSegmentBytes, 0644)
		if err != nil {
			return "", errors.Wrap(err, "Writing initialization segment to a file path '"+initSegmentFilePath+"'  failed")
		}
		return initSegmentFilePath, nil
	}
}

func extractInitSegment(mp4Path string) ([]byte, error) {
	fmt.Printf("Processing video from path %v \n", mp4Path)

	mp4Bytes, err := ioutil.ReadFile(mp4Path)
	if err != nil {
		return nil, err
	}

	ftypBoxSize := extractBytesToIntSegment(mp4Bytes, 0, 4)
	ftypBoxType := extractBytesToStringSegment(mp4Bytes, 4, 8)
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

func getMoovBox(mp4Bytes []byte, boxStartIndex int) ([]byte, error) {
	for {

		if boxStartIndex+8 > len(mp4Bytes) {
			break
		}

		boxSize := extractBytesToIntSegment(mp4Bytes, boxStartIndex, boxStartIndex+4)
		boxType := extractBytesToStringSegment(mp4Bytes, boxStartIndex+4, boxStartIndex+8)

		if boxType == "moov" {
			return mp4Bytes[boxStartIndex : boxStartIndex+boxSize], nil
		}

		boxStartIndex += int(boxSize) + boxStartIndex
	}

	return nil, errors.New("moov box not found")
}
