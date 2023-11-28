package labels

import (
	"fmt"
	"os"

	contacts "gcontact-labels/contacts"

	svg "github.com/ajstarks/svgo"
)

func GenerateAddressLabels(entries []*contacts.Contact, rows, columns int) {
	// Create a new SVG file
	file, err := os.Create("address_labels.svg")
	if err != nil {
		fmt.Println("Failed to create SVG file:", err)
		return
	}
	defer file.Close()

	canvas := svg.New(file)
	// A4 paper size in pixels (assuming 300 DPI)
	width := 2480
	height := 3508
	canvas.Start(width, height)
	leftMargin := 60
	topMargin := 60
	labelWidth := (width - 2*leftMargin) / columns
	labelHeight := (height - 2*topMargin) / rows
	fontSize := 40
	style := fmt.Sprintf("font-size:%dpx; font-family:Liberation Sans, Arial", fontSize)
	rowSpacing := 50

	// Generate the address labels
	for i := 0; i < len(entries); i++ {
		entry := entries[i]

		// Calculate the position of the label
		x := i % columns * labelWidth
		y := i / columns * labelHeight

		// Write the contact information on the label
		canvas.Text(x+leftMargin, y+topMargin, entry.Name, style)
		canvas.Text(x+leftMargin, y+topMargin+rowSpacing, entry.StreetAddress, style)
		canvas.Text(x+leftMargin, y+topMargin+2*rowSpacing, fmt.Sprintf("%s %s", entry.PostalCode, entry.City), style)
	}

	// End the SVG encoding
	canvas.End()
}
