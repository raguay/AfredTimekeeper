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
// TSDir                         This keeps the directory name for the time sheets. It is a complete path.
//
const (
	MAXPROJECTS int = 20
)

var TSDir = ""

//
// Function:           main
//
// Description:       This is the main function for the TimeKeeper program. It takes the command line
//                          and parses it for the proper functionality.
//
func main() {
	if len(os.Args) > 1 {
		switch os.Args[1][0] {
		case 'm':
			//
			// atk:month
			//
			SystemViewMonth()
		case 'l':
			//
			// atk:lastmonth
			//
			SystemViewLastMonth()
		case 'w':
			//
			// atk:week
			//
			SystemViewWeek()
		case 't':
			//
			// atk:current
			//
			SystemViewDate()
		case 'r':
			//
			// atk:remove
			//
			RemoveProject()
		case 'c':
			//
			// atk:project
			//
			ChangeProject()
		case 'b':
			//
			// atk:current
			//
			SystemViewDay()
		case 'a':
			//
			// atk:addproject
			//
			AddProject()
		case 'o':
			//
			// atk:state
			//
			StopStart()
		case 'p':
			//
			// Used for atk:project script filter
			//
			project()
		case 'T':
			//
			// atk:time
			//
			SystemAllProjects()
		case 's':
			fallthrough
		default:
			//
			// Used for the script filter on atk:state
			//
			state()
		}
	}
}

//
// Function:		getTimeSheetDir
//
// Description:		This function is used to cache a copy of the time
//					sheet directory and give it in the return.
//
func getTimeSheetDir() string {
	if strings.Contains("", TSDir) {
		Filename := goAlfred.Data() + "/dir.txt"
		buf, err := ioutil.ReadFile(Filename)
		if err == nil {
			//
			// Convert the directory path to a string and trim it.
			//
			TSDir = strings.TrimSpace(string(buf))
		} else {
			//
			// The Directory has not been set. Make it the default location.
			//
			TSDir = goAlfred.Data()
			err = ioutil.WriteFile(Filename, []byte(TSDir), 0666)
		}
	}

	//
	// Return the directory to the time sheets.
	//
	return (TSDir)
}

//
// Function:              SystemAllProjects
//
// Description:          This function will display to the terminal the time for all projects on
//                             the day given on the command line next.
//
func SystemAllProjects() {
	//
	// Get the current date in case there isn't one on the command line.
	//
	tm := time.Now()
	if len(os.Args) > 2 {
		if strings.Contains("today", os.Args[2]) {
			//
			// Today's date.
			//
			tm = time.Now()
		} else if strings.Contains("yesterday", os.Args[2]) {
			//
			// Yesterday is today minus one day.
			//
			tm = time.Now()
			tm = tm.AddDate(0, 0, -1)
		} else {
			//
			// Parse the date string given.
			//
			tm, _ = time.Parse("2006-Jan-02", os.Args[2])
		}
	}

	//
	// Get the list of project names.
	//
	proj := GetListOfProjects()

	//
	// For each project, get the time spent on it for the given day.
	//
	numproj := len(proj) - 1
	for i := 0; i < numproj; i++ {
		fmt.Printf("%s: %s\n", proj[i], formatTimeString(GetTimeAtDate(proj[i], tm)))
	}
}

//
// Function:           SystemViewMonth
//
// Description:       This function will calculate the time the current month for all the projects.
//
func SystemViewMonth() {
	//
	// Get the current project.
	//
	currentProject := GetCurrentProject()

	//
	// Get the time on that project for this month. The current time gives the current month.
	//
	now := time.Now()
	year, month, _ := now.Date()
	tm := GetTimeAtMonth(currentProject, year, month)

	//
	// format the time string and print it out.
	//
	fmt.Print(formatTimeString(tm))
}

