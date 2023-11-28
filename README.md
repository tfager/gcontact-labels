# gcontact-labels

Simple CLI to generate SVG for printing mailing labels from Google contacts.

## Usage

```
# On first run you'll get local web server address, follow that to authorize Google API
go run main/main.go groups

# Pick the desired group (its resource ID) from the list, and:
go run main/main.go -group=contactGroups/xxxxx contacts

# The result will be in address_labels.svg
```
