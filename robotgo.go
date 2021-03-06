// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// http://www.
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

/*
//#if defined(IS_MACOSX)
	#cgo darwin CFLAGS: -x objective-c  -Wno-deprecated-declarations -I/usr/local/opt/libpng/include -I/usr/local/opt/zlib/include
	#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit -framework Carbon -framework CoreFoundation -L/usr/local/opt/libpng/lib -lpng -L/usr/local/opt/zlib/lib -lz
//#elif defined(USE_X11)
	#cgo linux CFLAGS:-I/usr/src
	#cgo linux LDFLAGS:-L/usr/src -lpng -lz -lX11 -lXtst -lX11-xcb -lxcb -lxcb-xkb -lxkbcommon -lxkbcommon-x11 -lm
//#endif
	#cgo windows LDFLAGS: -lgdi32 -luser32 -lpng -lz
//#include <AppKit/NSEvent.h>
#include "screen/goScreen.h"
#include "mouse/goMouse.h"
#include "key/goKey.h"
#include "bitmap/goBitmap.h"
#include "event/goEvent.h"
#include "window/goWindow.h"
*/
import "C"

import (
	. "fmt"
	"reflect"
	"unsafe"
	// "runtime"
	// "syscall"
)

/*
      _______.  ______ .______       _______  _______ .__   __.
    /       | /      ||   _  \     |   ____||   ____||  \ |  |
   |   (----`|  ,----'|  |_)  |    |  |__   |  |__   |   \|  |
    \   \    |  |     |      /     |   __|  |   __|  |  . `  |
.----)   |   |  `----.|  |\  \----.|  |____ |  |____ |  |\   |
|_______/     \______|| _| `._____||_______||_______||__| \__|
*/

type Bit_map struct {
	ImageBuffer   *uint8
	Width         int
	Height        int
	Bytewidth     int
	BitsPerPixel  uint8
	BytesPerPixel uint8
}

func GetPixelColor(x, y int) string {
	cx := C.size_t(x)
	cy := C.size_t(y)
	color := C.aGetPixelColor(cx, cy)
	// color := C.aGetPixelColor(x, y)
	gcolor := C.GoString(color)
	defer C.free(unsafe.Pointer(color))
	return gcolor
}

func GetScreenSize() (int, int) {
	size := C.aGetScreenSize()
	// Println("...", size, size.width)
	return int(size.width), int(size.height)
}

func GetXDisplayName() string {
	name := C.aGetXDisplayName()
	gname := C.GoString(name)
	defer C.free(unsafe.Pointer(name))
	return gname
}

func SetXDisplayName(name string) string {
	cname := C.CString(name)
	str := C.aSetXDisplayName(cname)
	gstr := C.GoString(str)
	return gstr
}

func CaptureScreen(args ...int) C.MMBitmapRef {
	var x C.size_t
	var y C.size_t
	var w C.size_t
	var h C.size_t
	Try(func() {
		x = C.size_t(args[0])
		y = C.size_t(args[1])
		w = C.size_t(args[2])
		h = C.size_t(args[3])
	}, func(e interface{}) {
		// Println("err:::", e)
		x = 0
		y = 0
		//Get screen size.
		var displaySize C.MMSize
		displaySize = C.getMainDisplaySize()
		w = displaySize.width
		h = displaySize.height
	})

	bit := C.aCaptureScreen(x, y, w, h)
	// Println("...", bit.width)
	return bit
}

func Capture_Screen(args ...int) Bit_map {
	var x C.size_t
	var y C.size_t
	var w C.size_t
	var h C.size_t
	Try(func() {
		x = C.size_t(args[0])
		y = C.size_t(args[1])
		w = C.size_t(args[2])
		h = C.size_t(args[3])
	}, func(e interface{}) {
		// Println("err:::", e)
		x = 0
		y = 0
		//Get screen size.
		var displaySize C.MMSize
		displaySize = C.getMainDisplaySize()
		w = displaySize.width
		h = displaySize.height
	})

	bit := C.aCaptureScreen(x, y, w, h)
	// Println("...", bit)
	bit_map := Bit_map{
		ImageBuffer:   (*uint8)(bit.imageBuffer),
		Width:         int(bit.width),
		Height:        int(bit.height),
		Bytewidth:     int(bit.bytewidth),
		BitsPerPixel:  uint8(bit.bitsPerPixel),
		BytesPerPixel: uint8(bit.bytesPerPixel),
	}

	return bit_map
}

