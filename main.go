package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

func main() {
	inputImageFile := "input.jpg"
	outputImageFile := "output.jpg"
	secretMessage := "This is a secret message."

	img, err := readJPEG(inputImageFile)
	if err != nil {
		fmt.Println("Error reading input image:", err)
		return
	}

	encodedImg := encodeMessage(img, secretMessage)
	err = writeJPEG(outputImageFile, encodedImg)
	if err != nil {
		fmt.Println("Error writing output image:", err)
		return
	}

	fmt.Println("Secret message encoded in", outputImageFile)
}

func readJPEG(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func writeJPEG(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
}

func encodeMessage(img image.Image, message string) image.Image {
	bounds := img.Bounds()
	rect := image.Rect(0, 0, bounds.Dx(), bounds.Dy())
	encodedImg := image.NewRGBA(rect)

	// Convert the message to binary
	binaryMessage := ""
	for _, char := range message {
		binaryMessage += fmt.Sprintf("%08b", char)
	}

	messageIndex := 0
	messageLen := len(binaryMessage)

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			if messageIndex < messageLen {
				c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
				c.R = setLSB(c.R, binaryMessage[messageIndex])
				messageIndex++

				if messageIndex < messageLen {
					c.G = setLSB(c.G, binaryMessage[messageIndex])
					messageIndex++
				}

				if messageIndex < messageLen {
					c.B = setLSB(c.B, binaryMessage[messageIndex])
					messageIndex++
				}

				encodedImg.SetRGBA(x, y, c)
			} else {
				encodedImg.Set(x, y, img.At(x, y))
			}
		}
	}

	return encodedImg
}

func setLSB(value uint8, bit byte) uint8 {
	if bit == '1' {
		return value | 1
	} else {
		return value &^ 1
	}
}
