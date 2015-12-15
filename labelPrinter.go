package SimpleCUPSPrintService

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	titleMarker      = "<TITLE>"
	accessCodeMarker = "<CODE>"
	firstNameMarker  = "<FIRSTNAME>"
	lastNameMarker   = "<LASTNAME>"
	notesMarker      = "<NOTES>"
	dateMarker       = "<DATE>"
	photoMarker      = "<photography>"
	allergyMarker    = "<allergies>"
	maxNameLength    = 16
)

//LabelPrinter ...
type LabelPrinter struct {
	isToBeSentToPrinter bool
	printerName,
	templatePath,
	outputPath string
	errorLog,
	infoLog *log.Logger
}

//NewLabelPrinter  ...
func NewLabelPrinter(isToBeSentToPrinter bool,
	printerName,
	templatePath,
	outputPath string,
	errorLog,
	infoLog *log.Logger) *LabelPrinter {
	var result = LabelPrinter{isToBeSentToPrinter,
		printerName,
		templatePath,
		outputPath,
		errorLog,
		infoLog}
	return &result
}

func formatLabelPart(maxLength int, labelPart string) string {
	partLength := len(labelPart)
	if partLength >= maxLength {
		return labelPart[0:maxLength-3] + "..."
	}
	return strings.Title(fmt.Sprintf("%s", labelPart))
}

func sendLabelToPrinter(filePath string, printerName string) (err error) {
	args := []string{filePath, printerName, "media=w36h89"}
	err = exec.Command("lpr", args[0], "-P", args[1], "-o", args[2]).Run()
	return err
}

func writeLabel(lableInfo LabelInfo,
	templatePath string, outputPath string) (err error) {

	inputStream, err := os.Open(templatePath)
	if err != nil {
		err = fmt.Errorf("Error opening print template at %s. Got error %v \n", templatePath, err)
		return
	}
	outputStream, err := os.Create(outputPath)
	if err != nil {
		err = fmt.Errorf("Error creating label file  %v \n", err)
		return
	}
	defer outputStream.Close()
	defer inputStream.Close()

	scanner := bufio.NewScanner(inputStream)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, titleMarker) {
			line = strings.Replace(line, titleMarker, strings.Title(lableInfo.Title), -1)
		} else if strings.Contains(line, accessCodeMarker) {
			line = strings.Replace(line, accessCodeMarker, strings.ToUpper(formatLabelPart(maxNameLength, lableInfo.Id)), -1)
		} else if strings.Contains(line, dateMarker) {
			line = strings.Replace(line, dateMarker, time.Now().Format("2006-01-02 15:04:05"), -1)
		} else if strings.Contains(line, firstNameMarker) {
			line = strings.Replace(line, firstNameMarker, formatLabelPart(maxNameLength, lableInfo.FirstName), -1)
		} else if strings.Contains(line, lastNameMarker) {
			line = strings.Replace(line, lastNameMarker, formatLabelPart(maxNameLength, lableInfo.LastName), -1)
		} else if strings.Contains(line, notesMarker) {
			line = strings.Replace(line, notesMarker, lableInfo.AdditionalInfo, -1)
		} else if strings.Contains(line, photoMarker) {
			if lableInfo.CanPhotograph {
				line = strings.Replace(line, photoMarker, "", -1)
			} else {
				line = strings.Replace(line, photoMarker, "Don't photograph", -1)
			}
		} else if strings.Contains(line, allergyMarker) {
			if lableInfo.HasAllergies {
				line = strings.Replace(line, allergyMarker, "Allergies", -1)
			} else {
				line = strings.Replace(line, allergyMarker, "", -1)
			}
		}
		_, err = io.WriteString(outputStream, line+"\n")
		if err != nil {
			return
		}
	}

	if scanner.Err() != nil {
		err = fmt.Errorf("Error reading label template")
		return
	}
	return nil
}

//Print prints labels
func (labelPrinter *LabelPrinter) Print(printJob LabelInfo,
	numberOfCopies int) (err error) {

	for i := 1; i <= numberOfCopies; i++ {
		labelPrinter.infoLog.Printf("Printing label: %s ... \n", printJob)
		timeStamp := time.Now()
		mockGUID := strconv.Itoa(timeStamp.Nanosecond())
		labelFileName := labelPrinter.outputPath + printJob.Id + "_" + timeStamp.Format("2006_01_02_15_04_05") + "_" + mockGUID + ".ps"
		if i > 1 {
			printJob.Title = "Parent / Guardian Copy"
		}

		err = writeLabel(printJob, labelPrinter.templatePath, labelFileName)
		if err != nil {
			labelPrinter.errorLog.Printf("Error printing label  %v \n", err)
			return err
		}
		if !labelPrinter.isToBeSentToPrinter {
			return nil
		}

		err = sendLabelToPrinter(labelFileName, labelPrinter.printerName)
		if err != nil {
			labelPrinter.errorLog.Printf("Error printing label  %v \n", err)
			return err
		}

		labelPrinter.infoLog.Println("Printed label successfully. Please collect at Label Station.")
	}
	return nil
}
