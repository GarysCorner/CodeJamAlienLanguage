//File:		codejamalienlanguage.go
//Author:	Gary Bezet
//Date:		2016-06-13
//Desc:		This program is designed to solve Google Code Jam "Alien Language"  I believe I could more easily solve this by constructing regex, but it felt like cheating
//Problem:	https://code.google.com/codejam/contest/90101/dashboard

package main

import (
		"fmt"
		"time"
		"os"
		"flag"
		"bufio"
		"strings"
		"strconv"
	)

//global variables
var infileopt, outfileopt string  //input and output filenames
var infile, outfile *os.File  //input and output file pointers
var totalcases int  //number of cases
var dictlen int //length of dictionary
var wordlen int //length of each alien word 


var dict []string
var testcases []testcase


//structures
type testcase struct {
	casenum int  
	letter [][]rune	
	solution int
	solvetime time.Duration
}





//program entry point
func main() {

	starttime := time.Now()  //start time for stats

	defer infile.Close()
	defer outfile.Close()

	initflags()  //initialize the command line args
	
	openFiles() //open the files
	
	processFile()
	
	startsolve := time.Now() //time we started solving
	for _, v := range testcases {
		v.solve()
		printErrln("Solved #", v.casenum, "in", v.solvetime," Ans=", v.solution)
		fmt.Fprintf(outfile,"Case #%d: %d\n", v.casenum, v.solution)
	}
	printErrln(totalcases, "cases solved in", time.Now().Sub(startsolve))
	

	
	printErrln("FINISHED!  Elapsed: ", time.Now().Sub(starttime))
	
}


//get the flags from command line
func initflags() {
	flag.StringVar(&infileopt, "if", "", "Input file (required)")
	flag.StringVar(&outfileopt, "of", "-", "Output file, defaults to stdout" )

	flag.Parse()

	if infileopt == "" {
		printErrln("You must supply an input file\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	

}

//print error to console
func printErrln( line ...interface{} ) {
	fmt.Fprintln( os.Stderr, line... )
}

func openFiles() {
	
	var err error
	
	infile, err = os.Open(infileopt)

	if err != nil {
		printErrln( "Error:  Could not open:  ", infileopt)
		printErrln( "\tError: ", err  )
		os.Exit(2)
	}

	if outfileopt == "-"  {
		outfile = os.Stdout
		outfileopt = "Stdout"
	} else {
		outfile, err = os.Create(outfileopt)

		if err != nil {
			printErrln( "Error:  Could not create:  ", outfileopt)
			printErrln( "\tError: ", err  )
			os.Exit(3)
		} 
	}

	printErrln("InFile:\t", infileopt)
	printErrln("OutFile:\t", outfileopt, "\n")
		
}


func processFile() {  //process the input file into data structure

	proctime := time.Now() //for time to load data

	var err error
	var line string
	
	reader := bufio.NewReader(infile)
	
	line, err = reader.ReadString('\n')
	if err != nil {
		printErrln("Couldn't read first line from:  ", infileopt)
		printErrln("\tError:  ", err )
		os.Exit(4)
		
	}
	
	initdata := strings.Split(line, " ")  //get an array with the initial data from line 1
	
	if len(initdata) != 3 {
		printErrln( "Could not read 3 values from first line from:  ", infileopt )
		os.Exit(5)
	}
	
	wordlen, err = strconv.Atoi( strings.TrimSpace(initdata[0]) )  //input wordlength
	if err != nil  { //if error reading number of cases
		printErrln("Couldn't read word length numbers from:  ", infileopt)
		printErrln("\tError:  ", err )
		os.Exit(5)
	}
	
	dictlen, err = strconv.Atoi( strings.TrimSpace(initdata[1]) )  //input dictionary length
	if err != nil  { //if error reading number of cases
		printErrln("Couldn't read dictionary length numbers from:  ", infileopt)
		printErrln("\tError:  ", err )
		os.Exit(5)
	}
	
	totalcases, err = strconv.Atoi( strings.TrimSpace(initdata[2]) )  //input number of test cases
	if err != nil  { //if error reading number of cases
		printErrln("Couldn't read case numbers from:  ", infileopt)
		printErrln("\tError:  ", err )
		os.Exit(5)
	}
	
	
	dict = make( []string, dictlen )
	for i := 0; i < dictlen; i++ {  //read the dictionary in
		
		line, err = reader.ReadString('\n')
		if err != nil {  
			printErrln("Fatal error reading dictionary from:  ", infileopt)
			printErrln("\tError:  ", err)
			os.Exit(6)
		}
		
		dict[i] = strings.TrimSpace( line )
		
		if len(dict[i]) != wordlen {  //make sure wordlength is correct
			printErrln("Bad word length (", len(dict[i]), ") on word#:  ", i+1)
			os.Exit(7)
		}
		
	}
	
	
	testcases = make([]testcase, totalcases)
	for i:= 0; i < totalcases; i++ {  //read in cases
	
		var inbyte rune
	
		testcases[i].casenum = i + 1
		
		testcases[i].letter = make([][]rune, wordlen)
	
	
		for c := 0; c < wordlen; c++ {  //read single character
			inbyte, _, err = reader.ReadRune()
		
			if err != nil {
				printErrln("Error reading case#:", i+1, "from file:  ", infileopt)
				printErrln("\tError:  ", err)
				os.Exit(8)
			}
		
			if inbyte != '(' {  //if single byte
			
				testcases[i].letter[c] = make([]rune, 1)
				testcases[i].letter[c][0] = inbyte
				
			} else {  //multibyte
				
				inbytes, err := reader.ReadString(')')
				if err != nil {
					printErrln("Error reading case#:", i+1, "from file:  ", infileopt)
					printErrln("\tError:  ", err)
					os.Exit(8)
				}
				
				inbytes = inbytes[0:len(inbytes)-1]  //chomp trailing ")"
				testcases[i].letter[c] = make([]rune, len(inbytes))
				
				for iter, inbyte := range inbytes {
				
					testcases[i].letter[c][iter] = inbyte
				
				}
				
				
			}
		
		}
		
		_, _, _ = reader.ReadRune()  //discard newline character	

	
	}
	
	printErrln("Input file processed in", time.Now().Sub(proctime))

}

//solve a case
func (self *testcase) solve() {
	
	starttime := time.Now()	//starttime for this solution


	toploop:
	for _, w := range dict {  //test all words
	
		for i, l := range w {  //test all letters
			if search( l, self.letter[i]) == false {
				continue toploop //if letter is not found then skip on to the next dictionary word
			}
		}		
		
		self.solution++  //if we get to this point solution found, increment the variable till all solutions found
	
	}

	self.solvetime = time.Now().Sub(starttime)

}

//return bool/true if needle found
func search( needle rune, haystack []rune) bool {
	for _, v := range haystack {
		if v == needle {
			return true  //return true if
		}
	}
	
	return false //if no match found return false
}











