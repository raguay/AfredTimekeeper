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
	"regexp"
)

//
// Setup and constants that are used.
// 
// MAXPROJECTS            This is the maximum number of projects allowed. 
//
const (
	MAXPROJECTS int =10
)

//
// Function:           main 
//
// Description:       This is the main function for the TimeKeeper program. It taked the command line
//                            and parses it for the proper functionality. 
//
func main() {
	if(len(os.Args) > 1) {
		switch os.Args[1][0] {
			case 'p':  project()
			case 's': fallthrough
			default:  state()
		}
	}
}

//
// Function:           state 
//
// Description:       This function gives the proper output for changing the state. The state
//                            first is the one opposite from the current state. 
//
func state() {
	//
	// Get the last state of the current project. 
	//
	stateFile := goAlfred.Data() + "/laststate.txt"
	buf, err := ioutil.ReadFile(stateFile)
	if(err != nil) {
		fmt.Print("Can not read the state file!")
		os.Exit(1)
	}
	curState := string(buf)

        //
        // Set the first command to the opposite of the current state. That way
        // the user simply pushes return to toggle states.  
        //
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
	fmt.Print("\n")
}

//
// Function:           project 
//
// Description:       This function creates a list of the projects and displays the ones
//                            similar to the input. 
//
func project() {
	//
	// See if we have any input other then the command. 
	//
	proj := ""
	if(len(os.Args) > 2) {
		proj = strings.ToLower(os.Args[2])
	}

	//
	// Set our default string. 
	//
	goAlfred.SetDefaultString("Alfred Time Keeper:  Sorry, no match...")

	//
	// Create the projects array and populate it. 
	//
	projects := make([]string,MAXPROJECTS)
	projectFile := goAlfred.Data() + "/projects.txt"
	buf, err := ioutil.ReadFile(projectFile)
	if(err != nil) {
		fmt.Print("Can not read the projects file!")
		os.Exit(1)
	}

         //
         // Split out the different project names into separate strings. 
         //
	projects = regexp.MustCompile("\n|\r").Split(string(buf), -1)

         //
         // The regexp split statement gives one string more than was split out. The last 
        // string is a catchall. It does not need to be included. 
        //
         numproj := len(projects) -1

        //
         // For each project, create a result line. 
         //
         for i := 0; i < numproj; i++ {
         	goAlfred.AddResultsSimilar(proj, projects[i], projects[i], projects[i], "", "icon.png", "yes", "", "") 
         }

	//
	// Print out the xml string. 
	//
	fmt.Print(goAlfred.ToXML())
	fmt.Print("\n")
}