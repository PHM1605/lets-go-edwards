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
