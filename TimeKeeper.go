package main

//
// Program:              TimeKeeper.go
//
// Description:         This program runs the feedback logic for selecting  the on/off state for
//                              the currently timed project. 
//

//
// Import the libraries we use for this program.
//
import (
	"github.com/raguay/goAlfred"
	"os"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	//
	// Get the last state of the current project. 
	//
	stateFile := goAlfred.Data() + "/laststate.txt"
	buf, err := ioutil.ReadFile(stateFile)
	if(err != nil) {
		fmt.Print("Can not read the state file!")
		os.Exit(1)
	}

	//
	// Split out the different project names into separate strings. 
	//
	curState := string(buf)
	if(strings.Contains(curState ,"start")) {
		goAlfred.AddResult( "stop", "stop", "stop", "", "icon.png", "yes", "", "") 		
		goAlfred.AddResult( "start", "start", "start",  "", "icon.png", "yes", "", "") 		
	} else {
		goAlfred.AddResult( "start", "start", "start",  "", "icon.png", "yes", "", "") 		
		goAlfred.AddResult( "stop", "stop", "stop", "", "icon.png", "yes", "", "") 		
	}

	//
	// Print out the xml string. 
	//
	fmt.Print(goAlfred.ToXML())
}