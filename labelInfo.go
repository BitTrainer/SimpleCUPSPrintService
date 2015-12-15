package SimpleCUPSPrintService

import "fmt"

//LabelInfo is used to transport label printing instructions
type LabelInfo struct {
	LabelType      string `json:"labelType"`
	Title          string `json:"title"`
	Id             string `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	AdditionalInfo string `json:"additionalInfo"`
	Code           string
	CanPhotograph  bool `json:"canPhotograph"`
	Photography    string
	HasAllergies   bool `json:"hasAllergies"`
	Allergies      string
	Date           string `json:"Date"`
	Notes          string
}

//String creates a string representing an instance of LabelInfo
func (labelInfo LabelInfo) String() string {
	return fmt.Sprintf("LabelType = %s, Title = %s, Id= %s, FirstName = %s, LastName = %s, AdditionalInfo = %s",
		labelInfo.LabelType,
		labelInfo.Title,
		labelInfo.Id,
		labelInfo.FirstName,
		labelInfo.LastName,
		labelInfo.AdditionalInfo)
}
