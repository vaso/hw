package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sync"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	buffSize := 1024

	flag.Parse()
	//fmt.Println("params: ", from, to, limit, offset)

	// open input file
	inputFile, err := os.Open(from)
	if err != nil {
		log.Fatalf("error opening source file: %v", err)
	}
	defer func() {
		if err := inputFile.Close(); err != nil {
			panic(err)
		}
	}()

	// get input file size
	fileInfo, err := inputFile.Stat()
	if err != nil {
		log.Fatalf("error getting source file info: %v", err)
	}
	fileInfo.Mode()
	fileSize := fileInfo.Size()
	if fileSize < offset {
		log.Fatalf("offset is bigger than input file size")
	}

	//init FileSeeker
	inputFileSeeker := io.ReadSeeker(inputFile)
	_, err = inputFileSeeker.Seek(offset, io.SeekStart)
	if err != nil {
		return
	}

	//// check if output file exists
	//existingToFileInfo, err := os.Stat(to)
	//if err == nil && existingToFileInfo.Size() > 0 {
	//	log.Fatalf("error writing to existing file: %s", to)
	//}

	// init output file
	outputFile, err := os.Create(to)
	if err != nil {
		log.Fatalf("error creating output file: %v", err)
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			panic(err)
		}
	}()

	var restSize int64
	if limit > 0 {
		restSize = int64(math.Min(float64(fileSize-offset), float64(limit)))
	} else {
		restSize = fileSize - offset
	}
	//fmt.Printf("copy %d bytes from file with original size %d bytes with offset %d bytes by chunks of %d bytes\n", restSize, fileSize, offset, buffSize)

	wg := sync.WaitGroup{}
	progressChan := initPB(restSize, &wg)

	copySize := restSize
	for copySize > 0 {
		length := math.Min(float64(buffSize), float64(copySize))
		n, err := io.CopyN(outputFile, inputFileSeeker, int64(length))
		if err != nil && err != io.EOF {
			log.Fatalf("error reading source file: %v", err)
		}
		if n == 0 {
			break
		}
		copySize -= n
		if copySize < 0 {
			copySize = 0
		}
		updatePB(copySize, progressChan)
		//time.Sleep(2 * time.Millisecond)
	}
	close(progressChan)
	wg.Wait()
	finishPB()
	fmt.Printf("\nsuccessfully copied file\n")
}

func initPB(restSize int64, wg *sync.WaitGroup) chan int64 {
	//totalStagesCount := int64(40)
	progressChan := make(chan int64, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		//fmt.Printf("\n")
		//curUndoneStagesCount := int64(0)
		for restSizeVal := range progressChan {
			completedPercent := 100.0 * float64(restSize-restSizeVal) / float64(restSize)
			fmt.Printf("\rcompleted - %.2f%%", completedPercent)

			//undoneStagesCount := int64(math.Ceil((float64(restSizeVal) / float64(restSize)) * float64(totalStagesCount)))
			//if curUndoneStagesCount == undoneStagesCount {
			//	continue
			//}
			//curUndoneStagesCount = undoneStagesCount
			//doneStagesCount := totalStagesCount - curUndoneStagesCount
			//doneStr := strings.Repeat("+", int(doneStagesCount))
			//undoneStr := strings.Repeat("_", int(undoneStagesCount))
			//			fmt.Printf("\rstatus: [%s%s]", doneStr, undoneStr)
			//fmt.Printf("status: [done - %d, undone - %d]\n", doneStagesCount, curUndoneStagesCount)
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
