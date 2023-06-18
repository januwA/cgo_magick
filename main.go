package main

// `pkg-config --cflags --libs MagickWand`

/*
#cgo pkg-config: MagickWand

#include <stdlib.h>
#include <MagickWand/MagickWand.h>
*/
import "C"
import (
	"os"
	"strconv"
	"unsafe"
)

func main() {
	var magick_wand *C.MagickWand
	var bg_wand *C.PixelWand
	_in := C.CString(os.Args[1])
	_out := C.CString(os.Args[2])
	_color := C.CString("rgba(0,0,0,0)")

	// 初始化 MagickWand 环境
	C.MagickWandGenesis()

	magick_wand = C.NewMagickWand()
	bg_wand = C.NewPixelWand()

	defer func() {
		// 释放malloc的内存
		C.free(unsafe.Pointer(_in))
		C.free(unsafe.Pointer(_out))
		C.free(unsafe.Pointer(_color))

		// 清理 wand
		C.DestroyPixelWand(bg_wand)
		C.DestroyMagickWand(magick_wand)

		// 终止 MagickWand 环境
		C.MagickWandTerminus()
	}()

	if C.MagickReadImage(magick_wand, _in) == C.MagickFalse {
		return
	}

	if C.PixelSetColor(bg_wand, _color) == C.MagickFalse {
		return
	}

	degrees, _ := strconv.ParseFloat(os.Args[3], 64)

	C.MagickResetIterator(magick_wand)
	for C.MagickNextImage(magick_wand) != C.MagickFalse {
		if C.MagickRotateImage(magick_wand, bg_wand, C.double(degrees)) == C.MagickFalse {
			return
		}
	}

	if C.MagickWriteImages(magick_wand, _out, 1) == C.MagickFalse {
		return
	}

}
