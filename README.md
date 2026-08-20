## Init
```sh
go mod init lets-go-edwards
```
## Send different kinds of request
To send a GET request `curl -i localhost:4000`. \
To send a GET request with ID `curl -i localhost:4000/snippet/view/11`. \
To send a HEAD request (GET but returns only header): `curl --head localhost:4000/`.\
To send a POST request `curl -i -d "" localhost:4000/xxx` ("" means "request body is empty string").\
To display a form to create Snippet: `curl -i localhost:4000/snippet/create`.\
To submit the snippet creating form: `curl -i -d "" localhost:4000/snippet/create`.\
To reach a `static` asset: `curl -i localhost:4000/static/xxx.png`.

## Run this application
```sh
go run ./cmd/web -addr=":4000"
```
We can run with environment variables too:
```sh
export SNIPPETBOX_ADDR=":9999"
go run ./cmd/web -addr=$SNIPPETBOX_ADDR
```

## Pick static files
```sh
curl https://www.alexedwards.net/static/sb-v2.tar.gz | tar -xvz -C ./ui/static/
```

## Avoid hacking tricks 
Purpose: when user types `localhost:4000/static` or `localhost:4000/static/css` do NOT list files.\
Create `index.html` for each `static/**` folder
```sh
find ./ui/static -type d -exec touch {}/index.html \;
```

## Parse command-line flags to a variable
### Parse to a single variable
```sh
addr := flag.String("addr", ":4000", "HTTP network address") 
flag.Parse()
```
### Parse to a struct of many variables
```sh
type config struct {
  addr string
  staticDir string
}
var cfg config
flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
flag.Parse()
```

## Structured logger
To log as `JSON` instead of plain text:
```sh
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
```
To set `minimum log level` to `Debug` (so that all `logger.Debug()` is printed):
```sh
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
```
To print out **file location** of source code & **line number** of this log:
```sh
logger := logger.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
```
To redirect logging to a file:
```sh
go run ./cmd/web >> /tmp/web.log
```
To include a **stack trace** in logging when error occurs:
```sh
logger.Error(err.Error(), "method", method, "uri", uri, "trace", string(debug.Stack()))
```

## Go and MySQL
Install: `go get github.com/go-sql-driver/mysql@v1`.\
Create ONE `sql.DB` object, which is a **pool of many connections**.\
### Create a new Snippet
With `curl POST` request: `curl -iL -d "" http://localhost:4000/snippet/create`.\
Check if DB inserts that new entry: `SELECT id, title, expires FROM snippets;`.

## To handle middleware chains
```sh
go get github.com/justinas/alice
```

## To get GET request's parameters
```sh
title := r.URL.Query().Get("title")
```

## To get parameters from BOTH GET and POST requests
```sh
title := r.Form.Get("title") // for both GET and POST request
```

## To install & use automatic form parser
```sh
go get github.com/go-playground/form
```

## Session - scs will store sessions in Server side (MySQL)
Install libraries
```sh
go get github.com/alexedwards/scs/v2
go get github.com/alexedwards/scs/mysqlstore
```

## TLS
Create TLS certificates with (in `tls` folder): 
```sh 
go run /opt/homebrew/Cellar/go/1.26.3/libexec/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```
This will create `cert.pem` and `key.pem` files

## Encryption password with `bcrypt`
```sh
go get golang.org/x/crypto/bcrypt@latest
```

## Add CSRF-token to a field in HTML; to prevent CSRF attack
Install `nosurf` package
```sh
go get github.com/justinas/nosurf
```

## Build the program
```sh
go build -o /tmp/web ./cmd/web/
cp -r ./tls /tmp/
cd /tmp/
./web
```

## Go Test
```sh
go test -v ./cmd/web
```
To run all tests files (`*_test.go`) in current project directory: `go test ./...`\
To choose which test to run: `go test -v -run="^TestPing$" ./cmd/web/`\
To run sub-test: `go test -v -run="^TestHumanDate$/^UTC$" ./cmd/web`\
To skip a test (run the rest): `go test -v -skip="^TestHumanDate$" ./cmd/web/`\
To specify how many times you want to run each test before caching: `go test -count=1 ./cmd/web`\
To clean test cache: `go clean -testcache`\
To exit test immediately after 1st failure `Errorf()`: `go test -failfast ./cmd/web`\
To run tests in parallel: 
```sh
func TestPing(t *testing.T) {
  t.Parallel()
  ...
}
```
We can specify how many goroutines to run in parallel: `go test -parallel=4 ./...`\
To detect race condition in BOTH code & test: `go test -race ./cmd/web/`

## Integration test
**Definition**: use a "dev-DB" to test full flow.\
**Setup**: like in `MySQL.md`.\
In this code base, we check if `models.UserModel.Exists()` is working correctly.\
Run 
```sh
go test -v ./internal/models
```
NOTE: to skip integration test (notice the `testing.Short()` part in `users_test.go`)
```sh
go test -v -short ./...
```

## Profiling test coverage
```sh
go test -cover ./...
```
We can export to a file
```sh
go test -coverprofile=profile.out ./...
```
Then view the coverage of that file:
```sh
go tool cover -func=profile.out
go tool cover -html=profile.out
```
To view in deeper details, see also **how many times each statement is executed during the test** (the brighter green the more that statement is executed)
```sh
go test -covermode=count -coverprofile=profile.out ./...
go tool cover -html=profile.out
```
If some tests are in parallel (not in our repos):
```sh
go test -covermode=atomic -coverprofile=profile.out ./...
go tool cover -html=profile.out
```