/*
.___  ___.   ______    __    __       _______. _______
|   \/   |  /  __  \  |  |  |  |     /       ||   ____|
|  \  /  | |  |  |  | |  |  |  |    |   (----`|  |__
|  |\/|  | |  |  |  | |  |  |  |     \   \    |   __|
|  |  |  | |  `--'  | |  `--'  | .----)   |   |  |____
|__|  |__|  \______/   \______/  |_______/    |_______|

*/

type MPoint struct {
	x int
	y int
}

//C.size_t  int
func MoveMouse(x, y int) {
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.aMoveMouse(cx, cy)
}

func DragMouse(x, y int) {
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.aDragMouse(cx, cy)
}

func MoveMouseSmooth(x, y int, args ...float64) {
	cx := C.size_t(x)
	cy := C.size_t(y)

	var (
		low  C.double
		high C.double
	)

	Try(func() {
		low = C.double(args[2])
		high = C.double(args[3])
	}, func(e interface{}) {
		// Println("err:::", e)
		low = 5.0
		high = 500.0
	})

	C.aMoveMouseSmooth(cx, cy, low, high)
}

func GetMousePos() (int, int) {
	pos := C.aGetMousePos()
	// Println("pos:###", pos, pos.x, pos.y)
	x := int(pos.x)
	y := int(pos.y)
	// return pos.x, pos.y
	return x, y
}

func MouseClick(args ...interface{}) {
	var button C.MMMouseButton
	var double C.bool
	Try(func() {
		// button = args[0].(C.MMMouseButton)
		if args[0].(string) == "left" {
			button = C.LEFT_BUTTON
		}
		if args[0].(string) == "center" {
			button = C.CENTER_BUTTON
		}
		if args[0].(string) == "right" {
			button = C.RIGHT_BUTTON
		}
		double = C.bool(args[1].(bool))
	}, func(e interface{}) {
		// Println("err:::", e)
		button = C.LEFT_BUTTON
		double = false
	})
	C.aMouseClick(button, double)
}

func MouseToggle(args ...interface{}) {
	var button C.MMMouseButton
	Try(func() {
		// button = args[1].(C.MMMouseButton)
		if args[1].(string) == "left" {
			button = C.LEFT_BUTTON
		}
		if args[1].(string) == "center" {
			button = C.CENTER_BUTTON
		}
		if args[1].(string) == "right" {
			button = C.RIGHT_BUTTON
		}
	}, func(e interface{}) {
		// Println("err:::", e)
		button = C.LEFT_BUTTON
	})
	down := C.CString(args[0].(string))
	C.aMouseToggle(down, button)
}

func SetMouseDelay(x int) {
	cx := C.size_t(x)
	C.aSetMouseDelay(cx)
}

func ScrollMouse(x int, y string) {
	cx := C.size_t(x)
	z := C.CString(y)
	C.aScrollMouse(cx, z)
	defer C.free(unsafe.Pointer(z))
}

/*
 __  ___  ___________    ____ .______     ______        ___      .______       _______
|  |/  / |   ____\   \  /   / |   _  \   /  __  \      /   \     |   _  \     |       \
|  '  /  |  |__   \   \/   /  |  |_)  | |  |  |  |    /  ^  \    |  |_)  |    |  .--.  |
|    <   |   __|   \_    _/   |   _  <  |  |  |  |   /  /_\  \   |      /     |  |  |  |
|  .  \  |  |____    |  |     |  |_)  | |  `--'  |  /  _____  \  |  |\  \----.|  '--'  |
|__|\__\ |_______|   |__|     |______/   \______/  /__/     \__\ | _| `._____||_______/

*/
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

