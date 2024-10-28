package main

import (
	"fmt"
	"strings"
	"time"
)

var input string = ""                     //Input String of the std_in
var done chan bool = make(chan bool)      //Main waits for the goroutine
var isMain chan bool = make(chan bool)
var startDone chan bool = make(chan bool) //finish the Start ticker goroutine
var isStart chan bool = make(chan bool)   //Needa look how to make this better where it is involved
var isOption chan bool = make(chan bool)  //Needa look how to make this better where it is involved

var formattedCurTime time.Duration 
//Var which will be shown in TUI and used in Options
//NEEDS TO BE CHANGED -> HAVE TO BE IN A CONF FILE!!
var workTime time.Duration = time.Duration(8 * time.Hour)
var breakTime time.Duration = time.Duration(30 * time.Minute)
var refreshRate time.Duration = time.Duration(5 * time.Second)
var endOfWork bool = false;
var workedTimeText string = "0 Seconds";

func main() {
	//start up text and clear after programm ends
	defer clearText()

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

        case input == "quit", input == "q":
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
            workedTimeHours := workedTime.Hours();
            workedTimeMinutes := workedTime.Minutes();
            workedTimeSeconds := workedTime.Seconds();
            workedTimeText = fmt.Sprint(workedTimeHours, "Hours", workedTimeMinutes, "Minutes", workedTimeSeconds, "Seconds")
            //TODO: fix printing with comma -> rounding but to 5 not if >= 5 then round to 5...
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

//changes rather Timer should also display when u can leave work
func toggleEndOfWork(){

    if(endOfWork == true){
        endOfWork = false
    }else{
        endOfWork = true
    }

}

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
    fmt.Println("\tP(ause) traching")
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

func clearText(){
		fmt.Print("\033[H\033[2J")
}
