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
	"fmt"
	"github.com/raguay/goAlfred"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//
// Setup and constants that are used.
//
// MAXPROJECTS            This is the maximum number of projects allowed.
//
const (
	MAXPROJECTS int = 10
)

//
// Function:           main
//
// Description:       This is the main function for the TimeKeeper program. It taked the command line
//                            and parses it for the proper functionality.
//
func main() {
	if len(os.Args) > 1 {
		switch os.Args[1][0] {
		case 'm':
			ViewMonth()
		case 'w':
			ViewWeek()
		case 't':
			ViewDate()
		case 'r':
			RemoveProject()
		case 'c':
			ChangeProject()
		case 'a':
			AddProject()
		case 'o':
			StopStart()
		case 'p':
			project()
		case 's':
			fallthrough
		default:
			state()
		}
	}
}

//
// Function:           ViewMonth
//
// Description:       This function will calculate the time the current month for all the projects.
//
func ViewMonth() {

}

//
// Function:           ViewWeek
//
// Description:       This function will calculate the time the current week for all the projects.
//
// Inputs:
// 		variable 	description
//
func ViewWeek() {

}

//
// Function:           ViewDate
//
// Description:       This function will calculate the time for projects at a certain date.
//
func ViewDate() {
	vdate := GetCommandLineString()
}

//
// Function:           RemoveProject
//
// Description:       This function will remove a project from the list a valid projects.
//
func RemoveProject() {
	//
	// Get the project name from the command line.
	//
	proj := GetCommandLineString()

	//
	// Get the list of project names.
	//
	projects := GetListOfProjects()

	//
	// Open the projects file in truncation mode to remove all the old stuff.
	//
	Filename := goAlfred.Data() + "/projects.txt"
	Fh, err := os.OpenFile(Filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		//
		// The file would not open. Error out.
		//
		fmt.Print("Could not open the  projects file: ", Filename, "\n")
		os.Exit(1)
	}

	//
	// Loop through all the projects.
	//
	for i := 0; i < len(projects); i++ {
		if !strings.Contains(proj, projects[i]) {
			//
			// It is not the project to be removed. Put it into the file.
			//
			Fh.WriteString(projects[i] + "\n")
		}
	}

	//
	// Close the file.
	//
	Fh.Close()

	//
	// Tell the user that the project has been removed.
	//
	fmt.Print(proj + " has been removed!")
}

//
// Function:           ChangeProject
//
// Description:       This function will change the currently active project. If the old
//                            project was started, it will stop it first, then set the new project
//                            and start it.
//
func ChangeProject() {
	//
	// Get the project name from the command line.
	//
	proj := GetCommandLineString()

	//
	// Get the current project.
	//
	currentProject := GetCurrentProject()

	//
	// Stop the current project.
	//
	StopStartProject(currentProject, "stop")

	//
	// Save the new project to the data file.
	//
	SaveProject(proj)

	//
	// Start the new project.
	//
	StopStartProject(proj, "start")

	//
	// Tell the user it is started.
	//
	fmt.Print("The current project is now " + proj + " and is  started.")
}

//
// Function:           GetCommandLineString
//
// Description:       This function is used to get the after the function if there is one.
//                             If not, then just return nothing.
//
func GetCommandLineString() string {
	//
	// See if we have any input other then the command.
	//
	clstring := ""
	if len(os.Args) > 2 {
		clstring = strings.TrimSpace(os.Args[2])
	}

	//
	// Return the the string.
	//
	return (clstring)
}

//
// Function:           AddProject
//
// Description:       This function will add a new project to the list of current projects.
//
func AddProject() {
	//
	// Get the project name from the command line.
	//
	proj := GetCommandLineString()

	//
	// Create the file name that contains all the projects.
	//
	projectFile := goAlfred.Data() + "/projects.txt"
	Fh, err := os.OpenFile(projectFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		Fh, err = os.Create(projectFile)
		if err != nil {
			//
			// The file would not open. Error out.
			//
			fmt.Print("Could not open the projects file: ", projectFile, "\n")
			os.Exit(1)
		}
	}

	//
	// Write the new command with the time stamp to the buffer.
	//
	_, err = io.WriteString(Fh, proj+"\n")

	//
	// Lose the file.
	//
	Fh.Close()

	//
	// Tell the user that the project is added.
	//
	fmt.Print("Added project " + proj + " to the list.")
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
	buf, _ := ioutil.ReadFile(stateFile)
	curState := string(buf)

	//
	// Set the first command to the opposite of the current state. That way
	// the user simply pushes return to toggle states.
	//
	if strings.Contains(curState, "start") {
		goAlfred.AddResult("stop", "stop", "stop", "", "icon.png", "yes", "", "")
		goAlfred.AddResult("start", "start", "start", "", "icon.png", "yes", "", "")
	} else {
		goAlfred.AddResult("start", "start", "start", "", "icon.png", "yes", "", "")
		goAlfred.AddResult("stop", "stop", "stop", "", "icon.png", "yes", "", "")
	}

	//
	// Print out the xml string.
	//
	fmt.Print(goAlfred.ToXML())
}

