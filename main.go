package main

import(
    "fmt" 
    "strings"
    "time"
)

var input string = ""                                   //Input String of the std_in
var done chan bool = make(chan bool)                    //Main waits for the goroutine
var start_done chan bool = make(chan bool)              //finish the Start ticker goroutine
var is_start chan bool = make(chan bool)              //Needa look how to make this better where it is involved 
var is_option chan bool = make(chan bool)              //Needa look how to make this better where it is involved 


func main(){
    //start up text and clear after programm ends
    defer Show_text("clear")

    //checks for the std_input
    Check_input_main()

    //if user is tracking time check its inputs
    for{
        select{
            case <- is_start:
                Check_input_start()
            case <- is_option:
                Check_input_option()
        }
    }

    <- done
}

/*-----------------Input Checks------------------*/

func Check_input_main(){
    Show_text("main")
    //label for breaking mid switch statement
    loop:
    for{
        //get std_input and make it lowercase
        fmt.Scan(&input)
        input = strings.ToLower(input)
        //decide through the input what the user can do
        switch{
            //Starting timer in thread 
            case input == "start", input == "s":
                go Start_timer()
                break loop
            //Change options of the timer 
            case input == "options", input == "o":
                go Options()
                break loop
            //Add further option like get overtime and current time worked and...

            //Wrong Input
            default:
                Show_text("clear")
                fmt.Println("WTF is going on!!")
                fmt.Println("Can't you read???")
                break
        }
    }
}

//Checks the input if the Timer was started 
func Check_input_start(){
    Show_text("start")
    loop:
    for{
        //get user Input and LowerCase it 
        fmt.Scan(&input)
        input = strings.ToLower(input)
        //decide through the input what the user can do
        switch{
            //Stops the timer
            case input == "stop", input == "s":
                go Stop_timer()
                break loop
            //Wrong input
            default:
                Show_text("clear")
                fmt.Println("WTF is going on!!")
                fmt.Println("Can't you read???")
                break
        }
    }
}

//Checks the input if the Timer was started 
func Check_input_option(){
    Show_text("option")
    loop:
    for{
        //get user Input and LowerCase it 
        fmt.Scan(&input)
        input = strings.ToLower(input)
        //decide through the input what the user can do
        switch{
            //Stops the timer
            case input == "", input == "":
                go Stop_timer()
                break loop
            //Wrong input
            default:
                Show_text("clear")
                fmt.Println("WTF is going on!!")
                fmt.Println("Can't you read???")
                break
        }
    }
}

/*-----------------Functions------------------*/

//Start the timer
func Start_timer(){ 
    start_show_time := time.NewTicker(5 * time.Second)           //Setups Ticker of X Seconds
    is_start <- true                                              //tell goroutine is running
    start_time := time.Now()                                        //get the time work started
    //fmt.Println("Tracks time...")                                 //Debug
    //Every time Ticker -> Prints current working time
    for{
        select{
            //finish goroutine
        case <-start_done:
            return

            //Print the current working time
        case cur_time := <-start_show_time.C:
            Show_text("start")
            fmt.Print("Your current Work time is: ", cur_time.Sub(start_time))
            break
        }
    }
}

//Stop the timer 
func Stop_timer(){
    for{
        done <- true
    }
}

/*-----------------TUI-Text------------------*/

//Shows text to the cmd according to the label
func Show_text(label string){
    label = strings.ToLower(label)
    switch{
        case label == "main":    
            Show_text("clear")
            fmt.Println("U wanna Track your time?")
            fmt.Println("Here are your options:")
            fmt.Println("Start tracking: \"S(tart)\"")
            fmt.Println("Options for the Tracker: \"O(ptions)\"")
            break
        case label == "start":
            Show_text("clear")
            fmt.Println("Your time gets currently tracked...")
            fmt.Println("Your current options:")
            fmt.Println("Stop tracking: \"Stop\"")
            break
        case label == "clear", label == "clr":
            fmt.Print("\033[H\033[2J")
            break
    }
}
