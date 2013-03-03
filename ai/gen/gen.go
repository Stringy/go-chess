//Package gen is responsible for generating a list of moves for a particular game state.
//It uses bitboards to contain the information regarding each piece type and a single 32 bit integer
//to contain movement information. It is only designed to generate from a single state, to a depth of one.
//Any further than this must be extended outside of this package. 
package gen

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	logger = log.New(ioutil.Discard, "", 0)
)

// init Initialises all global data including attack masks
func init() {
	InitAllData()        //init globals
	InitialiseAllMasks() //init attack masks
	file, err := os.Open("gen.log")
	if err != nil {
		//do nothing
	} else {
		SetLogger(file)
	}
}

func SetLogger(w io.Writer) {
	logger = log.New(w, "gen", log.LstdFlags|log.Lmicroseconds)
}
