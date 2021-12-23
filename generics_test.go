package chai

import (
	"errors"
	"fmt"
	"testing"
)

func err() error {
	return errors.New("error!")
}

func nerr() error {
	return nil
}

func isnil(err error) bool {
	return err == nil
}

type errzzcomp interface {
	comparable
	ZZ()
}

type errzz interface {
	ZZ()
}

type errzz2 interface {
	error
}

type errzz2comp interface {
	comparable
	error
}

type zzz struct {
}

func (zzz) ZZ() {
}

// func gen[T errzz](terr T) {
// 	spew.Dump(terr)
// 	fmt.Printf("terr: %+v\n", terr)
// 	fmt.Printf("terr Is: %+v\n", errors.Is(terr, nil))

// 	if terr == (error)(nil) && isnil(terr) {
// 		fmt.Printf("terr is nil!\n")
// 	} else {
// 		fmt.Printf("terr is NOT nil!\n")
// 	}
// }

func comparezz(z errzz) {
	if z != nil {
		fmt.Println("Not nil")
	} else {
		fmt.Println("nil!")
	}
}

func gen[T errzz](terr T) {
	// if terr == nil {

	// }

	comparezz(terr)
}

func TestZZ(t *testing.T) {
	fmt.Println("----err----")
	gen(&zzz{})
	var z *zzz
	gen(z)
	var zz errzz
	gen(zz)

	// gen(err())
	fmt.Println("----nerr----")
	// gen(nerr())
}
