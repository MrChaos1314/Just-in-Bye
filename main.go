package main

import (
	"fmt"
	"strings"
	"time"
    "encoding/json"
    "os"
)

var input string = ""                     //Input String of the std_in
var done chan bool = make(chan bool)      //Main waits for the goroutine
var isMain chan bool = make(chan bool)
var startDone chan bool = make(chan bool) //finish the Start ticker goroutine
var isStart chan bool = make(chan bool)   //Needa look how to make this better where it is involved
var isOption chan bool = make(chan bool)  //Needa look how to make this better where it is involved

var formattedCurTime time.Duration 
//Var which will be shown in TUI and used in Options
//Default value ... else it gets loaded in the beginning from a conf
var workTime time.Duration = time.Duration(8 * time.Hour)
var breakTime time.Duration = time.Duration(30 * time.Minute)
var refreshRate time.Duration = time.Duration(5 * time.Second)
var endOfWork bool = false;
var workedTimeText string = "0s";

//Config Json-Struct 
type conf struct{

    WorkTime time.Duration
    BreakTime time.Duration

    RefreshRate time.Duration

    EndOfWork bool

};

func main() {

	//start up text and clear after programm ends
	defer clearText()

    loadConf()

	//checks for the std_input
    clearText()
    checkInputMain()

    //if user is tracking time check its inputs
loop:
    for {
        select {
        case <-isStart:
            clearText()
            checkInputStart()
            break
        case <-isOption:
            clearText()
            checkInputOption()
            break
        case <-isMain:
            clearText()
            checkInputMain()
            break
        case <-done:
            fmt.Print("test")
            break loop 
		}
	}
}

/*-----------------Input Checks------------------*/

func checkInputMain() {
	showTextMain()
	//label for breaking mid switch statement
loop:
	for {
		//get std_input and make it lowercase
		fmt.Scan(&input)
		input = strings.ToLower(input)
		//decide through the input what the user can do
		switch {

		//Starting timer in thread
		case input == "start", input == "s":
			go StartTimer()
			break loop

		//Change options of the timer
		case input == "options", input == "o":
			go Options()
			break loop

        //TODO: Add further option like get overtime and current time worked and...

        case input == "quit", input == "q":
            go Quit()
            break loop

		//Wrong Input
        //TODO: Currently not working properly - where should the TUI be reseted?
		default:
			showTextMain()
			fmt.Println("WTF is going on!!")
			fmt.Println("Can't you read???")
		}
	}
}

// Checks the input if the Timer was started
func checkInputStart() {
	showTextStart()
loop:
	for {
		//get user Input and LowerCase it
		fmt.Scan(&input)
		input = strings.ToLower(input)
		//decide through the input what the user can do
		switch {
		//Stops the timer
		case input == "stop", input == "s":
			go StopTimer()
			break loop
		//Wrong input
		default:
            clearText()
			fmt.Println("WTF is going on!!")
			fmt.Println("Can't you read???")
		}
	}
}

// Checks the input if the Timer was started
func checkInputOption() {
loop:
	for {
        clearText()
        showTextOptions()
		//get user Input and LowerCase it
		fmt.Scan(&input)
		input = strings.ToLower(input)
		//decide through the input what the user can do
		switch {

		//Stops the timer
		case input == "end", input == "e":
            toggleEndOfWork()
			break 

		//Change the WorkTime
		case input == "work", input == "w":
            changeWorkTime()
			break 

		//change the BreakTime
		case input == "break", input == "b":
            toggleEndOfWork()
			break 

        case input == "refresh", input == "r":
            toggleEndOfWork()
            break

        case input == "quit", input == "q":
            go BackToMain()
            break loop

        case input == "save", input == "s":
            writeConf()
            go BackToMain()
            break loop
		//Wrong input
		default:
            break
		}
	}
}

/*-----------------Functions------------------*/

