package DB

type IDB interface {
	Init()        // Initialize the database connection
	Health() bool // Check if the database is connected
	Close()       // Close the database connection
}
