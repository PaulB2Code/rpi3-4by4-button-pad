package KeyPad

import (
	"errors"
	"fmt"
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

func (kp *Keypad) Close() {
	// Unmap gpio memory when done
	rpio.Close()
}

type Keypad struct {
	Reading bool
}

func New() (Keypad, error) {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		return Keypad{}, err
	}

	for i, val := range row {
		//	fmt.Println("Row ", i, "  -  ", val)
		rowPin = append(rowPin, rpio.Pin(val))
		rowPin[i].Output()
		rowPin[i].High()
	}
	for i, val := range column {
		//	fmt.Println("Col ", i, "  -  ", val)
		colPin = append(colPin, rpio.Pin(val))
		colPin[i].Input()
		colPin[i].PullDown()
	}
	return Keypad{}, nil
}

func getLetter(row int, col int) (string, error) {
	if len(keys) < row || len(keys) < col {
		return "", errors.New(fmt.Sprintf("No data for row %v and col %v", row, col))
	}
	val := keys[row][col]

	return val, nil
}

var ticker *time.Ticker

func (kp *Keypad) TrackClicked(timeScanning time.Duration, c chan int) {
	ticker = time.NewTicker(timeScanning)
	go func() {
		for _ = range ticker.C {
			if !kp.Reading {
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
func (kp *Keypad) StopTracking() {
	ticker.Stop()
}
func (kp *Keypad) restartTracking(timeScanning time.Duration) {
	ticker = time.NewTicker(timeScanning)
}

func (kp *Keypad) GetValueWithColumn(i int) (string, error) {
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
