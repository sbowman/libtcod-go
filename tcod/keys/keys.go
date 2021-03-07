package keys

/*
 #include "../include/libtcod.h"
*/
import "C"

const (
	Zero        = C.TCODK_0
	One         = C.TCODK_1
	Two         = C.TCODK_2
	Three       = C.TCODK_3
	Four        = C.TCODK_4
	Five        = C.TCODK_5
	Six         = C.TCODK_6
	Seven       = C.TCODK_7
	Eight       = C.TCODK_8
	Nine        = C.TCODK_9
	Alt         = C.TCODK_ALT
	Apps        = C.TCODK_APPS
	Backspace   = C.TCODK_BACKSPACE
	CapsLock    = C.TCODK_CAPSLOCK
	Char        = C.TCODK_CHAR
	Control     = C.TCODK_CONTROL
	Delete      = C.TCODK_DELETE
	Down        = C.TCODK_DOWN
	End         = C.TCODK_END
	Enter       = C.TCODK_ENTER
	ESCAPE      = C.TCODK_ESCAPE
	Pressed     = C.TCOD_KEY_PRESSED
	Released    = C.TCOD_KEY_RELEASED
	F1          = C.TCODK_F1
	F10         = C.TCODK_F10
	F11         = C.TCODK_F11
	F12         = C.TCODK_F12
	F2          = C.TCODK_F2
	F3          = C.TCODK_F3
	F4          = C.TCODK_F4
	F5          = C.TCODK_F5
	F6          = C.TCODK_F6
	F7          = C.TCODK_F7
	F8          = C.TCODK_F8
	F9          = C.TCODK_F9
	Home        = C.TCODK_HOME
	Insert      = C.TCODK_INSERT
	KP0         = C.TCODK_KP0
	KP1         = C.TCODK_KP1
	KP2         = C.TCODK_KP2
	KP3         = C.TCODK_KP3
	KP4         = C.TCODK_KP4
	KP5         = C.TCODK_KP5
	KP6         = C.TCODK_KP6
	KP7         = C.TCODK_KP7
	KP8         = C.TCODK_KP8
	KP9         = C.TCODK_KP9
	KPAdd       = C.TCODK_KPADD
	KPDed       = C.TCODK_KPDEC
	KPDiv       = C.TCODK_KPDIV
	KPEnter     = C.TCODK_KPENTER
	KPMul       = C.TCODK_KPMUL
	KPSub       = C.TCODK_KPSUB
	Left        = C.TCODK_LEFT
	LWin        = C.TCODK_LWIN
	None        = C.TCODK_NONE
	NumLock     = C.TCODK_NUMLOCK
	PgDown      = C.TCODK_PAGEDOWN
	PgUp        = C.TCODK_PAGEUP
	Pause       = C.TCODK_PAUSE
	PrintScreen = C.TCODK_PRINTSCREEN
	Right       = C.TCODK_RIGHT
	RWin        = C.TCODK_RWIN
	ScrollLock  = C.TCODK_SCROLLLOCK
	Shift       = C.TCODK_SHIFT
	Space       = C.TCODK_SPACE
	Tab         = C.TCODK_TAB
	Up          = C.TCODK_UP
)
