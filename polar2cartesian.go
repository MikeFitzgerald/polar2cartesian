/*
LESSONS

1. Eventhough the polar and cartesian structs happen to have the same types data fields
   (2 float64 values), it is not possible to implicitly convert between them.  
   This is intentional in the Go language design and supports defensive programming.

2. If a package has one or more init() functions they are automatically executed BEFORE
   the main package's main() function is called.

3. init() functions must NEVER be called explicitly.

4. Channels are modeled on Unix pipes and provide one or two-way communication of data.

5. Chanels befave like FIFO queues, hence preserving order or the items sent into them.

6. Channels are created with the make() function
       Example:
           messages := make(chan string, 10)        (where 10 is the buffer size)
    If buffer size reaches 0 (meaning all channels are used, the channel blocks until 
    at least one item is received from it).

7. A channel of buffer size of 0 can only send an item if the other end is waiting for 
   an item.

 */

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
//	"time"
)

type polar struct {
	radius float64
	theta float64
}

type cartesian struct {
	x float64
	y float64
}

var prompt = "Enter a radius and an angle (in degrees) e.g., 12.5, 90, or %s to quit."

func init() {
	if runtime.GOOS == "windows" {
		prompt = fmt.Sprintf(prompt, "Ctrl+Z, Enter")
	} else {
		prompt = fmt.Sprintf(prompt, "Ctrl+D")
	}
}

func main () {
	questions := make(chan polar)
	defer close(questions)
	answers := createSolver(questions)
	defer close(answers)
	interact(questions, answers)
}

func createSolver(questions chan polar) chan cartesian {
	answers := make(chan cartesian)
	go func() {
		for {
			polarCoord := <-questions
			// time.Sleep(5000 * time.Millisecond)
			theta := polarCoord.theta * math.Pi / 180.0 // degrees to radians
			x := polarCoord.radius * math.Cos(theta)
			y := polarCoord.radius * math.Sin(theta)
			answers <- cartesian{x, y}
		}
	}()
	return answers
}

const result = "Polar radius=%.02f theta=%.02f degrees -> Cartesian x=%.02f y =%.02f\n"

func interact(questions chan polar, answers chan cartesian) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)
	for {
		fmt.Printf("Radius and angle: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		var radius, theta float64
		if _, err := fmt.Sscanf(line, "%f %f", &radius, &theta); err != nil {
			fmt.Fprintln(os.Stderr, "invalid input")
			continue
		}
		questions <- polar{radius, theta}
		coord := <-answers
		fmt.Printf(result, radius, theta, coord.x, coord.y)
	}
	fmt.Println()
}