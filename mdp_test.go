package rltetris

import "testing"

func TestOneHotAppend(t *testing.T) {
	arr := make([]byte, 0, 100)
	arr = appendOneHot(arr, 0, 7, 0)
	if len(arr) != 8 {
		t.Errorf("No se produjo un arreglo del tamaño correcto, esperaba %v, obtuve %v", 8, len(arr))
	}

	arr = make([]byte, 0, 100)
	arr = appendOneHot(arr, 0, 7, 7)
	if len(arr) != 8 {
		t.Errorf("No se produjo un arreglo del tamaño correcto, esperaba %v, obtuve %v", 8, len(arr))
	}
}
