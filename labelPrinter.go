package SimpleCUPSPrintService

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	maxNameLength = 16
)

//LabelPrinter ...
type LabelPrinter struct {
	isToBeSentToPrinter bool
	printerName,
	templatePath,
	outputPath string
	errorLog,
	infoLog *log.Logger
	labelTemplate *template.Template
}

//NewLabelPrinter  ...
func NewLabelPrinter(isToBeSentToPrinter bool,
	printerName,
	templatePath,
	outputPath string,
	errorLog,
	infoLog *log.Logger) *LabelPrinter {
	templatePathParts := strings.Split(templatePath, string(os.PathSeparator))
	templateName := templatePathParts[len(templatePathParts)-1]
	labelTemplate, err := template.New(templateName).ParseFiles(templatePath)
	if err != nil {
		return nil
	}
	var result = LabelPrinter{isToBeSentToPrinter,
		printerName,
		templatePath,
		outputPath,
		errorLog,
		infoLog,
		labelTemplate}

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

func prepareLabelInfoForPrinting(lableInfo *LabelInfo) {
	lableInfo.Code = strings.ToUpper(formatLabelPart(maxNameLength, lableInfo.Id))
	lableInfo.Date = time.Now().Format("2006-01-02 15:04:05")
	lableInfo.FirstName = formatLabelPart(maxNameLength, lableInfo.FirstName)
	lableInfo.LastName = formatLabelPart(maxNameLength, lableInfo.LastName)
	lableInfo.Notes = lableInfo.AdditionalInfo
	if lableInfo.CanPhotograph {
		lableInfo.Photography = ""
	} else {
		lableInfo.Photography = "Don't photograph"
	}
	if lableInfo.HasAllergies {
		lableInfo.Allergies = "Allergies"
	} else {
		lableInfo.Allergies = ""
	}
}

func writeLabel(lableInfo *LabelInfo,
	labelTemplate *template.Template,
	outputPath string) (err error) {
	prepareLabelInfoForPrinting(lableInfo)
	w, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer w.Close()
	labelTemplate.Execute(w, lableInfo)
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

		err = writeLabel(&printJob, labelPrinter.labelTemplate, labelFileName)
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
