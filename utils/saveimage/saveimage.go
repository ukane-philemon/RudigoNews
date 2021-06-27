package saveimage

import (
	"image"
	"io"
	"log"
	"mime/multipart"
	"os"
)


func SaveImage(FeaturedImage multipart.File, header *multipart.FileHeader) (im image.Config, err error) {
defer FeaturedImage.Close()
			filetype := header.Header.Get("Content-Type")

			switch filetype {
			case "image/jpeg", "image/jpg", "image/gif", "image/png", "mage/bmp", "image/webp", "image/tiff":
				// Creating file in folder
				out, err := os.Create("upload/" + header.Filename)
				if err != nil {
					log.Println( "Unable to create the file for writing. Check your write access privilege")
					return im, err
				}
				defer out.Close()
				// write the content from POST REQUEST to the file
				_, err = io.Copy(out, FeaturedImage)
				if err != nil {
					log.Println(err)

				}
			default: 
			return 

			}
			//open file to read image width and height
			imagefile, err := os.Open("upload/" + header.Filename)
			if err != nil {
				log.Println(err)
				return
			}
			defer imagefile.Close()
			im, _, err = image.DecodeConfig(imagefile)

			if err != nil {
				log.Println(err)
				return
			}
 return im, err
}