//
// Function:           SystemViewLastMonth
//
// Description:       This function will calculate the time last
//							 month for the current project.
//
func SystemViewLastMonth() {
	//
	// Get the current project.
	//
	currentProject := GetCurrentProject()

	//
	// Get the time on that project for this month. The current time gives the current month.
	//
	now := time.Now()
	year, month, _ := now.Date()
	month = month - 1

	tm := GetTimeAtMonth(currentProject, year, month)

	//
	// format the time string and print it out.
	//
	fmt.Print(formatTimeString(tm))
}

//
// Function:           GetTimeAtMonth
//
// Description:       This function will take a project and calculate
//							 the time spent on that project for a particular
//							 month.
//
func GetTimeAtMonth(project string, year int, month time.Month) int64 {
	tm := int64(0)
	dateStart := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	dateEnd := time.Date(year, month, 32, 0, 0, 0, 0, time.UTC)
	lastDay := 32 - dateEnd.Day()

	//
	// Get the time added up for the whole month.
	//
	for i := 0; i < lastDay; i++ {
		tm += GetTimeAtDate(project, dateStart.AddDate(0, 0, i))
	}

	//
	// Return the amount of time calculated.
	//
	return (tm)
}

//
// Function:           SystemViewWeek
//
// Description:       This function will calculate the time the current week for all the projects.
//
// Inputs:
// 		variable 	description
//
func SystemViewWeek() {
	currentProject := GetCurrentProject()
	tm := GetTimeAtWeek(currentProject, time.Now())
	fmt.Print(formatTimeString(tm))
}

//
// Function:           GetTimeAtDate
//
// Description:       This function will take a project and calculate the time spent
//                          on that project for a particular date.
//
func GetTimeAtWeek(project string, date time.Time) int64 {
	tm := int64(0)
	dateStart := date
	dateEnd := date
	switch date.Weekday() {
	case 0:
		{
			dateEnd = dateEnd.AddDate(0, 0, 6)
		}
	case 1:
		{
			dateStart = dateStart.AddDate(0, 0, -1)
			dateEnd = dateEnd.AddDate(0, 0, 5)
		}
	case 2:
		{
			dateStart = dateStart.AddDate(0, 0, -2)
			dateEnd = dateEnd.AddDate(0, 0, 4)
		}
	case 3:
		{
			dateStart = dateStart.AddDate(0, 0, -3)
			dateEnd = dateEnd.AddDate(0, 0, 3)
		}
	case 4:
		{
			dateStart = dateStart.AddDate(0, 0, -4)
			dateEnd = dateEnd.AddDate(0, 0, 2)
		}
	case 5:
		{
			dateStart = dateStart.AddDate(0, 0, -5)
			dateEnd = dateEnd.AddDate(0, 0, 1)
		}
	case 6:
		{
			dateStart = dateStart.AddDate(0, 0, -6)
		}
	}
	//
	// Get the time added up for th whole week.
	//
	for i := 0; i < 7; i++ {
		tm += GetTimeAtDate(project, dateStart.AddDate(0, 0, i))
	}
	return (tm)
}

//
// Function:           SystemViewDate
//
// Description:       This function will calculate the time for projects at a certain date.
//
func SystemViewDate() {
	currentProject := GetCurrentProject()
	tm := GetTimeAtDate(currentProject, time.Now())
	fmt.Print(formatTimeString(tm))
}

//
// function:            SystemViewDay
//
// Description:        This function is for displaying a nice time for the current project.
//
func SystemViewDay() {
	currentProject := GetCurrentProject()
	tm := GetTimeAtDate(currentProject, time.Now())
	ctime := formatTimeString(tm)
	state := GetCurrentState()
	fmt.Printf("The current time on %s is %s. Current state is %s.", currentProject, ctime, state)
}

//
// Function:           GetTimeAtDate
//
// Description:       This function will take a project and calculate the time spent
//                          on that project for a particular date.
//
func GetTimeAtDate(project string, date time.Time) int64 {
	//
	// Get the current project.
	//
	filename := generateTimeLogFileName(project, date)
	tm := readDayTime(filename)
	return tm
}

