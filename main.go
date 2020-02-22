package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

func main() {
	speech := htgotts.Speech{
		Folder:   "audio",
		Language: "en",
	}

	var outTE *walk.TextEdit
	var generatedNums []int
	MainWindow{
		Title:   "Get ready for interview",
		MinSize: Size{100, 100},
		MaxSize: Size{200, 200},
		Layout:  VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TextEdit{AssignTo: &outTE, ReadOnly: true},
				},
			},
			PushButton{
				Text: "Play question",
				OnClicked: func() {

					// Сигнал по нажатию на кнопку
					questions, err := getDataFromFile("text_for_converting_to_speech.txt")
					if err != nil {
						log.Fatal(err, "unable to get data from file: %s",
							"text_for_converting_to_speech.txt",
						)
					}

					for len(generatedNums) <= len(questions) {
						log.Info("generated nums", generatedNums)
						randomNumber := createRandomNumber(0, len(questions))
						if validateRandomNumber(generatedNums, randomNumber) {
							generatedNums = append(generatedNums, randomNumber)
							speech.Speak(questions[randomNumber])
							outTE.SetText(questions[randomNumber])
							break
						} else if len(generatedNums) == len(questions) {
							speech.Speak("Questions are over")
							outTE.SetText("Questions are over!!!!!!!!!")
							generatedNums = nil
							break
						} else {
							continue
						}
					}

				},
			},
		},
	}.Run()
}

func validateRandomNumber(array []int, randomNumber int) bool {
	log.Info("array", array)
	log.Info("randomNumber", randomNumber)
	for i := range array {
		if array[i] == randomNumber {
			return false
		}
	}

	return true
}

func createRandomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

type valueMinMax struct {
	min int
	max int
}

func getDataFromFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, karma.Format(err, "unable to open a file")
	}

	defer file.Close()

	dataFromFile, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, karma.Format(err, "unable to get data from file")
	}

	textFromFile := strings.Split(string([]byte(dataFromFile)), "\n")
	return textFromFile, nil
}