func KeyTap(args ...interface{}) {
	var akey string
	var akeyt string
	// var ckeyarr []*C.char
	ckeyarr := make([](*_Ctype_char), 0)
	var keyarr []string
	var num int

	Try(func() {
		if reflect.TypeOf(args[1]) == reflect.TypeOf(keyarr) {

			keyarr = args[1].([]string)

			num = len(keyarr)

			for i, _ := range keyarr {
				ckeyarr = append(ckeyarr, (*C.char)(unsafe.Pointer(C.CString(keyarr[i]))))
			}

		} else {
			akey = args[1].(string)

			Try(func() {
				akeyt = args[2].(string)

			}, func(e interface{}) {
				// Println("err:::", e)
				akeyt = "null"
			})
		}

	}, func(e interface{}) {
		// Println("err:::", e)
		akey = "null"
		keyarr = []string{"null"}
	})
	// }()

	zkey := C.CString(args[0].(string))

	if akey == "" && len(keyarr) != 0 {
		C.aKey_Tap(zkey, (**_Ctype_char)(unsafe.Pointer(&ckeyarr[0])), C.int(num))
	} else {
		// zkey := C.CString(args[0])
		amod := C.CString(akey)
		amodt := C.CString(akeyt)

		C.aKeyTap(zkey, amod, amodt)

		defer C.free(unsafe.Pointer(zkey))
		defer C.free(unsafe.Pointer(amod))
		defer C.free(unsafe.Pointer(amodt))
	}

}

func KeyToggle(args ...string) {
	var adown string
	var amkey string
	var amkeyt string
	Try(func() {
		adown = args[1]
		Try(func() {
			amkey = args[2]
			Try(func() {
				amkeyt = args[3]
			}, func(e interface{}) {
				// Println("err:::", e)
				amkeyt = "null"
			})
		}, func(e interface{}) {
			// Println("err:::", e)
			amkey = "null"
		})
	}, func(e interface{}) {
		// Println("err:::", e)
		adown = "null"
	})

	ckey := C.CString(args[0])
	cadown := C.CString(adown)
	camkey := C.CString(amkey)
	camkeyt := C.CString(amkeyt)
	// defer func() {
	str := C.aKeyToggle(ckey, cadown, camkey, camkeyt)
	Println(str)
	// }()
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(unsafe.Pointer(cadown))
	defer C.free(unsafe.Pointer(camkey))
	defer C.free(unsafe.Pointer(camkeyt))
}

func TypeString(x string) {
	cx := C.CString(x)
	C.aTypeString(cx)
	defer C.free(unsafe.Pointer(cx))
}

func TypeStringDelayed(x string, y int) {
	cx := C.CString(x)
	cy := C.size_t(y)
	C.aTypeStringDelayed(cx, cy)
	defer C.free(unsafe.Pointer(cx))
}

func SetKeyboardDelay(x int) {
	C.aSetKeyboardDelay(C.size_t(x))
}

/*
.______    __  .___________..___  ___.      ___      .______
|   _  \  |  | |           ||   \/   |     /   \     |   _  \
|  |_)  | |  | `---|  |----`|  \  /  |    /  ^  \    |  |_)  |
|   _  <  |  |     |  |     |  |\/|  |   /  /_\  \   |   ___/
|  |_)  | |  |     |  |     |  |  |  |  /  _____  \  |  |
|______/  |__|     |__|     |__|  |__| /__/     \__\ | _|
*/
func FindBitmap(args ...interface{}) (int, int) {
	var bit C.MMBitmapRef
	bit = args[0].(C.MMBitmapRef)

	var rect C.MMRect
	Try(func() {
		rect.origin.x = C.size_t(args[1].(int))
		rect.origin.y = C.size_t(args[2].(int))
		rect.size.width = C.size_t(args[3].(int))
		rect.size.height = C.size_t(args[4].(int))
	}, func(e interface{}) {
		// Println("err:::", e)
		// rect.origin.x = x
		// rect.origin.y = y
		// rect.size.width = w
		// rect.size.height = h
	})

	pos := C.aFindBitmap(bit, rect)
	// Println("pos----", pos)
	return int(pos.x), int(pos.y)
}

func OpenBitmap(args ...interface{}) C.MMBitmapRef {
	path := C.CString(args[0].(string))
	var mtype C.uint16_t
	Try(func() {
		mtype = C.uint16_t(args[1].(int))
	}, func(e interface{}) {
		// Println("err:::", e)
		mtype = 1
	})
	bit := C.aOpenBitmap(path, mtype)
	// Println("opening...", bit)
	return bit
	// defer C.free(unsafe.Pointer(path))
}

func SaveBitmap(args ...interface{}) {
	var mtype C.uint16_t
	Try(func() {
		mtype = C.uint16_t(args[2].(int))
	}, func(e interface{}) {
		// Println("err:::", e)
		mtype = 1
	})

	path := C.CString(args[1].(string))
	savebit := C.aSaveBitmap(args[0].(C.MMBitmapRef), path, mtype)
	Println("saved...", savebit)
	// return bit
	// defer C.free(unsafe.Pointer(path))
}

