package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jfreymuth/oggvorbis"
)

func loadOggFloat32(path string) ([]float32, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r, err := oggvorbis.NewReader(file)
	// handle error

	fmt.Println(r.SampleRate())
	fmt.Println(r.Channels())

	buffer := make([]float32, 8192)
	for {
		//n, err := r.Read(buffer)

		// use buffer[:n]
		//return buffer[:n], nil

		if err == io.EOF {
			break
		}
		if err != nil {
			// handle error
			return nil, err
		}
	}
	return buffer, nil
}

func loadOggByte(path string) ([]byte, error) {

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	/*
		bytesread, err := file.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	*/

	return buffer, nil

}
