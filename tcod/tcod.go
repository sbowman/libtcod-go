/*
libtcod-go is a go library providing bindings for libtcod

Most of the API is wrapped except:
- custom containers - Go has it's own containers
- threads, mutexes and semaphores - they are replaced by goroutines and channels
- SDL renderer - currently Go has very cumbersome C callback mechanism
*/
package tcod

/*
 #cgo LDFLAGS:-ltcod
 #include <stdio.h>
 #include <stdlib.h>
 #include <string.h>
 #include "include/libtcod.h"

 // This is a workaround for cgo disability to process varargs
 // These functions are copied verbatim from console_c and replaced ... with simple string
 // Formatting will be done in Go functions

 void _TCOD_console_print(TCOD_console_t con,int x, int y, char *s) {
 	TCOD_console_print(con,x,y,"%s",s);
 }

 void  _TCOD_console_print_ex(TCOD_console_t con,int x, int y, TCOD_bkgnd_flag_t flag, TCOD_alignment_t alignment, const char *s) {
 	TCOD_console_print_ex(con,x,y,flag,alignment,"%s",s);
 }

 int _TCOD_console_print_rect(TCOD_console_t con,int x, int y, int w, int h, char *s) {
 	return TCOD_console_print_rect(con,x,y,w,h,"%s",s);
 }


 int _TCOD_console_print_rect_ex(TCOD_console_t con,int x, int y, int w, int h, TCOD_bkgnd_flag_t flag, TCOD_alignment_t alignment, const char *s) {
	return TCOD_console_print_rect_ex(con, x, y, w, h, flag, alignment, "%s", s);
 }

 int _TCOD_console_height_rect(TCOD_console_t con,int x, int y, int w, int h, char *s) {
 	return TCOD_console_get_height_rect(con,x,y,w,h, "%s", s);
 }

 void _TCOD_console_print_frame(TCOD_console_t con,int x,int y,int w,int h, bool empty, TCOD_bkgnd_flag_t flag, char *s) {
 	TCOD_console_put_char(con,x,y,TCOD_CHAR_NW,flag);
 	TCOD_console_put_char(con,x+w-1,y,TCOD_CHAR_NE,flag);
 	TCOD_console_put_char(con,x,y+h-1,TCOD_CHAR_SW,flag);
 	TCOD_console_put_char(con,x+w-1,y+h-1,TCOD_CHAR_SE,flag);
 	TCOD_console_hline(con,x+1,y,w-2,flag);
 	TCOD_console_hline(con,x+1,y+h-1,w-2,flag);

 	if ( h > 2 ) {
 		TCOD_console_vline(con,x,y+1,h-2,flag);
 		TCOD_console_vline(con,x+w-1,y+1,h-2,flag);
 		if ( empty ) {
 			TCOD_console_rect(con,x+1,y+1,w-2,h-2,true,flag);
 		}
 	}
 	if (s) {
 		int xs;
 		TCOD_color_t tmp;
 		xs = x + (w-strlen(s)-2)/2;

		tmp = TCOD_console_get_default_background(con);
		TCOD_console_set_default_background(con, TCOD_console_get_default_foreground(con));
		TCOD_console_set_default_foreground(con, tmp);

 		TCOD_console_print(con,xs,y," %s ",s);

		tmp = TCOD_console_get_default_background(con);
		TCOD_console_set_default_background(con, TCOD_console_get_default_foreground(con));
		TCOD_console_set_default_foreground(con, tmp);
 	}
   }


  float _TCOD_heightmap_get_nth_value(const TCOD_heightmap_t *hm, int nth) {
  	return hm->values[nth];
  }

  void _TCOD_heightmap_set_nth_value(const TCOD_heightmap_t *hm, int nth, float val) {
  	hm->values[nth] = val;
  }

  typedef struct {
  	char *name;
  	TCOD_value_type_t value_type;
  	TCOD_value_t value;
  } _prop_t;

  bool _TCOD_sys_file_exists(const char * filename) {
	return TCOD_sys_file_exists(filename);
  }

*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

type void unsafe.Pointer

//
// Misc functions
//

type BkgndFlag C.TCOD_bkgnd_flag_t

func BkgndAlpha(alpha float32) BkgndFlag {
	return BkgndFlag(BkgndAlph | (((uint8)(alpha * 255)) << 8))
}

func BkgndAddAlpha(alpha float32) BkgndFlag {
	return BkgndFlag(BkgndAdda | (((uint8)(alpha * 255)) << 8))
}

func If(condition bool, tv, fv interface{}) interface{} {
	if condition {
		return tv
	} else {
		return fv
	}
}

func Clamp(a, b, x int) int {
	return If(x < a, a, If(x > b, b, x).(int)).(int)
}

func ClampF(a, b, x float32) float32 {
	return If(x < a, a, If(x > b, b, x).(float32)).(float32)
}

//
// TODO should free those strings?
func toStringSlice(l C.TCOD_list_t, free bool) (result []string) {
	size := C.TCOD_list_size(l)

	result = make([]string, int(size))
	for i := 0; i < int(size); i++ {
		c := (*C.char)(C.TCOD_list_get(l, C.int(i)))
		result[i] = C.GoString(c)
		if free {
			C.free(unsafe.Pointer(c))
		}
	}
	if free {
		C.TCOD_list_delete(l)
	}
	return
}


//
//
// Event
//
//

// Event is a system event, such as a keypress, mouse movement, or finger tap.
type Event int

func SysCheckForEvent(eventMask int, key *Key, mouse *Mouse) Event {
	var cKey C.TCOD_key_t
	var cMouse C.TCOD_mouse_t

	event := C.TCOD_sys_check_for_event(C.int(eventMask), &cKey, &cMouse)

	if key != nil {
		toKeyPtr(cKey, key)
	}

	if mouse != nil {
		toMousePtr(cMouse, mouse)
	}

	return Event(int(event))
}

//
//
// Key handling
//
//

type KeyCode C.TCOD_keycode_t

type Key struct {
	VK      KeyCode
	C       byte
	Pressed bool
	LAlt    bool
	LCtrl   bool
	RAlt    bool
	RCtrl   bool
	Shift   bool
}

func toKey(k C.TCOD_key_t) (result Key) {
	result.VK = KeyCode(k.vk)
	result.C = byte(k.c)
	result.Pressed = toBool(k.pressed)
	result.LAlt = toBool(k.lalt)
	result.LCtrl = toBool(k.lctrl)
	result.RAlt = toBool(k.ralt)
	result.RCtrl = toBool(k.rctrl)
	result.Shift = toBool(k.shift)
	return
}

func toKeyPtr(k C.TCOD_key_t, result *Key) {
	result.VK = KeyCode(k.vk)
	result.C = byte(k.c)
	result.Pressed = toBool(k.pressed)
	result.LAlt = toBool(k.lalt)
	result.LCtrl = toBool(k.lctrl)
	result.RAlt = toBool(k.ralt)
	result.RCtrl = toBool(k.rctrl)
	result.Shift = toBool(k.shift)
}

func fromKey(k Key) (result C.TCOD_key_t) {
	result.vk = C.TCOD_keycode_t(k.VK)
	result.c = C.char(k.C)
	result.pressed = fromBool(k.Pressed)
	result.lalt = fromBool(k.LAlt)
	result.lctrl = fromBool(k.LCtrl)
	result.ralt = fromBool(k.RAlt)
	result.rctrl = fromBool(k.RCtrl)
	result.shift = fromBool(k.Shift)
	return
}

//
//
// Bool handling
//
func toBool(b C.bool) bool {
	return bool(b)
}

func fromBool(b bool) C.bool {
	return C.bool(b)
}

//
// Color handling
//
type Color struct {
	R uint8
	G uint8
	B uint8
}

type ColCtrl C.TCOD_colctrl_t

func fromColor(c Color) (result C.TCOD_color_t) {
	result.r = C.uint8_t(c.R)
	result.g = C.uint8_t(c.G)
	result.b = C.uint8_t(c.B)
	return
}

func toColor(c C.TCOD_color_t) (result Color) {
	result.R = uint8(c.r)
	result.G = uint8(c.g)
	result.B = uint8(c.b)
	return
}

func NewColorRGB(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b}
}

func NewColorHSV(h, s, v float32) Color {
	return toColor(C.TCOD_color_HSV(C.float(h), C.float(s), C.float(v)))
}

// basic operations
func (color Color) Equals(c2 Color) bool {
	cc1 := fromColor(color)
	cc2 := fromColor(c2)
	return toBool(C.TCOD_color_equals(cc1, cc2))
}

func (color Color) Add(c2 Color) Color {
	cc1 := fromColor(color)
	cc2 := fromColor(c2)
	return toColor(C.TCOD_color_add(cc1, cc2))
}

func (color Color) Subtract(c2 Color) Color {
	cc1 := fromColor(color)
	cc2 := fromColor(c2)
	return toColor(C.TCOD_color_subtract(cc1, cc2))
}

func (color Color) Multiply(c2 Color) Color {
	cc1 := fromColor(color)
	cc2 := fromColor(c2)
	return toColor(C.TCOD_color_multiply(cc1, cc2))
}

func (color Color) MultiplyScalar(value float32) Color {
	c := fromColor(color)
	return toColor(C.TCOD_color_multiply_scalar(c, C.float(value)))
}

func (color Color) Lerp(c2 Color, coef float32) Color {
	cc1 := fromColor(color)
	cc2 := fromColor(c2)
	return toColor(C.TCOD_color_lerp(cc1, cc2, C.float(coef)))
}

// HSV transformations

func (color Color) Lighten(ratio float32) Color {
	return color.Lerp(White, ratio)
}

func (color Color) Darken(ratio float32) Color {
	return color.Lerp(Black, ratio)
}

func (color Color) SetHSV(h float32, s float32, v float32) Color {
	c := C.TCOD_color_t{}
	C.TCOD_color_set_HSV(&c, C.float(h), C.float(s), C.float(v))
	return toColor(c)
}

func (color Color) GetHue() float32 {
	return float32(C.TCOD_color_get_hue(fromColor(color)))
}

func (color Color) SetHue(h float32) Color {
	c := C.TCOD_color_t{}
	C.TCOD_color_set_hue(&c, C.float(h))
	return toColor(c)
}

func (color Color) GetSaturation() float32 {
	return float32(C.TCOD_color_get_saturation(fromColor(color)))
}

func (color Color) SetSaturation(h float32) Color {
	c := C.TCOD_color_t{}
	C.TCOD_color_set_saturation(&c, C.float(h))
	return toColor(c)
}

func (color Color) GetValue() float32 {
	return float32(C.TCOD_color_get_value(fromColor(color)))
}

func (color Color) SetValue(h float32) Color {
	c := C.TCOD_color_t{}
	C.TCOD_color_set_value(&c, C.float(h))
	return toColor(c)
}

func (color Color) ShiftHue(hshift float32) Color {
	c := C.TCOD_color_t{}
	C.TCOD_color_shift_hue(&c, C.float(hshift))
	return toColor(c)
}

func (color Color) ScaleHSV(scoef, vcoef float32) Color {
	c := C.TCOD_color_t{}
	C.TCOD_color_scale_HSV(&c, C.float(scoef), C.float(vcoef))
	return toColor(c)
}

func (color Color) GetHSV() (h, s, v float32) {
	var ch, cs, sv C.float
	C.TCOD_color_get_HSV(fromColor(color), &ch, &cs, &sv)
	h = float32(ch)
	s = float32(cs)
	v = float32(sv)
	return
}

func ColorGenMap(cmap []Color, nbKey int, keyColor []Color, keyIndex []int) {
	for segment := 0; segment < nbKey-1; segment++ {
		idxStart := keyIndex[segment]
		idxEnd := keyIndex[segment+1]
		for idx := idxStart; idx <= idxEnd; idx++ {
			cmap[idx] = keyColor[segment].Lerp(keyColor[segment+1], float32(idx-idxStart)/float32(idxEnd-idxStart))
		}
	}
}

//
//
// Mouse
//
//
type Mouse struct {
	X, Y           int
	Dx, Dy         int
	Cx, Cy         int
	Dcx, Dcy       int
	LButton        bool
	RButton        bool
	MButton        bool
	LButtonPressed bool
	RButtonPressed bool
	MButtonPressed bool
	WheelUp        bool
	WheelDown      bool
}

func fromMouse(m Mouse) (result C.TCOD_mouse_t) {
	result.x = C.int(m.X)
	result.y = C.int(m.Y)
	result.dx = C.int(m.Dx)
	result.dy = C.int(m.Dy)
	result.cx = C.int(m.Cx)
	result.cy = C.int(m.Cy)
	result.dcx = C.int(m.Dcx)
	result.dcy = C.int(m.Dcy)
	result.lbutton = fromBool(m.LButton)
	result.rbutton = fromBool(m.RButton)
	result.mbutton = fromBool(m.MButton)
	result.lbutton_pressed = fromBool(m.LButtonPressed)
	result.rbutton_pressed = fromBool(m.RButtonPressed)
	result.mbutton_pressed = fromBool(m.MButtonPressed)
	result.wheel_up = fromBool(m.WheelUp)
	result.wheel_down = fromBool(m.WheelDown)
	return
}

func toMouse(m C.TCOD_mouse_t) (result Mouse) {
	result.X = int(m.x)
	result.Y = int(m.y)
	result.Dx = int(m.dx)
	result.Dy = int(m.dy)
	result.Cx = int(m.cx)
	result.Cy = int(m.cy)
	result.Dcx = int(m.dcx)
	result.Dcy = int(m.dcy)
	result.LButton = toBool(m.lbutton)
	result.RButton = toBool(m.rbutton)
	result.MButton = toBool(m.mbutton)
	result.LButtonPressed = toBool(m.lbutton_pressed)
	result.RButtonPressed = toBool(m.rbutton_pressed)
	result.MButtonPressed = toBool(m.mbutton_pressed)
	result.WheelUp = toBool(m.wheel_up)
	result.WheelDown = toBool(m.wheel_down)
	return
}

func toMousePtr(m C.TCOD_mouse_t, result *Mouse) {
	result.X = int(m.x)
	result.Y = int(m.y)
	result.Dx = int(m.dx)
	result.Dy = int(m.dy)
	result.Cx = int(m.cx)
	result.Cy = int(m.cy)
	result.Dcx = int(m.dcx)
	result.Dcy = int(m.dcy)
	result.LButton = toBool(m.lbutton)
	result.RButton = toBool(m.rbutton)
	result.MButton = toBool(m.mbutton)
	result.LButtonPressed = toBool(m.lbutton_pressed)
	result.RButtonPressed = toBool(m.rbutton_pressed)
	result.MButtonPressed = toBool(m.mbutton_pressed)
	result.WheelUp = toBool(m.wheel_up)
	result.WheelDown = toBool(m.wheel_down)
}


func MouseGetStatus() Mouse {
	return toMouse(C.TCOD_mouse_get_status())
}

func MouseShowCursor(visible bool) {
	C.TCOD_mouse_show_cursor(fromBool(visible))
}

func MouseIsCursorVisible() bool {
	return toBool(C.TCOD_mouse_is_cursor_visible())
}

func MouseMove(x, y int) {
	C.TCOD_mouse_move(C.int(x), C.int(y))
}

//
//
// Console
//
//
type Alignment C.TCOD_alignment_t
type Renderer C.TCOD_renderer_t

type IConsole interface {
	GetData() C.TCOD_console_t
	GetDefaultBackground() Color
	GetDefaultForeground() Color
	SetDefaultForeground(color Color)
	SetDefaultBackground(color Color)
	Clear()
	GetCharBackground(x, y int) Color
	GetCharForeground(x, y int) Color
	SetCharBackground(x, y int, color Color, flag BkgndFlag)
	SetCharForeground(x, y int, color Color)
	SetChar(x, y int, c int)
	PutChar(x, y, c int, flag BkgndFlag)
	PutCharEx(x, y, c int, fore, back Color)
	Print(x, y int, fmts string, v ...interface{})
	PrintEx(x, y int, flag BkgndFlag, alignment Alignment, fmts string, v ...interface{})
	PrintRect(x, y, w, h int, fmts string, v ...interface{}) int
	PrintRectEx(x, y, w, h int, flag BkgndFlag, alignment Alignment, fmts string, v ...interface{}) int
	HeightRect(x, y, w, h int, fmts string, v ...interface{}) int
	SetBackgroundFlag(flag BkgndFlag)
	GetBackgroundFlag() BkgndFlag
	SetAlignment(alignment Alignment)
	GetAlignment() Alignment
	Rect(x, y, w, h int, clear bool, flag BkgndFlag)
	Hline(x, y, l int, flag BkgndFlag)
	Vline(x, y, l int, flag BkgndFlag)
	PrintFrame(x, y, w, h int, empty bool, flag BkgndFlag, fmts string, v ...interface{})
	GetChar(x, y int) int
	GetWidth() int
	GetHeight() int
	SetKeyColor(color Color)
	Blit(xSrc, ySrc, wSrc, hSrc int, dst IConsole, xDst, yDst int, foregroundAlpha, backgroundAlpha float32)
}

// Console

type Console struct {
	Data C.TCOD_console_t
}

func deleteConsole(c *Console) {
	C.TCOD_console_delete(c.Data)
}

func NewConsole(w, h int) *Console {
	result := &Console{C.TCOD_console_new(C.int(w), C.int(h))}
	runtime.SetFinalizer(result, deleteConsole)
	return result
}

func (console *Console) GetData() C.TCOD_console_t {
	return console.Data
}

func (console *Console) SetDefaultBackground(color Color) {
	C.TCOD_console_set_default_background(console.Data, fromColor(color))
}

func (console *Console) SetDefaultForeground(color Color) {
	C.TCOD_console_set_default_foreground(console.Data, fromColor(color))
}

func (console *Console) Clear() {
	C.TCOD_console_clear(console.Data)
}

func (console *Console) SetCharBackground(x, y int, color Color, flag BkgndFlag) {
	ccolor := fromColor(color)
	C.TCOD_console_set_char_background(console.Data, C.int(x), C.int(y), ccolor, C.TCOD_bkgnd_flag_t(flag))
}

func (console *Console) SetCharForeground(x, y int, color Color) {
	ccolor := fromColor(color)
	C.TCOD_console_set_char_foreground(console.Data, C.int(x), C.int(y), ccolor)
}

func (console *Console) SetChar(x, y int, c int) {
	C.TCOD_console_set_char(console.Data, C.int(x), C.int(y), C.int(c))
}

func (console *Console) PutChar(x, y, c int, flag BkgndFlag) {
	C.TCOD_console_put_char(console.Data, C.int(x), C.int(y), C.int(c), C.TCOD_bkgnd_flag_t(flag))
}

func (console *Console) PutCharEx(x, y, c int, fore, back Color) {
	forec := fromColor(fore)
	backc := fromColor(back)
	C.TCOD_console_put_char_ex(console.Data, C.int(x), C.int(y), C.int(c), forec, backc)
}

func (console *Console) Print(x, y int, fmts string, v ...interface{}) {
	s := fmt.Sprintf(fmts, v...)
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C._TCOD_console_print(console.Data, C.int(x), C.int(y), cs)
}

func (console *Console) PrintEx(x, y int, flag BkgndFlag, alignment Alignment, fmts string, v ...interface{}) {
	s := fmt.Sprintf(fmts, v...)
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C._TCOD_console_print_ex(console.Data, C.int(x), C.int(y), C.TCOD_bkgnd_flag_t(flag), C.TCOD_alignment_t(alignment), cs)
}

func (console *Console) PrintRect(x, y, w, h int, fmts string, v ...interface{}) int {
	s := fmt.Sprintf(fmts, v...)
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C._TCOD_console_print_rect(console.Data, C.int(x), C.int(y), C.int(w), C.int(h), cs))
}

func (console *Console) PrintRectEx(x, y, w, h int, flag BkgndFlag, alignment Alignment, fmts string, v ...interface{}) int {
	s := fmt.Sprintf(fmts, v...)
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C._TCOD_console_print_rect_ex(console.Data, C.int(x), C.int(y), C.int(w), C.int(h), C.TCOD_bkgnd_flag_t(flag),
		C.TCOD_alignment_t(alignment), cs))
}

func (console *Console) SetBackgroundFlag(flag BkgndFlag) {
	C.TCOD_console_set_background_flag(console.Data, C.TCOD_bkgnd_flag_t(flag))
}

func (console *Console) GetBackgroundFlag() BkgndFlag {
	return BkgndFlag(C.TCOD_console_get_background_flag(console.Data))
}

func (console *Console) SetAlignment(alignment Alignment) {
	C.TCOD_console_set_alignment(console.Data, C.TCOD_alignment_t(alignment))
}

func (console *Console) GetAlignment() Alignment {
	return Alignment(C.TCOD_console_get_alignment(console.Data))
}

func (console *Console) HeightRect(x, y, w, h int, fmts string, v ...interface{}) int {
	s := fmt.Sprintf(fmts, v...)
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C._TCOD_console_height_rect(console.Data, C.int(x), C.int(y), C.int(w), C.int(h), cs))
}

func (console *Console) Rect(x, y, w, h int, clear bool, flag BkgndFlag) {
	C.TCOD_console_rect(console.Data, C.int(x), C.int(y), C.int(w), C.int(h), fromBool(clear), C.TCOD_bkgnd_flag_t(flag))
}

func (console *Console) Hline(x, y, l int, flag BkgndFlag) {
	C.TCOD_console_hline(console.Data, C.int(x), C.int(y), C.int(l), C.TCOD_bkgnd_flag_t(flag))
}

func (console *Console) Vline(x, y, l int, flag BkgndFlag) {
	C.TCOD_console_hline(console.Data, C.int(x), C.int(y), C.int(l), C.TCOD_bkgnd_flag_t(flag))
}

func (console *Console) PrintFrame(x, y, w, h int, empty bool, flag BkgndFlag, fmts string, v ...interface{}) {
	s := fmt.Sprintf(fmts, v...)
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C._TCOD_console_print_frame(console.Data, C.int(x), C.int(y), C.int(w), C.int(h),
		fromBool(empty), C.TCOD_bkgnd_flag_t(flag), cs)

}

// TODO check unicode support
// TCODLIB_API void TCOD_console_map_string_to_font_utf(const wchar_t *s, int fontCharX, int fontCharY);
// TCODLIB_API void TCOD_console_print_left_utf(TCOD_console_t con,int x, int y, TCOD_bkgnd_flag_t flag, const wchar_t *fmt, ...);
// TCODLIB_API void TCOD_console_print_right_utf(TCOD_console_t con,int x, int y, TCOD_bkgnd_flag_t flag, const wchar_t *fmt, ...);
// TCODLIB_API void TCOD_console_print_center_utf(TCOD_console_t con,int x, int y, TCOD_bkgnd_flag_t flag, const wchar_t *fmt, ...);
// TCODLIB_API int TCOD_console_print_left_rect_utf(TCOD_console_t con,int x, int y, int w, int h, TCOD_bkgnd_flag_t flag, const wchar_t *fmt, ...);
// TCODLIB_API int TCOD_console_print_right_rect_utf(TCOD_console_t con,int x, int y, int w, int h, TCOD_bkgnd_flag_t flag, const wchar_t *fmt, ...);
// TCODLIB_API int TCOD_console_print_center_rect_utf(TCOD_console_t con,int x, int y, int w, int h, TCOD_bkgnd_flag_t flag, const wchar_t *fmt, ...);
// TCODLIB_API int TCOD_console_height_left_rect_utf(TCOD_console_t con,int x, int y, int w, int h, const wchar_t *fmt, ...);
// TCODLIB_API int TCOD_console_height_right_rect_utf(TCOD_console_t con,int x, int y, int w, int h, const wchar_t *fmt, ...);
// TCODLIB_API int TCOD_console_height_center_rect_utf(TCOD_console_t con,int x, int y, int w, int h, const wchar_t *fmt, ...);
// #endif

func (console *Console) GetDefaultBackground() Color {
	return toColor(C.TCOD_console_get_default_background(console.Data))
}

func (console *Console) GetDefaultForeground() Color {
	return toColor(C.TCOD_console_get_default_foreground(console.Data))
}

func (console *Console) GetCharBackground(x, y int) Color {
	return toColor(C.TCOD_console_get_char_background(console.Data, C.int(x), C.int(y)))
}

func (console *Console) GetCharForeground(x, y int) Color {
	return toColor(C.TCOD_console_get_char_foreground(console.Data, C.int(x), C.int(y)))
}

func (console *Console) GetChar(x, y int) int {
	return int(C.TCOD_console_get_char(console.Data, C.int(x), C.int(y)))
}

func (console *Console) GetWidth() int {
	return int(C.TCOD_console_get_width(console.Data))
}

func (console *Console) GetHeight() int {
	return int(C.TCOD_console_get_height(console.Data))
}

func (console *Console) SetKeyColor(color Color) {
	ccolor := fromColor(color)
	C.TCOD_console_set_key_color(console.Data, ccolor)
}

func (console *Console) Blit(xSrc, ySrc, wSrc, hSrc int, dst IConsole, xDst, yDst int, foregroundAlpha, backgroundAlpha float32) {
	C.TCOD_console_blit(console.Data, C.int(xSrc), C.int(ySrc), C.int(wSrc), C.int(hSrc),
		dst.GetData(), C.int(xDst), C.int(yDst), C.float(foregroundAlpha), C.float(backgroundAlpha))
}

// RootConsole

type RootConsole struct {
	Console
}

func NewRootConsole(w, h int, title string, fullscreen bool, renderer Renderer) *RootConsole {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.TCOD_console_init_root(C.int(w), C.int(h), ctitle, fromBool(fullscreen), C.TCOD_renderer_t(renderer))
	// in root console, Data field is nil
	return &RootConsole{}
}

func NewRootConsoleWithFont(w, h int, title string, fullscreen bool, fontFile string, fontFlags, nbCharHoriz,
	nbCharVertic int, renderer Renderer) *RootConsole {
	cfontFile := C.CString(fontFile)
	defer C.free(unsafe.Pointer(cfontFile))
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.TCOD_console_set_custom_font(cfontFile, C.int(fontFlags), C.int(nbCharHoriz), C.int(nbCharVertic))
	C.TCOD_console_init_root(C.int(w), C.int(h), ctitle, fromBool(fullscreen), C.TCOD_renderer_t(renderer))
	// in root console, Data field is nil
	return &RootConsole{}
}

func (root *RootConsole) SetWindowTitle(title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.TCOD_console_set_window_title(ctitle)

}

func (root *RootConsole) SetFullscreen(fullscreen bool) {
	C.TCOD_console_set_fullscreen(fromBool(fullscreen))
}

func (root *RootConsole) IsFullscreen() bool {
	return toBool(C.TCOD_console_is_fullscreen())
}

func (root *RootConsole) IsWindowClosed() bool {
	return toBool(C.TCOD_console_is_window_closed())
}

func (root *RootConsole) SetCustomFont(fontFile string, flags int, nbCharHoriz int, nbCharVertic int) {
	cfontFile := C.CString(fontFile)
	defer C.free(unsafe.Pointer(cfontFile))
	C.TCOD_console_set_custom_font(cfontFile, C.int(flags), C.int(nbCharHoriz), C.int(nbCharVertic))
}

func (root *RootConsole) MapAsciiCodeToFont(asciiCode, fontCharX, fontCharY int) {
	C.TCOD_console_map_ascii_code_to_font(C.int(asciiCode), C.int(fontCharX), C.int(fontCharY))
}

func (root *RootConsole) MapAsciiCodesToFont(asciiCode, fontCharX, fontCharY int) {
	C.TCOD_console_map_ascii_code_to_font(C.int(asciiCode), C.int(fontCharX), C.int(fontCharY))
}

func (root *RootConsole) MapStringToFont(s string, fontCharX, fontCharY int) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.TCOD_console_map_string_to_font(cs, C.int(fontCharX), C.int(fontCharY))
}

func (root *RootConsole) SetDirty(x, y, w, h int) {
	C.TCOD_console_set_dirty(C.int(x), C.int(y), C.int(w), C.int(h))
}

func (root *RootConsole) SetFade(val uint8, fade Color) {
	ccolor := fromColor(fade)
	C.TCOD_console_set_fade(C.uint8_t(val), ccolor)
}

func (root *RootConsole) GetFade() uint8 {
	return uint8(C.TCOD_console_get_fade())
}

func (root *RootConsole) GetFadingColor() Color {
	return toColor(C.TCOD_console_get_fading_color())
}

func (root *RootConsole) Flush() {
	C.TCOD_console_flush()
}

func (root *RootConsole) SetColorControl(ctrl ColCtrl, fore, back Color) {
	forec := fromColor(fore)
	backc := fromColor(back)
	C.TCOD_console_set_color_control(C.TCOD_colctrl_t(ctrl), forec, backc)
}

func (root *RootConsole) CheckForKeypress(flags int) Key {
	return toKey(C.TCOD_console_check_for_keypress(C.int(flags)))
}

func (root *RootConsole) WaitForKeypress(flush bool) Key {
	return toKey(C.TCOD_console_wait_for_keypress(fromBool(flush)))
}

func (root *RootConsole) SetKeyboardRepeat(initialDelay, interval int) {
	C.TCOD_console_set_keyboard_repeat(C.int(initialDelay), C.int(interval))
}

func (root *RootConsole) DisableKeyboardRepeat() {
	C.TCOD_console_disable_keyboard_repeat()
}

func (root *RootConsole) IsKeyPressed(keyCode KeyCode) bool {
	return toBool(C.TCOD_console_is_key_pressed(C.TCOD_keycode_t(keyCode)))
}

func (root *RootConsole) Credits() {
	C.TCOD_console_credits()
}

func (root *RootConsole) ResetCredits() {
	C.TCOD_console_credits_reset()
}

func (root *RootConsole) RenderCredits(x, y int, alpha bool) bool {
	return toBool(C.TCOD_console_credits_render(C.int(x), C.int(y), fromBool(alpha)))
}

//
// Bresenham line algorithm
// Fully ported to Go for easier callbacks
//
//

type LineListener func(x, y int, userData interface{}) bool

type Point struct {
	x, y int
}

// thread-safe versions
type BresenhamData struct {
	stepx  int
	stepy  int
	e      int
	deltax int
	deltay int
	origx  int
	origy  int
	destx  int
	desty  int
}

var bresenhamData BresenhamData

func lineInitMt(xFrom, yFrom, xTo, yTo int, data *BresenhamData) {
	data.origx = xFrom
	data.origy = yFrom
	data.destx = xTo
	data.desty = yTo
	data.deltax = xTo - xFrom
	data.deltay = yTo - yFrom
	if data.deltax > 0 {
		data.stepx = 1
	} else if data.deltax < 0 {
		data.stepx = -1
	} else {
		data.stepx = 0
	}
	if data.deltay > 0 {
		data.stepy = 1
	} else if data.deltay < 0 {
		data.stepy = -1
	} else {
		data.stepy = 0
	}
	if data.stepx*data.deltax > data.stepy*data.deltay {
		data.e = data.stepx * data.deltax
		data.deltax *= 2
		data.deltay *= 2
	} else {
		data.e = data.stepy * data.deltay
		data.deltax *= 2
		data.deltay *= 2
	}
}

func lineStepMt(xCur, yCur *int, data *BresenhamData) bool {
	if data.stepx*data.deltax > data.stepy*data.deltay {
		if data.origx == data.destx {
			return true
		}
		data.origx += data.stepx
		data.e -= data.stepy * data.deltay
		if data.e < 0 {
			data.origy += data.stepy
			data.e += data.stepx * data.deltax
		}
	} else {
		if data.origy == data.desty {
			return true
		}
		data.origy += data.stepy
		data.e -= data.stepx * data.deltax
		if data.e < 0 {
			data.origx += data.stepx
			data.e += data.stepy * data.deltay
		}
	}
	*xCur = data.origx
	*yCur = data.origy
	return false
}

func lineInit(xFrom, yFrom, xTo, yTo int) {
	lineInitMt(xFrom, yFrom, xTo, yTo, &bresenhamData)
}

func lineStep(xCur, yCur *int) bool {
	return lineStepMt(xCur, yCur, &bresenhamData)
}

func LineMt(xo, yo, xd, yd int, listener LineListener, userData interface{}, data *BresenhamData) bool {
	lineInitMt(xo, yo, xd, yd, data)
	if !listener(xo, yo, userData) {
		return false
	}
	for !lineStepMt(&xo, &yo, data) {
		if !listener(xo, yo, userData) {
			return false
		}
	}
	return true
}

func Line(xo, yo, xd, yd int, userData interface{}, listener LineListener) bool {
	return LineMt(xo, yo, xd, yd, listener, userData, &bresenhamData)
}

// returns slice of Points where the line was drawn
func LinePoints(xo, yo, xd, yd int) []Point {
	result := []Point{}
	Line(xo, yo, xd, yd, nil, func(x, y int, data interface{}) bool {
		result = append(result, Point{x, y})
		return true
	})
	return result
}

//
//
// Name generator
//

//
func NamegenParse(filename string, random *Random) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	C.TCOD_namegen_parse(cfilename, random.Data)
}

// generate a name
func NamegenGenerate(name string) string {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	c := C.TCOD_namegen_generate(cname, fromBool(true))
	defer C.free(unsafe.Pointer(c))
	return C.GoString(c)
}

// generate a name using a custom generation rule
func NamegenGenerateCustom(name, rule string) string {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	crule := C.CString(rule)
	defer C.free(unsafe.Pointer(crule))
	c := C.TCOD_namegen_generate_custom(cname, crule, fromBool(true))
	defer C.free(unsafe.Pointer(c))
	return C.GoString(c)
}

// retrieve the list of all available syllable set names
func NamegenGetSets() []string {
	return toStringSlice(C.TCOD_namegen_get_sets(), false)
}

// delete a generator
func NamegenDestroy() {
	C.TCOD_namegen_destroy()
}

//
//
// Text field
// TODO this is available only in debug version? (1.5.0)
//
//
type Text struct {
	Data C.TCOD_text_t
}

func deleteText(t *Text) {
	C.TCOD_text_delete(t.Data)
}

func NewText(x, y, w, h, maxChars int) *Text {
	result := &Text{C.TCOD_text_init(C.int(x), C.int(y), C.int(w), C.int(h), C.int(maxChars))}
	runtime.SetFinalizer(result, deleteText)
	return result
}

func (txt *Text) SetProperties(cursorChar int, blinkInterval int, prompt string, tabSize int) {
	cprompt := C.CString(prompt)
	defer C.free(unsafe.Pointer(cprompt))
	C.TCOD_text_set_properties(txt.Data, C.int(cursorChar), C.int(blinkInterval), cprompt, C.int(tabSize))
}

func (txt *Text) SetColors(fore, back Color, backTransparency float32) {
	forec := fromColor(fore)
	backc := fromColor(back)
	C.TCOD_text_set_colors(txt.Data, forec, backc, C.float(backTransparency))
}

func (txt *Text) Update(key Key) {
	C.TCOD_text_update(txt.Data, fromKey(key))
}

func (txt *Text) Render(console IConsole) {
	C.TCOD_text_render(txt.Data, console.GetData())
}

func (console *Console) RenderText(text *Text) {
	C.TCOD_text_render(text.Data, console.Data)
}

func (txt *Text) Get() string {
	t := C.TCOD_text_get(txt.Data)
	return C.GoString(t)

}

func (txt *Text) Reset() {
	C.TCOD_text_reset(txt.Data)
}

func SysElapsedMilliseconds() uint32 {
	return uint32(C.TCOD_sys_elapsed_milli())
}

func SysElapsedSeconds() float32 {
	return float32(C.TCOD_sys_elapsed_seconds())
}

func SysSleepMilliseconds(val uint32) {
	C.TCOD_sys_sleep_milli(C.uint32_t(val))
}

func SysSaveScreenshot() {
	C.TCOD_sys_save_screenshot(nil)
}

func SysSaveScreenshotToFile(filename string) {
	if filename == "" {
		C.TCOD_sys_save_screenshot(nil)
	} else {
		cfilename := C.CString(filename)
		defer C.free(unsafe.Pointer(cfilename))
		C.TCOD_sys_save_screenshot(cfilename)
	}
}

func SysForceFullscreenResolution(width, height int) {
	C.TCOD_sys_force_fullscreen_resolution(C.int(width), C.int(height))
}

func SysSetFps(val int) {
	C.TCOD_sys_set_fps(C.int(val))
}

func SysGetFps() int {
	return int(C.TCOD_sys_get_fps())
}

func SysGetLastFrameLength() float32 {
	return float32(C.TCOD_sys_get_last_frame_length())
}

func SysGetCurrentResolution() (w, h int) {
	var cw, ch C.int
	C.TCOD_sys_get_current_resolution(&cw, &ch)
	w, h = int(cw), int(ch)
	return
}

func SysGetFullscreenOffsets() (offx, offy int) {
	var coffx, coffy C.int
	C.TCOD_sys_get_fullscreen_offsets(&coffx, &coffy)
	offx, offy = int(coffx), int(coffy)
	return
}

func SysUpdateChar(asciiCode, fontx, fonty int, img Image, x, y int) {
	C.TCOD_sys_update_char(C.int(asciiCode), C.int(fontx), C.int(fonty), img.Data, C.int(x), C.int(y))
}

func SysGetCharSize() (w, h int) {
	var cw, ch C.int
	C.TCOD_sys_get_char_size(&cw, &ch)
	w, h = int(cw), int(ch)
	return
}

// filesystem stuff
func SysCreateDirectory(path string) bool {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return toBool(C.TCOD_sys_create_directory(cpath))
}

func SysDeleteFile(path string) bool {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return toBool(C.TCOD_sys_delete_file(cpath))
}

func SysDeleteDirectory(path string) bool {
	return toBool(C.TCOD_sys_delete_directory(C.CString(path)))
}

func SysIsDirectory(path string) bool {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return toBool(C.TCOD_sys_is_directory(cpath))
}

func SysGetDirectoryContent(path, pattern string) []string {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	cpattern := C.CString(pattern)
	defer C.free(unsafe.Pointer(cpattern))
	return toStringSlice(
		C.TCOD_sys_get_directory_content(
			cpath, cpattern),
		true)
}

func SysFileExists(filename string) bool {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	return toBool(C._TCOD_sys_file_exists(cfilename))
}

func SysGetNumCores() int {
	return int(C.TCOD_sys_get_num_cores())
}

// Clipboard

func SysClipboardSet(value string) {
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	C.TCOD_sys_clipboard_set(cvalue)
}

func SysClipboardGet() string {
	return C.GoString(C.TCOD_sys_clipboard_get())
}

//
// Field Of View Map
//

type Map struct {
	Data C.TCOD_map_t
}

type FovAlgorithm C.TCOD_fov_algorithm_t

// destroy a map
func deleteMap(m *Map) {
	C.TCOD_map_delete(m.Data)
}

func NewMap(width, height int) *Map {
	result := &Map{C.TCOD_map_new(C.int(width), C.int(height))}
	runtime.SetFinalizer(result, deleteMap)
	return result
}

// set all cells as solid rock (cannot see through nor walk)
func (m *Map) Clear(isTransparent bool, isWalkable bool) {
	C.TCOD_map_clear(m.Data, fromBool(isTransparent), fromBool(isWalkable))
}

// copy a map to another, reallocating it when needed
func (m *Map) Copy(dest Map) {
	C.TCOD_map_copy(m.Data, dest.Data)
}

// change a cell properties
func (m *Map) SetProperties(x, y int, isTransparent bool, isWalkable bool) {
	C.TCOD_map_set_properties(m.Data, C.int(x), C.int(y), fromBool(isTransparent), fromBool(isWalkable))
}

// calculate the field of view (potentially visible cells from player_x,player_y)
func (m *Map) ComputeFov(playerX, playerY, maxRadius int, lightWalls bool, algo FovAlgorithm) {
	C.TCOD_map_compute_fov(m.Data, C.int(playerX), C.int(playerY),
		C.int(maxRadius), fromBool(lightWalls),
		C.TCOD_fov_algorithm_t(algo))
}

// check if a cell is in the last computed field of view
func (m *Map) IsInFov(x, y int) bool {
	return toBool(C.TCOD_map_is_in_fov(m.Data, C.int(x), C.int(y)))
}

func (m *Map) SetInFov(x, y int, fov bool) {
	C.TCOD_map_set_in_fov(m.Data, C.int(x), C.int(y), fromBool(fov))
}

// retrieve properties from the map

func (m *Map) IsTransparent(x, y int) bool {
	return toBool(C.TCOD_map_is_transparent(m.Data, C.int(x), C.int(y)))
}

func (m *Map) IsWalkable(x, y int) bool {
	return toBool(C.TCOD_map_is_walkable(m.Data, C.int(x), C.int(y)))
}

func (m *Map) GetWidth() int {
	return int(C.TCOD_map_get_width(m.Data))
}

func (m *Map) GetHeight() int {
	return int(C.TCOD_map_get_height(m.Data))
}

func (m *Map) GetNbCells() int {
	return int(C.TCOD_map_get_nb_cells(m.Data))
}

//
// BSP Dungeon generation
//
//
type Bsp struct {
	X, Y, W, H         int   // node position & size
	Position           int   // position of splitting
	Level              uint8 // level in the tree
	Horizontal         bool  // horizontal splitting ?
	next, father, sons *Bsp  // BSP tree hierarchy structuring
}

type BspListener func(node *Bsp, userData interface{}) bool

func (bsp *Bsp) AddSon(son *Bsp) {
	lastson := bsp.sons
	son.father = bsp
	for lastson != nil && lastson.next != nil {
		lastson = lastson.next
	}
	if lastson != nil {
		lastson.next = son
	} else {
		bsp.sons = son
	}
}

func NewBspWithSize(x, y, w, h int) (result *Bsp) {
	result = new(Bsp)
	*result = Bsp{X: x, Y: y, W: w, H: h}
	return
}

func (bsp *Bsp) Left() *Bsp {
	return bsp.sons
}

func (bsp *Bsp) Right() *Bsp {
	if bsp.sons != nil {
		return bsp.sons.next
	} else {
		return nil
	}
}

func (bsp *Bsp) Father() *Bsp {
	return bsp.father
}

func (bsp *Bsp) IsLeaf() bool {
	return bsp.sons == nil
}

func NewBspIntern(father *Bsp, left bool) *Bsp {
	bsp := new(Bsp)
	if father.Horizontal {
		bsp.X = father.X
		bsp.W = father.W
		if left {
			bsp.Y = father.Y
		} else {
			bsp.Y = father.Position
		}
		if left {
			bsp.H = father.Position - bsp.Y
		} else {
			bsp.H = father.Y + father.H - father.Position
		}
	} else {
		bsp.Y = father.Y
		bsp.H = father.H
		if left {
			bsp.X = father.X
		} else {
			bsp.X = father.Position
		}
		if left {
			bsp.W = father.Position - bsp.X
		} else {
			bsp.W = father.X + father.W - father.Position
		}
	}
	bsp.Level = father.Level + 1
	return bsp
}

func (bsp *Bsp) TraversePreOrder(listener BspListener, userData interface{}) bool {
	if !listener(bsp, userData) {
		return false
	}
	if bsp.Left() != nil && !bsp.Left().TraversePreOrder(listener, userData) {
		return false
	}
	if bsp.Right() != nil && !bsp.Right().TraversePreOrder(listener, userData) {
		return false
	}
	return true
}

func (bsp *Bsp) TraverseInOrder(listener BspListener, userData interface{}) bool {
	if bsp.Left() != nil && !bsp.Left().TraverseInOrder(listener, userData) {
		return false
	}
	if !listener(bsp, userData) {
		return false
	}
	if bsp.Right() != nil && !bsp.Right().TraverseInOrder(listener, userData) {
		return false
	}
	return true
}

func (bsp *Bsp) TraversePostOrder(listener BspListener, userData interface{}) bool {
	if bsp.Left() != nil && !bsp.Left().TraversePostOrder(listener, userData) {
		return false
	}
	if bsp.Right() != nil && !bsp.Right().TraversePostOrder(listener, userData) {
		return false
	}
	if !listener(bsp, userData) {
		return false
	}
	return true
}

func (bsp *Bsp) TraverseLevelOrder(listener BspListener, userData interface{}) bool {
	stack := []*Bsp{bsp}
	for len(stack) > 0 {
		node := stack[0]
		stack = stack[1:]
		if node.Left() != nil {
			stack = append(stack, node.Left())
		}
		if node.Right() != nil {
			stack = append(stack, node.Right())
		}
		if !listener(node, userData) {
			return false
		}
	}
	return true
}

// TODO can it store Go values in list structure??
// maybe replace it with record
func (bsp *Bsp) TraverseInvertedLevelOrder(listener BspListener, userData interface{}) bool {
	stack1 := []*Bsp{bsp}
	stack2 := []*Bsp{}
	for len(stack1) > 0 {
		node := stack1[0]
		stack1 = stack1[1:]
		stack2 = append(stack2, node)
		if node.Left() != nil {
			stack1 = append(stack1, node.Left())
		}
		if node.Right() != nil {
			stack1 = append(stack1, node.Right())
		}
	}
	for len(stack2) > 0 {
		node := stack2[len(stack2)-1]
		stack2 = stack2[:len(stack2)-1]
		if !listener(node, userData) {
			return false
		}
	}
	return true
}

func (bsp *Bsp) RemoveSons() {
	node := bsp.sons
	var nextNode *Bsp
	for node != nil {
		nextNode = node.next
		node.RemoveSons()
		node = nextNode
	}
	bsp.sons = nil
}

func (bsp *Bsp) SplitOnce(horizontal bool, position int) {
	bsp.Horizontal = horizontal
	bsp.Position = position
	bsp.AddSon(NewBspIntern(bsp, true))
	bsp.AddSon(NewBspIntern(bsp, false))
}

func (bsp *Bsp) SplitRecursive(randomizer *Random, nb int, minHSize int, minVSize int, maxHRatio float32, maxVRatio float32) {
	var horiz bool
	var position int
	if nb == 0 || (bsp.W < 2*minHSize && bsp.H < 2*minVSize) {
		return
	}
	// promote square rooms
	if bsp.H < 2*minVSize || float32(bsp.W) > float32(bsp.H)*maxHRatio {
		horiz = false
	} else if bsp.W < 2*minHSize || float32(bsp.H) > float32(bsp.W)*maxVRatio {
		horiz = true
	} else {
		horiz = (randomizer.GetInt(0, 1) == 0)
	}
	if horiz {
		position = randomizer.GetInt(bsp.Y+minVSize, bsp.Y+bsp.H-minVSize)
	} else {
		position = randomizer.GetInt(bsp.X+minHSize, bsp.X+bsp.W-minHSize)
	}
	bsp.SplitOnce(horiz, position)
	if bsp.Left() != nil {
		bsp.Left().SplitRecursive(randomizer, nb-1, minHSize, minVSize, maxHRatio, maxVRatio)
	}
	if bsp.Right() != nil {
		bsp.Right().SplitRecursive(randomizer, nb-1, minHSize, minVSize, maxHRatio, maxVRatio)
	}
}

func (bsp *Bsp) Resize(x, y, w, h int) {
	bsp.X, bsp.Y, bsp.W, bsp.H = x, y, w, h
	if bsp.Left() != nil {
		if bsp.Horizontal {
			bsp.Left().Resize(x, y, w, bsp.Position-y)
			if bsp.Right() != nil {
				bsp.Right().Resize(x, bsp.Position, w, y+h-bsp.Position)
			}
		} else {
			bsp.Left().Resize(x, y, bsp.Position-x, h)
			if bsp.Right() != nil {
				bsp.Right().Resize(bsp.Position, y, x+w-bsp.Position, h)
			}
		}
	}
}

func (bsp *Bsp) Contains(x, y int) bool {
	return x >= bsp.X && y >= bsp.Y && x < bsp.X+bsp.W && y < bsp.Y+bsp.H
}

func (bsp *Bsp) FindNode(x, y int) *Bsp {
	if !bsp.Contains(x, y) {
		return nil
	}
	if !bsp.IsLeaf() {
		var left, right *Bsp
		left = bsp.Left()
		if left.Contains(x, y) {
			return left.FindNode(x, y)
		}
		right = bsp.Right()
		if right.Contains(x, y) {
			return right.FindNode(x, y)
		}
	}
	return bsp
}

//
// HeightMap
//

type HeightMap struct {
	Data *C.TCOD_heightmap_t
}

func deleteHeightmap(h *HeightMap) {
	C.TCOD_heightmap_delete(h.Data)
}

func NewHeightMap(w, h int) *HeightMap {
	result := &HeightMap{C.TCOD_heightmap_new(C.int(w), C.int(h))}
	runtime.SetFinalizer(result, deleteHeightmap)
	return result
}

func (heightMap *HeightMap) GetValue(x, y int) float32 {
	return float32(C.TCOD_heightmap_get_value(heightMap.Data, C.int(x), C.int(y)))
}

func (heightMap *HeightMap) GetWidth() int {
	return int(heightMap.Data.w)
}

func (heightMap *HeightMap) GetHeight() int {
	return int(heightMap.Data.h)
}

func (heightMap *HeightMap) GetInterpolatedValue(x, y float32) float32 {
	return float32(C.TCOD_heightmap_get_interpolated_value(heightMap.Data, C.float(x), C.float(y)))
}

func (heightMap *HeightMap) SetValue(x, y int, value float32) {
	C.TCOD_heightmap_set_value(heightMap.Data, C.int(x), C.int(y), C.float(value))
}

func (heightMap *HeightMap) GetNthValue(nth int) float32 {
	return float32(C._TCOD_heightmap_get_nth_value(heightMap.Data, C.int(nth)))
}

func (heightMap *HeightMap) SetNthValue(nth int, value float32) {
	C._TCOD_heightmap_set_nth_value(heightMap.Data, C.int(nth), C.float(value))
}

func (heightMap *HeightMap) GetSlope(x, y int) float32 {
	return float32(C.TCOD_heightmap_get_slope(heightMap.Data, C.int(x), C.int(y)))
}

func (heightMap *HeightMap) GetNormal(x, y float32, n *[3]float32, waterLevel float32) {
	C.TCOD_heightmap_get_normal(heightMap.Data, C.float(x), C.float(y),
		(*C.float)(unsafe.Pointer(&n[0])),
		C.float(waterLevel))
}

func (heightMap *HeightMap) CountCells(min, max float32) int {
	return int(C.TCOD_heightmap_count_cells(heightMap.Data, C.float(min), C.float(max)))
}

func (heightMap *HeightMap) HasLandOnBorder(waterLevel float32) bool {
	return toBool(C.TCOD_heightmap_has_land_on_border(heightMap.Data, C.float(waterLevel)))
}

func (heightMap *HeightMap) GetMinMax() (min, max float32) {
	var cmin, cmax C.float
	C.TCOD_heightmap_get_minmax(heightMap.Data, &cmin, &cmax)
	min, max = float32(cmin), float32(cmax)
	return
}

func (heightMap *HeightMap) Copy(source *HeightMap) {
	C.TCOD_heightmap_copy(source.Data, heightMap.Data)
}

func (heightMap *HeightMap) Add(value float32) {
	C.TCOD_heightmap_add(heightMap.Data, C.float(value))
}

func (heightMap *HeightMap) Scale(value float32) {
	C.TCOD_heightmap_scale(heightMap.Data, C.float(value))
}

func (heightMap *HeightMap) Clamp(min, max float32) {
	C.TCOD_heightmap_clamp(heightMap.Data, C.float(min), C.float(max))
}

func (heightMap *HeightMap) Normalize() {
	heightMap.NormalizeRange(0, 1)
}

func (heightMap *HeightMap) NormalizeRange(min, max float32) {
	C.TCOD_heightmap_normalize(heightMap.Data, C.float(min), C.float(max))
}

func (heightMap *HeightMap) Clear() {
	C.TCOD_heightmap_clear(heightMap.Data)
}

func (heightMap *HeightMap) Lerp(hm1 *HeightMap, hm2 *HeightMap, coef float32) {
	C.TCOD_heightmap_lerp_hm(hm1.Data, hm2.Data, heightMap.Data, C.float(coef))
}

func (heightMap *HeightMap) AddHm(hm1 *HeightMap, hm2 *HeightMap) {
	C.TCOD_heightmap_add_hm(hm1.Data, hm2.Data, heightMap.Data)
}

func (heightMap *HeightMap) Multiply(hm1 *HeightMap, hm2 *HeightMap) {
	C.TCOD_heightmap_multiply_hm(hm1.Data, hm2.Data, heightMap.Data)
}

func (heightMap *HeightMap) AddHill(hx, hy, hradius, hheight float32) {
	C.TCOD_heightmap_add_hill(heightMap.Data, C.float(hx), C.float(hy), C.float(hradius), C.float(hheight))
}

func (heightMap *HeightMap) DigHill(hx, hy, hradius, hheight float32) {
	C.TCOD_heightmap_dig_hill(heightMap.Data, C.float(hx), C.float(hy), C.float(hradius), C.float(hheight))
}

func (heightMap *HeightMap) DigBezier(px, py *[4]int, startRadius, startDepth, endRadius, endDepth float32) {
	C.TCOD_heightmap_dig_bezier(heightMap.Data,
		(*C.int)(unsafe.Pointer(&px[0])),
		(*C.int)(unsafe.Pointer(&py[0])),
		C.float(startRadius), C.float(startDepth), C.float(endRadius), C.float(endDepth))
}

func (heightMap *HeightMap) RainErosion(nbDrops int, erosionCoef, sedimentationCoef float32, rnd *Random) {
	C.TCOD_heightmap_rain_erosion(heightMap.Data, C.int(nbDrops), C.float(erosionCoef), C.float(sedimentationCoef), rnd.Data)
}

func (heightMap *HeightMap) KernelTransform(kernelsize int, dx, dy []int, weight []float32, minLevel, maxLevel float32) {
	C.TCOD_heightmap_kernel_transform(heightMap.Data, C.int(kernelsize),
		(*C.int)(unsafe.Pointer(&dx[0])),
		(*C.int)(unsafe.Pointer(&dy[0])),
		(*C.float)(unsafe.Pointer(&weight[0])),
		C.float(minLevel),
		C.float(maxLevel))
}

func (heightMap *HeightMap) AddVoronoi(nbPoints, nbCoef int, coef []float32, rnd *Random) {
	C.TCOD_heightmap_add_voronoi(heightMap.Data, C.int(nbPoints), C.int(nbCoef), (*C.float)(unsafe.Pointer(&coef[0])), rnd.Data)
}

func (heightMap *HeightMap) AddFbm(noise *Noise, mulx, muly, addx, addy, octaves, delta, scale float32) {
	C.TCOD_heightmap_add_fbm(heightMap.Data, noise.Data, C.float(mulx),
		C.float(muly), C.float(addx), C.float(addy), C.float(octaves), C.float(delta), C.float(scale))
}

func (heightMap *HeightMap) ScaleFbm(noise *Noise, mulx, muly, addx, addy, octaves, delta, scale float32) {
	C.TCOD_heightmap_scale_fbm(heightMap.Data, noise.Data, C.float(mulx),
		C.float(muly), C.float(addx), C.float(addy), C.float(octaves), C.float(delta), C.float(scale))
}

func (heightMap *HeightMap) Islandify(seaLevel float32, random *Random) {
	C.TCOD_heightmap_islandify(heightMap.Data, C.float(seaLevel), random.Data)
}

//
// Image
//

type Image struct {
	Data C.TCOD_image_t
}

func deleteImage(img *Image) {
	C.TCOD_image_delete(img.Data)
}

func newImage(data C.TCOD_image_t) *Image {
	result := &Image{data}
	runtime.SetFinalizer(result, deleteImage)
	return result
}

func NewImage(width, height int) *Image {
	return newImage(C.TCOD_image_new(C.int(width), C.int(height)))
}

func NewImageFromConsole(console *Console) *Image {
	return newImage(C.TCOD_image_from_console(console.Data))
}

func (image *Image) RefreshConsole(console *Console) {
	C.TCOD_image_refresh_console(image.Data, console.Data)
}

func LoadImage(filename string) *Image {
	return newImage(C.TCOD_image_load(C.CString(filename)))
}

func (image *Image) Clear(color Color) {
	ccolor := fromColor(color)
	C.TCOD_image_clear(image.Data, ccolor)
}

func (image *Image) Invert() {
	C.TCOD_image_invert(image.Data)
}

func (image *Image) Hflip() {
	C.TCOD_image_hflip(image.Data)
}

func (image *Image) Rotate90(numRotations int) {
	C.TCOD_image_rotate90(image.Data, C.int(numRotations))
}

func (image *Image) Vflip() {
	C.TCOD_image_vflip(image.Data)
}

func (image *Image) Scale(neww, newh int) {
	C.TCOD_image_scale(image.Data, C.int(neww), C.int(newh))
}

func (image *Image) Save(filename string) {
	C.TCOD_image_save(image.Data, C.CString(filename))
}

func (image *Image) GetSize(w, h *int) {
	var cw, ch C.int
	C.TCOD_image_get_size(image.Data, &cw, &ch)
	*w = int(cw)
	*h = int(ch)
}

func (image *Image) GetPixel(x, y int) Color {
	return toColor(C.TCOD_image_get_pixel(image.Data, C.int(x), C.int(y)))
}

func (image *Image) GetAlpha(x, y int) int {
	return int(C.TCOD_image_get_alpha(image.Data, C.int(x), C.int(y)))
}

func (image *Image) GetMipmapPixel(x0, y0, x1, y1 float32) Color {
	return toColor(C.TCOD_image_get_mipmap_pixel(image.Data, C.float(x0), C.float(y0),
		C.float(x1), C.float(y1)))
}

func (image *Image) PutPixel(x, y int, color Color) {
	ccolor := fromColor(color)
	C.TCOD_image_put_pixel(image.Data, C.int(x), C.int(y), ccolor)
}

func (image *Image) Blit(console *Console, x, y float32, bkgndFlag BkgndFlag, scalex, scaley, angle float32) {
	C.TCOD_image_blit(image.Data, console.Data, C.float(x), C.float(y),
		C.TCOD_bkgnd_flag_t(bkgndFlag), C.float(scalex), C.float(scaley), C.float(angle))
}

func (image *Image) BlitRect(console *Console, x, y, w, h int, flag BkgndFlag) {
	C.TCOD_image_blit_rect(image.Data, console.Data, C.int(x), C.int(y), C.int(w), C.int(h), C.TCOD_bkgnd_flag_t(flag))
}

func (image *Image) Blit2x(dest *Console, dx, dy, sx, sy, w, h int) {
	C.TCOD_image_blit_2x(image.Data, dest.Data, C.int(dx), C.int(dy), C.int(sx), C.int(sy), C.int(w), C.int(h))
}

func (image *Image) SetKeyColor(keyColor Color) {
	ckeyColor := fromColor(keyColor)
	C.TCOD_image_set_key_color(image.Data, ckeyColor)
}

func (image *Image) IsPixelTransparent(x, y int) bool {
	return toBool(C.TCOD_image_is_pixel_transparent(image.Data, C.int(x), C.int(y)))
}

//
// Path
//
//
type Path struct {
	Data C.TCOD_path_t
}

func deletePath(path *Path) {
	C.TCOD_path_delete(path.Data)
}

func NewPathUsingMap(m *Map, diagonalCost float32) *Path {
	result := &Path{C.TCOD_path_new_using_map(m.Data, C.float(diagonalCost))}
	runtime.SetFinalizer(result, deletePath)
	return result
}

// Not implemented - go not supporting callbacks
// func PathNewUsingFunction() {
//	//TCODLIB_API TCOD_path_t
//  TCOD_path_new_using_function(int map_width, int map_height, TCOD_path_func_t func, void *user_Data, float diagonalCost);
// }

func (path *Path) Compute(ox, oy, dx, dy int) bool {
	return toBool(C.TCOD_path_compute(path.Data, C.int(ox), C.int(oy), C.int(dx), C.int(dy)))
}

func (path *Path) Walk(recalcWhenNeeded bool) (x, y int) {
	var cx, cy C.int
	C.TCOD_path_walk(path.Data, &cx, &cy, fromBool(recalcWhenNeeded))
	x, y = int(cx), int(cy)
	return
}

func (path *Path) IsEmpty() bool {
	return toBool(C.TCOD_path_is_empty(path.Data))
}

func (path *Path) Size() int {
	return int(C.TCOD_path_size(path.Data))
}

func (path *Path) Get(index int) (x, y int) {
	var cx, cy C.int
	C.TCOD_path_get(path.Data, C.int(index), &cx, &cy)
	x, y = int(cx), int(cy)
	return
}

func (path *Path) GetOrigin() (x, y int) {
	var cx, cy C.int
	C.TCOD_path_get_origin(path.Data, &cx, &cy)
	x, y = int(cx), int(cy)
	return
}

func (path *Path) GetDestination() (x, y int) {
	var cx, cy C.int
	C.TCOD_path_get_destination(path.Data, &cx, &cy)
	x, y = int(cx), int(cy)
	return
}

//
// Dijkstra path
//

type Dijkstra struct {
	Data C.TCOD_dijkstra_t
}

func deleteDijkstra(d *Dijkstra) {
	C.TCOD_dijkstra_delete(d.Data)
}

func NewDijkstraUsingMap(m *Map, diagonalCost float32) *Dijkstra {
	result := &Dijkstra{C.TCOD_dijkstra_new(m.Data, C.float(diagonalCost))}
	runtime.SetFinalizer(result, deleteDijkstra)
	return result
}

// Not implemented - go not supporting callbacks
// func DijkstraNewUsingFunction() {
//	//TCODLIB_API TCOD_Dijkstra_t
//   TCOD_Dijkstra_new_using_function(int map_width, int map_height, TCOD_Dijkstra_func_t func, void *user_Data, float diagonalCost);
// }

func (dijkstra *Dijkstra) Compute(rootX, rootY int) {
	C.TCOD_dijkstra_compute(dijkstra.Data, C.int(rootX), C.int(rootY))
}

func (dijkstra *Dijkstra) GetDistance(x, y int) float32 {
	return float32(C.TCOD_dijkstra_get_distance(dijkstra.Data, C.int(x), C.int(y)))
}

func (dijkstra *Dijkstra) PathSet(x, y int) bool {
	return toBool(C.TCOD_dijkstra_path_set(dijkstra.Data, C.int(x), C.int(y)))
}

func (dijkstra *Dijkstra) IsEmpty() bool {
	return toBool(C.TCOD_dijkstra_is_empty(dijkstra.Data))
}

func (dijkstra *Dijkstra) Size() int {
	return int(C.TCOD_dijkstra_size(dijkstra.Data))
}

func (dijkstra *Dijkstra) Get(index int) (x, y int) {
	var cx, cy C.int
	C.TCOD_dijkstra_get(dijkstra.Data, C.int(index), &cx, &cy)
	x, y = int(cx), int(cy)
	return
}

func (dijkstra *Dijkstra) PathWalk() (x, y int) {
	var cx, cy C.int
	C.TCOD_dijkstra_path_walk(dijkstra.Data, &cx, &cy)
	x, y = int(cx), int(cy)
	return
}

//
// Mersenne Random generator
//

type RandomAlgo C.TCOD_random_algo_t

type Distribution C.TCOD_distribution_t

type Random struct {
	Data C.TCOD_random_t
}

type Dice struct {
	Data C.TCOD_dice_t
}

func fromDice(d Dice) C.TCOD_dice_t {
	return d.Data
}

func toDice(d C.TCOD_dice_t) Dice {
	return Dice{d}
}

func deleteRandom(r *Random) {
	C.TCOD_random_delete(r.Data)
}

func newRandom(data C.TCOD_random_t) *Random {
	result := &Random{data}
	runtime.SetFinalizer(result, deleteRandom)
	return result
}

func GetRandomInstance() *Random {
	return newRandom(C.TCOD_random_get_instance())
}

func NewRandom() *Random {
	return newRandom(C.TCOD_random_new(C.TCOD_random_algo_t(RNG_MT)))
}

func NewRandomWithAlgo(algo RandomAlgo) *Random {
	return newRandom(C.TCOD_random_new(C.TCOD_random_algo_t(algo)))
}

func NewRandomFromSeedWithAlgo(seed uint32, algo RandomAlgo) *Random {
	return newRandom(C.TCOD_random_new_from_seed(C.TCOD_random_algo_t(algo), C.uint32_t(seed)))
}

func NewRandomFromSeed(seed uint32) *Random {
	return newRandom(
		C.TCOD_random_new_from_seed(
			C.TCOD_random_algo_t(RNG_MT),
			C.uint32_t(seed)))
}

func (random *Random) Save() *Random {
	result := newRandom(C.TCOD_random_save(random.Data))
	return result
}

func (random *Random) Restore(backup *Random) {
	C.TCOD_random_restore(random.Data, backup.Data)
}

func (random *Random) SetDistribution(distribution Distribution) {
	C.TCOD_random_set_distribution(random.Data, C.TCOD_distribution_t(distribution))
}

func (random *Random) GetInt(min, max int) int {
	return int(C.TCOD_random_get_int(random.Data, C.int(min), C.int(max)))
}

func (random *Random) GetFloat(min, max float32) float32 {
	return float32(C.TCOD_random_get_float(random.Data, C.float(min), C.float(max)))
}

func (random *Random) GetDouble(min, max float64) float64 {
	return float64(C.TCOD_random_get_double(random.Data, C.double(min), C.double(max)))
}

func (random *Random) GetIntMean(min, max, mean int) int {
	return int(C.TCOD_random_get_int_mean(random.Data, C.int(min), C.int(max), C.int(mean)))
}

func (random *Random) GetFloatMean(min, max, mean float32) float32 {
	return float32(C.TCOD_random_get_float_mean(random.Data, C.float(min), C.float(max), C.float(mean)))
}

func (random *Random) GetDoubleMean(min, max, mean float64) float64 {
	return float64(C.TCOD_random_get_double_mean(random.Data, C.double(min), C.double(max), C.double(mean)))
}

func NewDice(s string) *Dice {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	result := &Dice{Data: C.TCOD_random_dice_new(cs)}
	return result
}

func (self *Dice) Roll(random *Random) int {
	return int(C.TCOD_random_dice_roll(random.Data, self.Data))
}

func RollDice(random *Random, s string) int {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C.TCOD_random_dice_roll_s(random.Data, cs))
}

//
// Parser library
//

type ParserValueType C.TCOD_value_type_t

type ParserStruct struct {
	Data C.TCOD_parser_struct_t
}

type Parser struct {
	Data C.TCOD_parser_t
}

type ParserProperty struct {
	Name      string
	ValueType ParserValueType
	Value     interface{}
}

func (ps ParserStruct) GetName() string {
	return C.GoString(C.TCOD_struct_get_name(ps.Data))
}

func (ps ParserStruct) AddProperty(name string, valueType ParserValueType, mandatory bool) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	C.TCOD_struct_add_property(ps.Data, cname, C.TCOD_value_type_t(valueType), fromBool(mandatory))
}

func (ps ParserStruct) AddListProperty(name string, valueType ParserValueType, mandatory bool) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.TCOD_struct_add_list_property(ps.Data, cname, C.TCOD_value_type_t(valueType), fromBool(mandatory))
}

func (ps ParserStruct) AddValueList(name string, valueList []string, mandatory bool) {
	cvalueList := make([]*C.char, len(valueList))
	for i := range valueList {
		cvalueList[i] = C.CString(valueList[i])
	}
	C.TCOD_struct_add_value_list_sized(ps.Data, C.CString(name),
		(**C.char)(unsafe.Pointer(&cvalueList[0])), C.int(len(valueList)), fromBool(mandatory))

	for i := range cvalueList {
		C.free(unsafe.Pointer(cvalueList[i]))
	}

}

func (ps ParserStruct) AddFlag(propname string) {
	cpropname := C.CString(propname)
	defer C.free(unsafe.Pointer(cpropname))
	C.TCOD_struct_add_flag(ps.Data, cpropname)
}

func (ps ParserStruct) AddStructure(substruct ParserStruct) {
	// TODO is this necessary ??
	//	struct1 := ps.Data
	//	substruct2 := struct_.Data
	C.TCOD_struct_add_structure(ps.Data, substruct.Data)
}

func (ps *ParserStruct) IsMandatory(propname string) bool {
	cpropname := C.CString(propname)
	defer C.free(unsafe.Pointer(cpropname))
	return toBool(C.TCOD_struct_is_mandatory(ps.Data, cpropname))
}

func (ps *ParserStruct) GetType(propname string) ParserValueType {
	cpropname := C.CString(propname)
	defer C.free(unsafe.Pointer(cpropname))
	return ParserValueType(C.TCOD_struct_get_type(ps.Data, cpropname))
}

func deleteParser(p *Parser) {
	C.TCOD_parser_delete(p.Data)
}

func NewParser() *Parser {
	result := &Parser{C.TCOD_parser_new()}
	runtime.SetFinalizer(result, deleteParser)
	return result
}

func (parser *Parser) RegisterStruct(name string) ParserStruct {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return ParserStruct{C.TCOD_parser_new_struct(parser.Data, cname)}
}

// TODO custom parsers are not supported
// TCODLIB_API TCOD_value_type_t TCOD_parser_new_custom_type(TCOD_parser_t parser,TCOD_parser_custom_t custom_type_parser);

// TODO listeners are not supported
// Running parser return list of parsed properties
func (parser *Parser) Run(filename string) []ParserProperty {
	// run parser with default listeners
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	C.TCOD_parser_run(parser.Data, cfilename, nil)

	// extract properties to Go structures
	var cprop *C._prop_t
	var prop ParserProperty
	var l C.TCOD_list_t = C.TCOD_list_t(((*C.TCOD_parser_int_t)(parser.Data)).props)
	result := make([]ParserProperty, C.TCOD_list_size(l))

	for i := 0; i < int(C.TCOD_list_size(l)); i++ {

		cprop = (*C._prop_t)(unsafe.Pointer(C.TCOD_list_get(l, C.int(i))))

		prop.Name = C.GoString(cprop.name)
		prop.ValueType = ParserValueType(cprop.value_type)
		if cprop.value_type == TYPE_STRING ||
			(cprop.value_type >= TYPE_VALUELIST00 && cprop.value_type <= TYPE_VALUELIST15) {
			prop.Value = C.GoString(*((**C.char)(unsafe.Pointer(&cprop.value))))
		} else if cprop.value_type == TYPE_INT {
			prop.Value = int(*((*C.int)(unsafe.Pointer(&cprop.value))))
		} else if cprop.value_type == TYPE_FLOAT {
			prop.Value = float32(*((*C.float)(unsafe.Pointer(&cprop.value))))
		} else if cprop.value_type == TYPE_BOOL {
			prop.Value = toBool(*((*C.bool)(unsafe.Pointer(&cprop.value))))
		} else if cprop.value_type == TYPE_COLOR {
			prop.Value = toColor(*((*C.TCOD_color_t)(unsafe.Pointer(&cprop.value))))
		} else if cprop.value_type == TYPE_DICE {
			prop.Value = toDice(*((*C.TCOD_dice_t)(unsafe.Pointer(&cprop.value))))
		} else if cprop.value_type >= TYPE_LIST {
			elType := cprop.value_type - TYPE_LIST
			elList := C.TCOD_list_t(*(*C.TCOD_list_t)(unsafe.Pointer(&cprop.value)))
			elListSize := int(C.TCOD_list_size(elList))
			if elType == TYPE_STRING {
				prop.Value = make([]string, elListSize)
				for j := 0; j < elListSize; j++ {
					elValue := (*C.char)(unsafe.Pointer(C.TCOD_list_get(elList, C.int(j))))
					prop.Value.([]string)[j] = C.GoString(elValue)
				}
			} else if elType == TYPE_INT {
				prop.Value = make([]int, elListSize)
				for j := 0; j < elListSize; j++ {
					elValue := C.TCOD_list_get(elList, C.int(j))
					prop.Value.([]int)[j] = int(*(*C.int)(unsafe.Pointer(&elValue)))
				}
			} else if elType == TYPE_FLOAT {
				prop.Value = make([]float32, elListSize)
				for j := 0; j < elListSize; j++ {
					elValue := C.TCOD_list_get(elList, C.int(j))
					prop.Value.([]float32)[j] = float32(*(*C.float)(unsafe.Pointer(&elValue)))
				}
			} else if elType == TYPE_BOOL {
				prop.Value = make([]bool, elListSize)
				for j := 0; j < elListSize; j++ {
					elValue := C.TCOD_list_get(elList, C.int(j))
					prop.Value.([]bool)[j] = toBool(*(*C.bool)(unsafe.Pointer(&elValue)))
				}
			} else if elType == TYPE_DICE {
				prop.Value = make([]Dice, elListSize)
				for j := 0; j < elListSize; j++ {
					elValue := *(*C.TCOD_dice_t)(unsafe.Pointer(C.TCOD_list_get(elList, C.int(j))))
					prop.Value.([]Dice)[j] = toDice(elValue)
				}
			} else if elType == TYPE_COLOR {
				prop.Value = make([]Color, elListSize)
				for j := 0; j < elListSize; j++ {
					elValue := *(*C.TCOD_color_t)(unsafe.Pointer(C.TCOD_list_get(elList, C.int(j))))
					prop.Value.([]Color)[j] = toColor(elValue)
				}
			}
		}
		result[i] = prop
	}
	return result
}

//
// Perlin noise
//

// Noise NEW

const NOISE_MAX_OCTAVES = 128
const NOISE_MAX_DIMENSIONS = 4
const NOISE_DEFAULT_HURST = 0.5
const NOISE_DEFAULT_LACUNARITY = 2.0

type NoiseType C.TCOD_noise_type_t

type Noise struct {
	Data C.TCOD_noise_t
}

type FloatArray []float32

func deleteNoise(noise *Noise) {
	C.TCOD_noise_delete(noise.Data)
}

func newNoise(d C.TCOD_noise_t) *Noise {
	result := &Noise{d}
	runtime.SetFinalizer(result, deleteNoise)
	return result
}

func NewNoise(dimensions int, random *Random) *Noise {
	return newNoise(C.TCOD_noise_new(C.int(dimensions), C.float(NOISE_DEFAULT_HURST),
		C.float(NOISE_DEFAULT_LACUNARITY), random.Data))
}

func NewNoiseWithOptions(dimensions int, hurst float32, lacunarity float32, random *Random) *Noise {
	return newNoise(C.TCOD_noise_new(C.int(dimensions), C.float(hurst), C.float(lacunarity), random.Data))
}

func (noise *Noise) GetEx(f FloatArray, noiseType NoiseType) float32 {
	return float32(C.TCOD_noise_get_ex(noise.Data, (*C.float)(unsafe.Pointer(&f[0])), C.TCOD_noise_type_t(noiseType)))
}

func (noise *Noise) SetType(noiseType NoiseType) {
	C.TCOD_noise_set_type(noise.Data, C.TCOD_noise_type_t(noiseType))
}

func (noise *Noise) GetFbmEx(f FloatArray, octaves float32, noiseType NoiseType) float32 {
	return float32(C.TCOD_noise_get_fbm_ex(noise.Data, (*C.float)(unsafe.Pointer(&f[0])), C.float(octaves),
		C.TCOD_noise_type_t(noiseType)))
}

func (noise *Noise) GetTurbulenceEx(f FloatArray, octaves float32, noiseType NoiseType) float32 {
	return float32(C.TCOD_noise_get_turbulence_ex(noise.Data, (*C.float)(unsafe.Pointer(&f[0])), C.float(octaves),
		C.TCOD_noise_type_t(noiseType)))
}

func (noise *Noise) Get(f FloatArray) float32 {
	return float32(C.TCOD_noise_get(noise.Data, (*C.float)(unsafe.Pointer(&f[0]))))
}

func (noise *Noise) GetFbm(f FloatArray, octaves float32) float32 {
	return float32(C.TCOD_noise_get_fbm(noise.Data, (*C.float)(unsafe.Pointer(&f[0])), C.float(octaves)))
}

func (noise *Noise) GetTurbulence(f FloatArray, octaves float32) float32 {
	return float32(C.TCOD_noise_get_turbulence(noise.Data, (*C.float)(unsafe.Pointer(&f[0])), C.float(octaves)))
}

//
// Zip
//

type Zip struct {
	Data C.TCOD_zip_t
}

func deleteZip(zip *Zip) {
	C.TCOD_zip_delete(zip.Data)
}

func NewZip() *Zip {
	result := &Zip{C.TCOD_zip_new()}
	runtime.SetFinalizer(result, deleteZip)
	return result
}

// output interface

func (zip *Zip) PutChar(val byte) {
	C.TCOD_zip_put_char(zip.Data, C.char(val))
}

func (zip *Zip) PutInt(val int) {
	C.TCOD_zip_put_int(zip.Data, C.int(val))
}

func (zip *Zip) PutFloat(val float32) {
	C.TCOD_zip_put_float(zip.Data, C.float(val))
}

func (zip *Zip) PutString(val string) {
	cval := C.CString(val)
	defer C.free(unsafe.Pointer(cval))
	C.TCOD_zip_put_string(zip.Data, cval)
}

func (zip *Zip) PutColor(val Color) {
	cval := fromColor(val)
	C.TCOD_zip_put_color(zip.Data, cval)
}

func (zip *Zip) PutImage(val *Image) {
	C.TCOD_zip_put_image(zip.Data, val.Data)
}

func (zip *Zip) PutConsole(val *Console) {
	C.TCOD_zip_put_console(zip.Data, val.Data)
}

func (zip *Zip) PutData(nbBytes int, data unsafe.Pointer) {
	C.TCOD_zip_put_data(zip.Data, C.int(nbBytes), data)
}

func (zip *Zip) GetCurrentBytes() uint32 {
	return uint32(C.TCOD_zip_get_current_bytes(zip.Data))
}

func (zip *Zip) SaveToFile(filename string) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	C.TCOD_zip_save_to_file(zip.Data, cfilename)
}

// input interface

func (zip *Zip) LoadFromFile(filename string) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	C.TCOD_zip_load_from_file(zip.Data, cfilename)
}

func (zip *Zip) GetChar() byte {
	return byte(C.TCOD_zip_get_char(zip.Data))
}

func (zip *Zip) GetInt() int {
	return int(C.TCOD_zip_get_int(zip.Data))
}

func (zip *Zip) GetFloat() float32 {
	return float32(C.TCOD_zip_get_float(zip.Data))
}

func (zip *Zip) GetString() string {
	return C.GoString(C.TCOD_zip_get_string(zip.Data))
}

func (zip *Zip) GetColor() Color {
	return toColor(C.TCOD_zip_get_color(zip.Data))
}

func (zip *Zip) GetImage() *Image {
	return &Image{C.TCOD_zip_get_image(zip.Data)}
}

func (zip *Zip) GetConsole() *Console {
	return &Console{C.TCOD_zip_get_console(zip.Data)}
}

func (zip *Zip) GetData(nbBytes int, data unsafe.Pointer) int {
	return int(C.TCOD_zip_get_data(zip.Data, C.int(nbBytes), data))
}

func (zip *Zip) GetRemainingBytes() uint32 {
	return uint32(C.TCOD_zip_get_remaining_bytes(zip.Data))
}

func (zip *Zip) SkipBytes(nbBytes uint32) {
	C.TCOD_zip_skip_bytes(zip.Data, C.uint32_t(nbBytes))
}
