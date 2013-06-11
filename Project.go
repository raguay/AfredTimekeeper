package main

//
// Program:              Project.go
//
// Description:         This program runs the feedback logic for selecting the current project for
//                              Alfred Time Keeper. 
//

//
// Import the libraries we use for this program.
//
import (
	"github.com/raguay/goAlfred"
	"os"
	"fmt"
	"io/ioutil"
	"regexp"
)

const (
	MAXPROJECTS int =10
)

func main() {
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
		goAlfred.AddResult( projects[i], projects[i], projects[i], "", "icon.png", "yes", "", "") 
	}

	//
	// Print out the xml string. 
	//
	fmt.Print(goAlfred.ToXML())
}