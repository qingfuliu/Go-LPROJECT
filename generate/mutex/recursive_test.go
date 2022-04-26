<<<<<<< HEAD
package mutex

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestNewRecursiveMutex(t *testing.T) {
	sli := make([]int,0, 10)
	p := unsafe.Pointer(&sli)

	sliSize := unsafe.Sizeof(sli)

	intSize := unsafe.Sizeof(int(0))

	pointerSize := unsafe.Sizeof(p)

	cap_:=(*int64)(unsafe.Pointer((uintptr(pointerSize)+uintptr(p))))
	size_:=(*int64)(unsafe.Pointer((uintptr(pointerSize*2)+uintptr(p))))
	fmt.Println(sliSize,intSize,pointerSize,*cap_,*size_)
}
=======
package mutex

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestNewRecursiveMutex(t *testing.T) {
	sli := make([]int,0, 10)
	p := unsafe.Pointer(&sli)

	sliSize := unsafe.Sizeof(sli)

	intSize := unsafe.Sizeof(int(0))

	pointerSize := unsafe.Sizeof(p)

	cap_:=(*int64)(unsafe.Pointer((uintptr(pointerSize)+uintptr(p))))
	size_:=(*int64)(unsafe.Pointer((uintptr(pointerSize*2)+uintptr(p))))
	fmt.Println(sliSize,intSize,pointerSize,*cap_,*size_)
}
>>>>>>> master
