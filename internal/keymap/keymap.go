package keymap

import (
	"github.com/micmonay/keybd_event"
)

var KeyMap map[string]int

func init() {
	KeyMap = make(map[string]int)
	KeyMap["ESC"] = keybd_event.VK_ESC
	KeyMap["1"] = keybd_event.VK_1
	KeyMap["2"] = keybd_event.VK_2
	KeyMap["3"] = keybd_event.VK_3
	KeyMap["4"] = keybd_event.VK_4
	KeyMap["5"] = keybd_event.VK_5
	KeyMap["6"] = keybd_event.VK_6
	KeyMap["7"] = keybd_event.VK_7
	KeyMap["8"] = keybd_event.VK_8
	KeyMap["9"] = keybd_event.VK_9
	KeyMap["0"] = keybd_event.VK_0
	KeyMap["A"] = keybd_event.VK_A
	KeyMap["B"] = keybd_event.VK_B
	KeyMap["C"] = keybd_event.VK_C
	KeyMap["D"] = keybd_event.VK_D
	KeyMap["E"] = keybd_event.VK_E
	KeyMap["F"] = keybd_event.VK_F
	KeyMap["G"] = keybd_event.VK_G
	KeyMap["H"] = keybd_event.VK_H
	KeyMap["I"] = keybd_event.VK_I
	KeyMap["J"] = keybd_event.VK_J
	KeyMap["K"] = keybd_event.VK_K
	KeyMap["L"] = keybd_event.VK_L
	KeyMap["M"] = keybd_event.VK_M
	KeyMap["N"] = keybd_event.VK_N
	KeyMap["O"] = keybd_event.VK_O
	KeyMap["P"] = keybd_event.VK_P
	KeyMap["Q"] = keybd_event.VK_Q
	KeyMap["R"] = keybd_event.VK_R
	KeyMap["S"] = keybd_event.VK_S
	KeyMap["T"] = keybd_event.VK_T
	KeyMap["U"] = keybd_event.VK_U
	KeyMap["V"] = keybd_event.VK_V
	KeyMap["W"] = keybd_event.VK_W
	KeyMap["X"] = keybd_event.VK_X
	KeyMap["Y"] = keybd_event.VK_Y
	KeyMap["Z"] = keybd_event.VK_Z
	KeyMap["F1"] = keybd_event.VK_F1
	KeyMap["F2"] = keybd_event.VK_F2
	KeyMap["F3"] = keybd_event.VK_F3
	KeyMap["F4"] = keybd_event.VK_F4
	KeyMap["F5"] = keybd_event.VK_F5
	KeyMap["F6"] = keybd_event.VK_F6
	KeyMap["F7"] = keybd_event.VK_F7
	KeyMap["F8"] = keybd_event.VK_F8
	KeyMap["F9"] = keybd_event.VK_F9
	KeyMap["F10"] = keybd_event.VK_F10
	KeyMap["F11"] = keybd_event.VK_F11
	KeyMap["F12"] = keybd_event.VK_F12

	KeyMap["NUMLOCK"] = keybd_event.VK_NUMLOCK
	KeyMap["SCROLLLOCK"] = keybd_event.VK_SCROLLLOCK
	KeyMap["CAPSLOCK"] = keybd_event.VK_CAPSLOCK
	KeyMap["PLUS"] = keybd_event.VK_KPPLUS
	KeyMap["MINUS"] = keybd_event.VK_MINUS
	KeyMap["EQUAL"] = keybd_event.VK_EQUAL
	KeyMap["BACKSPACE"] = keybd_event.VK_BACKSPACE
	KeyMap["TAB"] = keybd_event.VK_TAB
	KeyMap["LEFTBRACE"] = keybd_event.VK_LEFTBRACE
	KeyMap["RIGHTBRACE"] = keybd_event.VK_RIGHTBRACE
	KeyMap["ENTER"] = keybd_event.VK_ENTER
	KeyMap["SEMICOLON"] = keybd_event.VK_SEMICOLON
	KeyMap["APOSTROPHE"] = keybd_event.VK_APOSTROPHE
	KeyMap["GRAVE"] = keybd_event.VK_GRAVE
	KeyMap["BACKSLASH"] = keybd_event.VK_BACKSLASH
	KeyMap["COMMA"] = keybd_event.VK_COMMA
	KeyMap["DOT"] = keybd_event.VK_DOT
	KeyMap["SLASH"] = keybd_event.VK_SLASH
	KeyMap["ASTERISK"] = keybd_event.VK_KPASTERISK
	KeyMap["SPACE"] = keybd_event.VK_SPACE
	KeyMap["PAGEUP"] = keybd_event.VK_PAGEUP
	KeyMap["PAGEDOWN"] = keybd_event.VK_PAGEDOWN
	KeyMap["END"] = keybd_event.VK_END
	KeyMap["HOME"] = keybd_event.VK_HOME
	KeyMap["LEFT"] = keybd_event.VK_LEFT
	KeyMap["UP"] = keybd_event.VK_UP
	KeyMap["RIGHT"] = keybd_event.VK_RIGHT
	KeyMap["DOWN"] = keybd_event.VK_DOWN
	KeyMap["PRINT"] = keybd_event.VK_PRINT
	KeyMap["INSERT"] = keybd_event.VK_INSERT
	KeyMap["DELETE"] = keybd_event.VK_DELETE
}
