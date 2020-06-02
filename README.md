To quickly create a new project from this boilerplate, clone the repository and replace all instances of `adiachenko/go-scaffold` with the name of your package. Then, replace remaining instances of `go-scaffold` with the name of your binary.

## Design Philosophy

This boilerplate is heavily inspired by [Package Oriented Design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html), namely in how most logic goes into `internal/` folder and the way imports are managed within that folder.

#### internal/

Packages that provide support for the different features the project owns (CRUD, services or business logic):

- CAN import packages from `internal/platform`
- CANNOT import packages from parent directories.
- CANNOT import other packages at the same level.
- NOT allowed to panic an application
- NOT allowed to stop execution due to an error
- Minority of handling errors happen here. You're encouraged to wrap errors with context and return it to the caller

#### internal/platform/

Packages that provide internal foundational support for the project (database, messaging, various helpers etc.)

- CAN import other packages from `internal/platform/`
- CANNOT import packages from parent directories.
- NOT allowed to panic an application.
- NOT allowed to wrap errors.
- Return only root cause error values

#### General Tips

- Question imports for the sake of sharing existing types (constants, data structures etc. should be scoped to a package)
- Question imports to others packages at the same level.
- If a package wants to import another package at the same level:
  - Question the current design choices of these packages.
  - If reasonable, move the package inside the source tree for the package that wants to import it.
  - Use the source tree to show the dependency relationships.

## Quick Start

Run the following command to setup your project:

```
cp .env.example .env
docker-compose up -d
docker-compose exec app go mod download
```

Start serving HTTP requests:

```
docker-compose exec app go run main.go serve
```

Build a binary for release:

```
docker-compose exec app go mod tidy
docker-compose exec app go build
```

## Editor Configuration

IntelliJ Settings:
1. Navigate to **Languages & Frameworks → Go → Go Modules (vgo)** and ensure you have checked "Enable Go Modules (vgo)" integration
2. Navigate to **Languages & Frameworks → Go → GOPATH** and add `.go` folder to the list of project GOPATHs.
3. Navigate to **Tools → File Watchers** and add _go fmt_ configuration to your project.

## Project Dependencies

Make sure to give a cursory glance at the following packages' documentation:

- https://github.com/spf13/cobra (CLI boilerplate)
- https://github.com/go-chi/chi (router)
- https://github.com/Vinelab/tracing-go (distributed tracing)
- https://github.com/sirupsen/logrus (logging)
- https://github.com/Vinelab/go-reporting (error reporting)

## Console Commands

All functionality that this application provides including http server is executed via a command line interface using Cobra. You can install Cobra on your host machine for easier scaffolding of new commands:

```sh
go get github.com/spf13/cobra/cobra
```

You can generate new console commands as demonstrated below:

```go
cobra add serve
cobra add config
cobra add create -p 'configCmd'
```

## Routing

Add new handler functions to `routes/handlers` package

## Error Reporting

If you're calling `errors.New(msg)` to create your own errors, make sure to declare them at the top of the file like in the example below (notice how each name starts with _Err_):

```go
package mypkg

var (
  ErrOperationIsNotAllowed = errors.New("operation not allowed")
)
```

The caller can later check these errors like so:

```go
errors.Is(err, mypkg.ErrOperationIsNotAllowed)
```

You can annotate any errors that you receive with a message, preserving original stack trace using `fmt.Errorf` (make sure to use `%w` for variable substitution, else the caller won't be able to check the type of the error):

```go
fmt.Errorf("operation failed: %w", err)
```

You should only wrap errors in `internal` packages. However, you must not wrap errors in `internal/platform` directory (return the original error instead).

## RabbitMQ

### Queue Producer

```go
// Available types: TopicMessage, FanoutMessage, DirectMessage
msg := rabbitmq.TopicMessage(
  body, // json body of the message
  "example-exchange",
  "example-routing-key",
)

if err := rabbitmq.Produce(msg); err != nil {
  return err
}
```

### Queue Subscriber

```go
err := rabbitmq.Subscribe(rabbitmq.Subscriber{
  ExchangeName: "example-exchange",
  ExchangeType: rabbitmq.ExchangeTopic,
  QueueName:    "example-queue-name",
  BindingKeys: []string{
    "example-binding-key",
  },
  Handler: func(delivery *amqp.Delivery) {
    // Unserialize payload into a struct and pass it to a package from /internal dir for processing...
    if err := internalpkg.DoSomething(payload); err != nil {
      rabbitmq.ProcessError(delivery, err)
    } else {
      rabbitmq.ProcessSuccess(delivery)
    }

  },
})

if err != nil {
  logrus.WithError(err).Fatal(err.Error())
}
```

Queue subsriber is to be created in a dedicated console command (each subscriber gets its own command).
