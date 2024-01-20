package main
import "github.com/chrpa-jakub/register-api/routes"
import "github.com/chrpa-jakub/register-api/database"

func main() {
 database.Init()
 routes.Init()
}
