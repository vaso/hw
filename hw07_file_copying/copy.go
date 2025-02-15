package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	buffSize := 1024

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

	//// check if output file exists
	// existingToFileInfo, err := os.Stat(toPath)
	// if err == nil && existingToFileInfo.Size() > 0 {
	//	 log.Fatalf("error writing to existing file: %s", toPath)
	// }

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
	progressChan := initPB(restSize, &wg)

	copySize := restSize
	for copySize > 0 {
		length := math.Min(float64(buffSize), float64(copySize))
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

func initPB(restSize int64, wg *sync.WaitGroup) chan int64 {
	// totalStagesCount := int64(40)
	progressChan := make(chan int64, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// curUndoneStagesCount := int64(0)
		for restSizeVal := range progressChan {
			completedPercent := 100.0 * float64(restSize-restSizeVal) / float64(restSize)
			fmt.Printf("\rcompleted - %.2f%%", completedPercent)

			// undoneStagesCount := int64(math.Ceil((float64(restSizeVal) / float64(restSize)) * float64(totalStagesCount)))
			// if curUndoneStagesCount == undoneStagesCount {
			//   continue
			// }
			// curUndoneStagesCount = undoneStagesCount
			// doneStagesCount := totalStagesCount - curUndoneStagesCount
			// doneStr := strings.Repeat("+", int(doneStagesCount))
			// undoneStr := strings.Repeat("_", int(undoneStagesCount))
			//			fmt.Printf("\rstatus: [%s%s]", doneStr, undoneStr)
			// fmt.Printf("status: [done - %d, undone - %d]\n", doneStagesCount, curUndoneStagesCount)
			//			fmt.Printf("status: [undone - %d]\n", restSizeVal)
		}
	}()

	return progressChan
}

func updatePB(copySize int64, pbChan chan int64) {
	select {
	case pbChan <- copySize:
	default:
	}
}

func finishPB() {
	fmt.Printf("\rcompleted - %.2f%%", 100.0)
}
