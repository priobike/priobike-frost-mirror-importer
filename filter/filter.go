package filter

import (
	"importer/log"
	"importer/things"
	"strconv"
	"strings"

	"github.com/thedatashed/xlsxreader"
)

const nodeString = "Knoten"
const connectionString = "betroffene connections"

// Filters synced things
func FilterThings() {
	log.Info.Println("Create filter...")

	// An arrays that contain all blocked things by their name (complete nodes, specific connections)
	var filterStringsNodes, filterStringsConnections = createFilterStrings()

	log.Info.Println("Filter created")
	// fmt.Println(filterStringsNodes, filterStringsConnections)

	// Filter things
	filterThings(filterStringsNodes, filterStringsConnections)

	log.Info.Println("Things filtered")
}

// Creates the arrays with filter values for complete nodes and specific connections
func createFilterStrings() ([]string, []string) {
	var filterStringsConnections = []string{}
	var filterStringsNodes = []string{}

	// Create an instance of the reader by opening the excel file
	xl, _ := xlsxreader.OpenFile("./Fehleranalyse_PrioBike_für_TUD.xlsx")

	// Ensure the file reader is closed once utilised
	defer xl.Close()

	// Iterate on the rows of data
	var indexNode = -1
	var indexConnection = -1
	for row := range xl.ReadRows(xl.Sheets[0]) {
		// Get the cell index on first row
		if row.Index == 1 {
			for index, cell := range row.Cells {
				if cell.Value == nodeString {
					indexNode = index
				}
				if cell.Value == connectionString {
					indexConnection = index
				}
			}
			if indexNode == -1 {
				panic("Node Column not found")
			}
			if indexConnection == -1 {
				panic("Connection Column not found")
			}
			continue
		}

		// Get the values and add them to the filter arrays
		nodeName := row.Cells[indexNode].Value
		nodeConnection := row.Cells[indexConnection].Value

		// Add node since all connections are blocked
		if nodeConnection == "alle" {
			filterStringsNodes = append(filterStringsNodes, nodeName)
			continue
		}

		// Create filterstring for sequences like 1;2;3;...;x
		splittedConnections := strings.Split(nodeConnection, ";")
		for _, splittedConnectionString := range splittedConnections {
			// Create filter strings for single values like 1-10
			if strings.Contains(splittedConnectionString, "-") {
				splittedNumbers := strings.Split(splittedConnectionString, "-")
				start, err := strconv.Atoi(splittedNumbers[0])
				if err != nil {
					panic(err)
				}
				end, err := strconv.Atoi(splittedNumbers[1])
				if err != nil {
					panic(err)
				}
				for i := start; i <= end; i++ {
					filterStringsConnections = append(filterStringsConnections, nodeName+"_"+strconv.Itoa(i))
				}
			} else {
				filterStringsConnections = append(filterStringsConnections, nodeName+"_"+splittedConnectionString)
			}
		}
	}
	return filterStringsNodes, filterStringsConnections
}

// Executes the filtering on the things map
func filterThings(filterStringsNodes []string, filterStringsConnections []string) {
	things.Things.Range(func(key, value interface{}) bool {

		name := value.(things.Thing).Name
		nodeName := value.(things.Thing).Name

		if strings.Contains(name, "_") {
			nodeName = strings.Split(name, "_")[0]
		}

		// Search in filter for all connections
		if stringInSlice(nodeName, filterStringsNodes) {
			// println(nodeName)
			things.Things.Delete(key)
			return true
		}

		// Search in filter for exact connection.´
		if stringInSlice(name, filterStringsConnections) {
			// println(name)
			things.Things.Delete(key)
			return true
		}

		return true
	})
}

// Same as array contains
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
