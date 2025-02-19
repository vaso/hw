package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
)

const BuffSize = 10

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 || limit < 0 {
		return errors.New("negative offset or limit values")
	}

	// open input file
	inputFile, err := os.Open(fromPath)
	defer func() {
		if err := inputFile.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		return err
	}

	// get input file size
	fileInfo, err := inputFile.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}
	fileInfo.Mode()
	fileSize := fileInfo.Size()
	if fileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	// init FileSeeker
	inputFileSeeker := io.ReadSeeker(inputFile)
	_, err = inputFileSeeker.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	// init output file
	outputFile, err := os.Create(toPath)
	defer func() {
		if err := outputFile.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		return err
	}

	var restSize int64
	if limit > 0 {
		restSize = int64(math.Min(float64(fileSize-offset), float64(limit)))
	} else {
		restSize = fileSize - offset
	}

	wg := sync.WaitGroup{}

	progressChan := make(chan int64, 1)
	initPB(restSize, progressChan, &wg)

	copySize := restSize
	for copySize > 0 {
		length := math.Min(float64(BuffSize), float64(copySize))
		n, err := io.CopyN(outputFile, inputFileSeeker, int64(length))
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if n == 0 {
			break
		}
		copySize -= n
		if copySize < 0 {
			copySize = 0
		}
		updatePB(copySize, progressChan)
	}

	close(progressChan)
	wg.Wait()
	finishPB()
	fmt.Printf("\nsuccessfully copied file\n")

	return nil
}

func initPB(restSize int64, progressChan chan int64, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for restSizeVal := range progressChan {
			completedPercent := 100.0 * float64(restSize-restSizeVal) / float64(restSize)
			displayCurrentPBStatus(completedPercent)
		}
	}()
}

func updatePB(copySize int64, pbChan chan int64) {
	select {
	case pbChan <- copySize:
	default:
	}
}

func finishPB() {
	displayCurrentPBStatus(100.0)
}

func displayCurrentPBStatus(completedPercent float64) {
	fmt.Printf("\rcompleted - %.2f%%", completedPercent)
}
