package main

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/h2non/bimg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Option struct {
	VipsOptions    bimg.Options
	FileNameSuffix string
}

func main() {

	options := []Option{
		{
			VipsOptions: bimg.Options{
				Width:         800,
				Quality:       90,
				Interlace:     true,
				StripMetadata: true,
			},
			FileNameSuffix: ".jpg",
		},
		{
			VipsOptions: bimg.Options{
				Width:     800,
				Quality:   90,
				Interlace: true,
			},
			FileNameSuffix: ".webp",
		},
		{
			VipsOptions: bimg.Options{
				Width:     1600,
				Quality:   90,
				Interlace: true,
			},
			FileNameSuffix: "@2x.jpg",
		},
		{
			VipsOptions: bimg.Options{
				Width:     1600,
				Quality:   90,
				Interlace: true,
			},
			FileNameSuffix: "@2x.webp",
		},
	}

	files, err := ioutil.ReadDir("/mount")
	if err != nil {
		log.Fatal("something went wrong: ", err)
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		name := strings.TrimSuffix(file.Name(), ext)
		cleanName := slug.Make(name)

		if !file.IsDir() && (ext == ".jpg" || ext == ".png") {

			buffer, err := bimg.Read("/mount/" + file.Name())
			if err != nil {
				fmt.Printf("error reading %s ::: %s\n", file.Name(), err.Error())
				continue
			}

			for _, option := range options {
				newImg, err := bimg.NewImage(buffer).Process(option.VipsOptions)
				if err != nil {
					fmt.Printf("error creating new file %s ::: %s\n", file.Name(), err.Error())
				}
				err = os.MkdirAll("/mount/output", os.ModePerm)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fullName := cleanName + option.FileNameSuffix
				err = ioutil.WriteFile("/mount/output/"+fullName, newImg, os.ModePerm)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("generating " + fullName)
			}

		}

	}
}