//
// Function:             formatTimeString
//
// Description:         This function takes the number of seconds and returns a string
//                            in hour:minute:seconds format with zero padding.
//
// Input:
//                            tm          time in seconds (an int64)
//
func formatTimeString(tm int64) string {
	min := int(tm / 60)
	sec := tm - int64(min*60)
	hr := min / 60
	min = min - (hr * 60)
	return fmt.Sprintf("%02d:%02d:%02d", hr, min, sec)
}

//
// Function:             readDayTime
//
// Description:         This function reads a time sheet file and calculates the time
//                            represented in that file.

func readDayTime(filename string) int64 {
	buf, _ := ioutil.ReadFile(filename)
	times := regexp.MustCompile("\n|\r").Split(string(buf), -1)

	//
	// Loop through all the time lines.
	//
	tmwork := int64(0)
	firsttime := int64(0)
	first := false
	for i := 0; i < len(times); i++ {
		if !strings.Contains("", times[i]) {
			//
			// Split by colon to time and action.
			//
			parts := strings.Split(times[i], ":")
			if strings.Contains("start", parts[1]) {
				firsttime, _ = strconv.ParseInt(parts[0], 10, 64)
				first = true
			} else {
				tm, _ := strconv.ParseInt(parts[0], 10, 64)
				tmwork += tm - firsttime
				first = false
			}
		}
	}

	//
	// If a start was the last thing processed, that means it is still being timed. Get the
	// current time to see the overall time. firsttime is the time stamp when the start
	// was given.
	//
	if first {
		currentTime := time.Now()
		ctime := currentTime.Unix()
		tmwork += ctime - firsttime
	}
	//
	// Return the final Time.
	//
	return tmwork
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
	Filename := getTimeSheetDir() + "/projects.txt"
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
	// Change to the specified project.
	//
	ChangeNamedProject(proj)
}

func ChangeNamedProject(proj string) {
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
	projectFile := getTimeSheetDir() + "/projects.txt"
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
	// Make it the new current project.
	//
	ChangeNamedProject(proj)

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
	stateFile := getTimeSheetDir() + "/laststate.txt"
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
	// Get the latest project.
	//
	latestproject := GetCurrentProject()

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
	// For each project, create a result line. Show all put the current project.
	//
	for i := 0; i < numproj; i++ {
		if !strings.Contains(projects[i], latestproject) || (latestproject == "") {
			goAlfred.AddResultsSimilar(proj, projects[i], projects[i], projects[i], "", "icon.png", "yes", "", "")
		}
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
	projectFile := getTimeSheetDir() + "/projects.txt"
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
	Filename := getTimeSheetDir() + "/project.txt"
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
	Filename := getTimeSheetDir() + "/project.txt"
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

	currentState := GetCurrentState()

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
		Filename := generateTimeLogFileName(currentProject, currentTime)
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
		str := fmt.Sprintf("%d:%s\n", currentTime.Unix(), cmd)
		_, err = io.WriteString(Fh, str)

		//
		// Lose the file.
		//
		Fh.Close()

		//
		// Write the laststate file with the new state.
		//
		ioutil.WriteFile(getTimeSheetDir()+"/laststate.txt", []byte(cmd), 0666)

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
// function:                GetCurrentState
//
// Description:           This function gets the current state of the project.
//
func GetCurrentState() string {
	//
	// Get the current state.
	//
	Filename := getTimeSheetDir() + "/laststate.txt"
	buf, err := ioutil.ReadFile(Filename)
	currentState := "stop"
	if err == nil {
		//
		// Convert the current project to a string and trim it.
		//
		currentState = strings.TrimSpace(string(buf))
	} else {
		//
		// The current state is stop if the file does not exist. Create the file
		// with that state.
		//
		err = ioutil.WriteFile(Filename, []byte(currentState), 0666)
	}
	return currentState
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
	filename := getTimeSheetDir() + "/" + proj + "_" + dt.Format("2006-01-02") + ".txt"
	return (filename)
}
