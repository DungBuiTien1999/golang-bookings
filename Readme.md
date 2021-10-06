# Bookings and Reservations

This is the repository for my bookings and reservations project demo.

- Built in Go version 1.17
- Uses the [chi router](github.com/go-chi/chi/v5)
- Uses [alex edwards SCS](github.com/alexedwards/scs/v2)
- Uses [nosurf](github.com/justinas/nosurf)

This project has only purpose to learn Web application with Golang language

Some commands:
go test -coverprofile=coverage.out && go tool cover -html=coverage.out (this command use to show test coverage)
chmod +x run.sh
./run.sh
go test -v ./... (root directory)

migrate database before run
create a new schema with name is golangbookings (whatever is up for you but remember edit code connect database)
`soda migrate` (notice: installs buffalo, link: https://gobuffalo.io/en/docs/getting-started/installation/)