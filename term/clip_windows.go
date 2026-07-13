//go:build windows

package term

/*
#include <windows.h>

	void copy_to_clipboard_utf16(const wchar_t* text, int length) {
	    if (!OpenClipboard(NULL)) return;
	    EmptyClipboard();
	    // Размер в байтах: (length+1) * sizeof(wchar_t)
	    HGLOBAL hMem = GlobalAlloc(GMEM_MOVEABLE, (length+1) * sizeof(wchar_t));
	    if (hMem) {
	        wchar_t* pMem = (wchar_t*)GlobalLock(hMem);
	        for (int i = 0; i <= length; i++) {
	            pMem[i] = text[i];
	        }
	        GlobalUnlock(hMem);
	        SetClipboardData(CF_UNICODETEXT, hMem);
	    }
	    CloseClipboard();
	}
*/
import "C"
import (
	"unicode/utf16"
	"unsafe"
)

func CopyToClipboard(text string) {
	utf16 := utf16.Encode([]rune(text))
	utf16 = append(utf16, 0)
	C.copy_to_clipboard_utf16((*C.wchar_t)(unsafe.Pointer(&utf16[0])), C.int(len(utf16)-1))
}
