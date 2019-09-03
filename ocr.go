// Package main provides ...
package main

// #cgo LDFLAGS: -L /usr/local/lib -ltesseract -llept
// #include "leptonica/allheaders.h"
// #include "tesseract/capi.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"os"
	"unsafe"
)

type Tess struct {
	tba *C.TessBaseAPI
}

// NewTess creates and returns a new tesseract instance.
func NewTess() (*Tess, error) {
	datapath := os.Getenv("TESSDATA_PREFIX")
	// create new empty TessBaseAPI
	tba := C.TessBaseAPICreate()

	// prepare string for C call
	cDatapath := C.CString(datapath)
	defer C.free(unsafe.Pointer(cDatapath))

	// prepare string for C call
	cLanguage := C.CString("eng")
	defer C.free(unsafe.Pointer(cLanguage))

	// initialize datapath and language on TessBaseAPI
	res := C.TessBaseAPIInit3(tba, cDatapath, cLanguage)
	if res != 0 {
		return nil, errors.New("could not initiate new Tess instance")
	}

	tess := &Tess{
		tba: tba,
	}

	return tess, nil
}

func (t *Tess) Image2text(data []byte) (string, error) {
	cpix := C.pixReadMem((*C.uchar)(unsafe.Pointer(&data[0])), C.size_t(len(data)))
	C.TessBaseAPISetImage2(t.tba, cpix)
	res := t.Text()
	return res, nil
}

// Text returns text after analysing the image(s)
func (t *Tess) Text() string {
	cText := C.TessBaseAPIGetUTF8Text(t.tba)
	defer C.free(unsafe.Pointer(cText))
	text := C.GoString(cText)
	return text
}
