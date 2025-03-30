// Local database for files.
//
// Example:
//
//	// Make a connection to the local DB
//	db := docdb.New("./app.db", "./files")
//	// Insert a file
//	doc, err := db.Insert(myFile)
//	// Retrieve that file using its ID
//	file, err := db.FindOne(doc.ID)
//	// Get 5 files sorted alphabetically by their name
//	files, err := db.FindMany(5, 0, col.Name, true)
package docdb
