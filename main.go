package main

import (
	"fmt"
	"strings"
	"time"
)

var input string = ""                     //Input String of the std_in
var done chan bool = make(chan bool)      //Main waits for the goroutine
var startDone chan bool = make(chan bool) //finish the Start ticker goroutine
var isStart chan bool = make(chan bool)   //Needa look how to make this better where it is involved
var isOption chan bool = make(chan bool)  //Needa look how to make this better where it is involved

func main() {
	//start up text and clear after programm ends
	defer showText("clear")

	//checks for the std_input
	checkInputMain()

	//if user is tracking time check its inputs
	//loop:
	for {
		select {
		case <-isStart:
			checkInputStart()
		case <-isOption:
			checkInputOption()
		case <-done:
			checkInputMain()
		}
	}
}

/*-----------------Input Checks------------------*/

func checkInputMain() {
	showText("main")
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
		//Add further option like get overtime and current time worked and...

		//Wrong Input
		default:
			showText("clear")
			fmt.Println("WTF is going on!!")
			fmt.Println("Can't you read???")
		}
	}
}

// Checks the input if the Timer was started
func checkInputStart() {
	showText("start")
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
			showText("clear")
			fmt.Println("WTF is going on!!")
			fmt.Println("Can't you read???")
		}
	}
}

// Checks the input if the Timer was started
func checkInputOption() {
	showText("option")
loop:
	for {
		//get user Input and LowerCase it
		fmt.Scan(&input)
		input = strings.ToLower(input)
		//decide through the input what the user can do
		switch {
		//Stops the timer
		case input == "", input == "":
			StopTimer()
			break loop
		//Wrong input
		default:
			showText("clear")
			fmt.Println("WTF is going on!!")
			fmt.Println("Can't you read???")
		}
	}
}

/*-----------------Functions------------------*/

// Start the timer
func StartTimer() {
	startShowTime := time.NewTicker(5 * time.Second) //Setups Ticker of X Seconds
	isStart <- true                                  //tell goroutine is running
	startTime := time.Now()                          //get the time work started
	//fmt.Println("Tracks time...")                                 //Debug
	//Every time Ticker -> Prints current working time
	for {
		select {
		//finish goroutine
		case <-startDone:
			return

			//Print the current working time
		case curTime := <-startShowTime.C:
			showText("start")
			fmt.Print("Your current Work time is: ", curTime.Sub(startTime))
		}
	}
}

// Stop the timer
func StopTimer() {
	done <- true
	startDone <- true
}

func Options() {

}

/*-----------------TUI-Text------------------*/

// Shows text to the cmd according to the label
func showText(label string) {
	label = strings.ToLower(label)
	switch {
	case label == "main":
		showText("clear")
		fmt.Println("U wanna Track your time?")
		fmt.Println("Here are your options:")
		fmt.Println("Start tracking: \"S(tart)\"")
		fmt.Println("Options for the Tracker: \"O(ptions)\"")
	case label == "start":
		showText("clear")
		fmt.Println("Your time gets currently tracked...")
		fmt.Println("Your current options:")
		fmt.Println("Stop tracking: \"Stop\"")
	case label == "clear", label == "clr":
		fmt.Print("\033[H\033[2J")
	}
}
