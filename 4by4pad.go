package KeyPad

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const ()

var row = [4]uint8{4, 17, 27, 22}     //{7, 11, 13, 15}
var column = [4]uint8{18, 23, 24, 25} //12, 16, 18, 2}

var keys = [4][4]string{
	{"1", "2", "3", "A"},
	{"4", "5", "6", "B"},
	{"7", "8", "9", "C"},
	{"*", "0", "#", "D"}}
var colPin = make([]rpio.Pin, 0, 0)
var rowPin = make([]rpio.Pin, 0, 0)

var reading = false

func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	for i, val := range row {
		fmt.Println("Row ", i, "  -  ", val)
		rowPin = append(rowPin, rpio.Pin(val))
		rowPin[i].Output()
		rowPin[i].High()
	}
	for i, val := range column {
		fmt.Println("Col ", i, "  -  ", val)
		colPin = append(colPin, rpio.Pin(val))
		colPin[i].Input()
		colPin[i].PullDown()
	}

	c := make(chan int)
	trackClicked(100*time.Millisecond, c)
	defer stopTracking()

	for {
		val := <-c
		reading = false

		letter, err := getValueWithColumn(val)
		if err != nil {
			log.Println("[ERROR], ", err)
		}
		log.Println("Letter clicked", letter)

	}

	for {
		time.Sleep(100 * time.Millisecond)
		//fmt.Println("Col 1  -", col1.Read(), "  Col2 - ", col2.Read(), "Col 3  -", col3.Read(), "  Col4 - ", col4.Read())
		for i, val := range colPin {
			if val.Read() == 1 {
				var rowId int
				//Test to know whos button is pushed
				for ii := 0; ii < len(rowPin); ii++ {
					rowPin[ii].Low()
					time.Sleep(10 * time.Millisecond)
					if val.Read() == 0 {
						rowId = ii
						goto endblock
					}
				}
				log.Println("Error Reading")
			endblock:

				letter, err := getLetter(rowId, i)
				if err != nil {
					log.Println("[ERROR], ", err)
				}
				log.Println("Letter clicked", letter)
				//Set High All
				time.Sleep(300 * time.Millisecond)
				for ii := 0; ii < rowId+1; ii++ {
					rowPin[ii].High()
				}
				goto end
			}
		}
	end:
	}
}

func getLetter(row int, col int) (string, error) {
	if len(keys) < row || len(keys) < col {
		return "", errors.New(fmt.Sprintf("No data for row %v and col %v", row, col))
	}
	val := keys[row][col]

	return val, nil
}

var ticker *time.Ticker

func trackClicked(timeScanning time.Duration, c chan int) {
	ticker = time.NewTicker(timeScanning)
	go func() {
		for _ = range ticker.C {
			if !reading {
				for i, val := range colPin {
					if val.Read() == 1 {
						c <- i
						goto end
					}
				}
			}
		end:
		}
	}()

}
func stopTracking() {
	ticker.Stop()
}
func restartTracking(timeScanning time.Duration) {
	ticker = time.NewTicker(timeScanning)
}

func getValueWithColumn(i int) (string, error) {
	var rowId int
	//Test to know whos button is pushed
	for ii := 0; ii < len(rowPin); ii++ {
		rowPin[ii].Low()
		time.Sleep(10 * time.Millisecond)
		if colPin[i].Read() == 0 {
			rowId = ii
			goto endblock
		}
	}
	return "", errors.New("Error Reading second number")
endblock:

	letter, err := getLetter(rowId, i)
	if err != nil {
		return "", err
	}
	//Set High All
	time.Sleep(300 * time.Millisecond)
	for ii := 0; ii < rowId+1; ii++ {
		rowPin[ii].High()
	}
	return letter, nil
}