//
// Function:           project
//
// Description:       This function creates a list of the projects and displays the ones
//                            similar to the input.
//
func project() {
	//
	// Get the project name from the command line.
	//
	proj := GetCommandLineString()

	//
	// Set our default string.
	//
	goAlfred.SetDefaultString("Alfred Time Keeper:  Sorry, no match...")

	//
	// Get the list of projects.
	//
	projects := make([]string, MAXPROJECTS)
	projects = GetListOfProjects()

	//
	// The regexp split statement gives one string more than was split out. The last
	// string is a catchall. It does not need to be included.
	//
	numproj := len(projects) - 1

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
}

//
// Function:           GetListOfProjects
//
// Description:       This function will return an array of string with the names of the project.
//
func GetListOfProjects() []string {
	//
	// Create the projects array and populate it.
	//
	projectFile := goAlfred.Data() + "/projects.txt"
	buf, _ := ioutil.ReadFile(projectFile)

	//
	// Split out the different project names into separate strings.
	//
	return (regexp.MustCompile("\n|\r").Split(string(buf), -1))
}

//
// Function:           StopStart
//
// Description:       This will place a start or stop time stamp for the current project and
//                            current date.
//
func StopStart() {
	//
	// See if we have any input other then the command.  If not, assume a stop command.
	//
	cmd := "stop"
	if len(os.Args) > 2 {
		cmd = strings.ToLower(os.Args[2])
	}

	//
	// Get the current project.
	//
	currentProject := GetCurrentProject()

	//
	// Run the appropriate function and print the results.
	//
	fmt.Print(StopStartProject(currentProject, cmd))
}

//
// Function:           GetCurrentProject
//
// Description:       This function will retrieve the current project from the
//                            state file.
//
func GetCurrentProject() string {
	//
	// Get the current project.
	//
	Filename := goAlfred.Data() + "/project.txt"
	buf, _ := ioutil.ReadFile(Filename)

	//
	// Convert the current project to a string, trim it, and return it.
	//
	return (strings.TrimSpace(string(buf)))
}

//
// Function:           SaveProject
//
// Description:       This function will save the given project name to the
//                            current project file.
//
// Inputs:
// 		proj 	     Name of the new project
//
func SaveProject(proj string) {
	//
	// Write the new project.
	//
	Filename := goAlfred.Data() + "/project.txt"
	err := ioutil.WriteFile(Filename, []byte(proj), 0666)
	if err != nil {
		fmt.Print("Can not write the project file: " + Filename)
		os.Exit(1)
	}
}

//
// Function:           StopStartProject
//
// Description:       This function is used to set the state for the given project.
//
// Inputs:
// 		currentProject 	The project to effect the state of.
//               cmd                    The start or stop command.
//
func StopStartProject(currentProject string, cmd string) string {
	//
	// Setup the result string.
	//
	resultStr := ""

	//
	// Get the current state.
	//
	Filename := goAlfred.Data() + "/laststate.txt"
	buf, err := ioutil.ReadFile(Filename)
	currentState := "stop"
	if err == nil {
		//
		// Convert the current project to a string and trim it.
		//
		currentState = strings.TrimSpace(string(buf))
	}

	//
	// Is the current state the same as the new state?
	//
	if strings.Contains(cmd, currentState) {
		//
		// It is already in that state. Do nothing, but give a message.
		//
		resultStr = "Already " + cmd + "\n"
	} else {
		//
		// Okay, we can proceed with writing the new state into the
		// dated project file. Open the file for writing.
		//
		currentTime := time.Now()
		Filename = generateTimeLogFileName(currentProject, currentTime)
		Fh, err := os.OpenFile(Filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			//
			// The file would not open. Error out.
			//
			fmt.Print("Could not open the dated project file: ", Filename, "\n")
			os.Exit(1)
		}

		//
		// Write the new command with the time stamp to the buffer.
		//
		str := strconv.FormatInt(currentTime.Unix(), 10) + ":" + cmd + "\n"
		_, err = io.WriteString(Fh, str)

		//
		// Lose the file.
		//
		Fh.Close()

		//
		// Write the laststate file with the new state.
		//
		ioutil.WriteFile(goAlfred.Data()+"/laststate.txt", []byte(cmd), 0666)

		//
		// Tell the user it is set.
		//
		resultStr = currentProject + " is now " + cmd
	}

	//
	// Return the resulting string.
	//
	return (resultStr)
}

//
// Function:           generateTimeLogFileName
//
// Description:       This functions creates the time log file based on the project name and
//                            date.
//
// Inputs:
// 		proj 	     Name of the project
//               dt           Date in question
//
func generateTimeLogFileName(proj string, dt time.Time) string {
	//
	// Generate the proper file name based on the project name and date.
	//
	filename := goAlfred.Data() + "/" + proj + "_" + dt.Format("2006-01-02") + ".txt"
	return (filename)
}