// Start the timer
func StartTimer() {
	isStart <- true                                 //switch TUI to startScreen
    showTime := time.NewTicker(refreshRate)     //Setups Ticker of X Seconds - TODO: Make it configureable
	startTime := time.Now()                         //get the time work started

	//Every time Ticker -> Prints current working time
	for {
		select {
		//finish goroutine
		case <-startDone:
			return

			//Print the current working time
		case curTime := <-showTime.C:
            clearText()
            showTextStart()
            workedTime := curTime.Sub(startTime)
            workedTimeText = fmt.Sprint(workedTime.Round(time.Second).String())
		}
	}
}

// Stop the timer
func StopTimer() {
	isMain <- true
	startDone <- true
}

//change TUI in main to Options
func Options() {
    isOption <- true
}

//change TUI back to Main
func BackToMain(){
    isMain <- true
}

//Quit whole Program
func Quit(){
    done <- true
}


/****************Option Functions****************************/

//changes rather Timer should also display when u can leave work
func toggleEndOfWork(){

    if(endOfWork == true){
        endOfWork = false
    }else{
        endOfWork = true
    }

}

func changeWorkTime(){
    clearText()
    showTextChangeWorkTime()
loop:
	for {
		//get user Input and LowerCase it
		fmt.Scan(&input)
		input = strings.ToLower(input)
		//decide through the input what the user can do
		switch {

		//Stops the timer
		case input == "abort", input == "a":
		    break loop

        case input == "quit", input == "q":
            break loop

		//Wrong input
		default:
            _, err := fmt.Scanf("%d", &input)
            if(err != nil){
                clearText()
                fmt.Println("WTF is going on!!")
                fmt.Println("Can't you read???")
            }
            //CONTINUE HERE - input before 5 lines needs to be changed maybe
		}
	}
}



func writeConf(){
    fd, err := os.Create("./config/options.conf")
    defer fd.Close();
    if(err != nil){
        panic(err);
    }

    config := &conf{
        
        WorkTime: workTime,
        BreakTime: breakTime,

        RefreshRate: refreshRate, 

        EndOfWork: endOfWork,
    }

    config_json, _ := json.Marshal(config);
    written, err := fd.Write(config_json)
    if(err != nil){
        panic(err)
    }
    fmt.Printf("Bytes written: %d\n", written)
}

func readConf() (time.Duration, time.Duration, time.Duration, bool){
    config := &conf{}
    fd, err := os.ReadFile("./config/options.conf")
    if(err != nil){
        panic(err);
    }
    json.Unmarshal(fd, config)

    return config.WorkTime, config.BreakTime, config.RefreshRate, config.EndOfWork
}

func loadConf(){
    workTime, breakTime, refreshRate, endOfWork = readConf()
}

/******************************************************************/


/*-----------------TUI-Text------------------*/


func showTextMain(){
    fmt.Println("U wanna Track your time?")
    fmt.Println("Here are your options:")
    fmt.Println()
    fmt.Println("\tS(tart) tracking")
    fmt.Println("\tO(ptions) for the Tracker")
    fmt.Println()
    fmt.Println("Q(uit)")
    fmt.Println()
}

func showTextStart(){
    fmt.Println("Your time gets currently tracked...")
    fmt.Println("Your current options:")
    fmt.Println()
    fmt.Println("\tS(top) tracking")
    fmt.Println("\tP(ause) tracking")
    fmt.Println()
    fmt.Println("Your current work time is: ", workedTimeText)
    fmt.Println()
}

func showTextOptions(){
    fmt.Println("Options:")
    fmt.Println()
    fmt.Println("\tW(ork) time in Hours: ", workTime.Hours(), "Hours")
    fmt.Println("\tB(reak) time in Minutes: ", breakTime.Minutes(), "Minutes")
    fmt.Println("\tR(efresh) rate in Seconds: ", refreshRate.Seconds(), "Seconds")
    fmt.Println("\tE(nd) of Work: ", endOfWork)
    fmt.Println("Q(uit)")
    fmt.Println("")
}

func showTextChangeWorkTime(){
    fmt.Println("Change the value for the work time:")
    fmt.Println()
    fmt.Println("\tCurrent value for work time in hours: ", workTime.Minutes())
    fmt.Println()
    fmt.Println("(A)bort - (Q)uit")
    fmt.Println()
    fmt.Print("Change value for work time in hours: ")
}

func clearText(){
		fmt.Print("\033[H\033[2J")
}