// func SaveBitmap(bit C.MMBitmapRef, gpath string, mtype C.MMImageType) {
// 	path := C.CString(gpath)
// 	savebit := C.aSaveBitmap(bit, path, mtype)
// 	Println("saveing...", savebit)
// 	// return bit
// 	// defer C.free(unsafe.Pointer(path))
// }

func TostringBitmap(bit C.MMBitmapRef) *C.char {
	str_bit := C.aTostringBitmap(bit)
	// Println("...", str_bit)
	return str_bit
}

func GetPortion(bit C.MMBitmapRef, x, y, w, h C.size_t) C.MMBitmapRef {
	var rect C.MMRect
	rect.origin.x = x
	rect.origin.y = y
	rect.size.width = w
	rect.size.height = h

	pos := C.aGetPortion(bit, rect)
	return pos
}

func Convert(args ...interface{}) {
	var mtype int
	Try(func() {
		mtype = args[2].(int)
	}, func(e interface{}) {
		Println("err:::", e)
		mtype = 1
	})
	//C.CString()
	opath := args[0].(string)
	spath := args[1].(string)
	bit_map := OpenBitmap(opath)
	// Println("a----", bit_map)
	SaveBitmap(bit_map, spath, mtype)
}

/*
------------ ---    ---  ------------ ----    ---- ------------
************ ***    ***  ************ *****   **** ************
----         ---    ---  ----         ------  ---- ------------
************ ***    ***  ************ ************     ****
------------ ---    ---  ------------ ------------     ----
****          ********   ****         ****  ******     ****
------------   ------    ------------ ----   -----     ----
************    ****     ************ ****    ****     ****

*/
func AddEvent(aeve string) int {
	cs := C.CString(aeve)
	eve := C.aEvent(cs)
	// Println("event@@", eve)
	geve := int(eve)
	return geve
}

func LEvent(aeve string) int {
	cs := C.CString(aeve)
	eve := C.aEvent(cs)
	// Println("event@@", eve)
	geve := int(eve)
	return geve
}

/*
____    __    ____  __  .__   __.  _______   ______   ____    __    ____
\   \  /  \  /   / |  | |  \ |  | |       \ /  __  \  \   \  /  \  /   /
 \   \/    \/   /  |  | |   \|  | |  .--.  |  |  |  |  \   \/    \/   /
  \            /   |  | |  . `  | |  |  |  |  |  |  |   \            /
   \    /\    /    |  | |  |\   | |  '--'  |  `--'  |    \    /\    /
    \__/  \__/     |__| |__| \__| |_______/ \______/      \__/  \__/

*/
func ShowAlert(title, msg string, args ...string) int {
	var (
		// title         string
		// msg           string
		defaultButton string
		cancelButton  string
	)

	Try(func() {
		// title = args[0]
		// msg = args[1]
		defaultButton = args[0]
		cancelButton = args[1]
	}, func(e interface{}) {
		defaultButton = "Ok"
		cancelButton = "Cancel"
	})
	atitle := C.CString(title)
	amsg := C.CString(msg)
	adefaultButton := C.CString(defaultButton)
	acancelButton := C.CString(cancelButton)

	cbool := C.aShowAlert(atitle, amsg, adefaultButton, acancelButton)
	ibool := int(cbool)
	return ibool
}

func IsValid() bool {
	abool := C.aIsValid()
	gbool := bool(abool)
	// Println("bool---------", gbool)
	return gbool
}

func SetActive(win C.MData) {
	C.aSetActive(win)
}

func GetActive() C.MData {
	mdata := C.aGetActive()
	// Println("active----", mdata)
	return mdata
}

func CloseWindow() {
	C.aCloseWindow()
}

func GetHandle() int {
	hwnd := C.aGetHandle()
	ghwnd := int(hwnd)
	// Println("gethwnd---", ghwnd)
	return ghwnd
}

// func SetHandle(hwnd int) {
// 	chwnd := C.uintptr(hwnd)
// 	C.aSetHandle(chwnd)
// }

func GetTitle() string {
	title := C.aGetTitle()
	gtittle := C.GoString(title)
	// Println("title...", gtittle)
	return gtittle
}
