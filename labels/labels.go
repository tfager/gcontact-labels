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

	// Initialize the SVG encoder
	canvas := svg.New(file)
	width := 595
	height := 842
	canvas.Start(width, height) // A4 paper size in pixels (assuming 72 DPI)
	labelWidth := width / columns
	labelHeight := height / rows

	// Generate the address labels
	for i := 0; i < len(entries); i++ {
		entry := entries[i]

		// Calculate the position of the label
		x := i % columns * labelWidth
		y := i / columns * labelHeight

		// Write the contact information on the label
		style := "font-size:10"
		canvas.Text(x+10, y+20, entry.Name, style)
		canvas.Text(x+10, y+35, entry.StreetAddress, style)
		canvas.Text(x+10, y+50, fmt.Sprintf("%s %s", entry.PostalCode, entry.City), style)
	}

	// End the SVG encoding
	canvas.End()
}
