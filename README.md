#Table reservation app

##Prerequisites

Create .env file in root directory and add following values:
```dotenv
POSTGRES_URI=postgres://<username>:<password>@127.0.0.1:5432/<db_name>
SESSION_SECRET=<any secret>
```

Use `go run ./cmd/app` to run project

###Instructions
Please follow my style of code (clean architecture)