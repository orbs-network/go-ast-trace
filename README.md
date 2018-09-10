# go-ast-trace

> AST-Based Tracing for Go

`go-ast-trace` is a tool for debugging and investigating issues in programs written in golang.

Investigation is based on automatic logging of certain events, running the program and then analyzing the logs. The CLI tool `go-ast-trace` focuses only on the generation of logs, their analysis is out of scope.

`go-ast-trace` works by manipulating the source code of your program and injecting the logging commands directly inside. It relies on a great feature of golang baked right in the standard library, the ability to [parse](https://golang.org/pkg/go/parser) golang sources to [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree), manipulate the AST and then [print](https://golang.org/pkg/go/printer) it back to source code.

## Installation

1. Go get the tool with `go get github.com/orbs-network/go-ast-trace`

2. Make sure `GOPATH/bin` is in your path.

3. Verify installation by typing in terminal `go-ast-trace` (it should show CLI usage).

## Usage

1. Open the root directory of your golang project in terminal.

2. Make sure the current state of your project is committed to git.

3. Run `go-ast-trace <trace-type> <input-files>`

    For example: `go-ast-trace locks *.go`
    
4. Run your project.

5. Trace logs will be printed to stdout, analyze them.

6. The tool will alter the source code of your project. Once finished, git reset everything back to normal.

&nbsp;
## Tracing mutex and channel locks

#### When is this useful?

Debugging synchronization issues is one of the most difficult and time consuming tasks a developer can face. If your project suffers from a deadlock for example, narrowing down the exact locking operation causing the deadlock is not easy. By tracing all locking operation to log, you can identify the exact lock causing the issue and the stack trace leading to it.

#### How to add lock tracing?

Run `go-ast-trace locks *.go` on your project source files.

#### How will your source code change?

Mutex locks (both `sync.Mutex` and `sync.RWMutex`) and blocking operations on channels (reads and writes), such as this:

```go
func normalLockExample(mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()
}
```

Will change after running the CLI tool to:

```go
func normalLockExample(mutex *sync.Mutex) {
	astraceInjection.BeforeLock(15352856648520921629)
	mutex.Lock()
	astraceInjection.AfterLock(15352856648520921629)
	
	defer mutex.Unlock()
}
```

Note: The number `15352856648520921629` is a random ID identifying this specific lock (to help with measuring time, etc).

#### What will the logs look like?

```cgo
ASTRACE 10:48:22.933061 BeforeLock main.normalLockExample [examples/mutex.go:19 examples/mutex.go:12 runtime/proc.go:198 runtime/asm_amd64.s:2361]
ASTRACE 10:48:22.933201 AfterLock main.normalLockExample 2ms [examples/mutex.go:25 examples/mutex.go:12 runtime/proc.go:198 runtime/asm_amd64.s:2361]
ASTRACE 10:48:23.157577 BeforeLock main.channelReadExample [examples/channel.go:18 examples/channel.go:12 runtime/proc.go:198 runtime/asm_amd64.s:2361]
ASTRACE 00:48:23.157758 AfterLock main.channelReadExample 3ms [examples/channel.go:24 examples/channel.go:12 runtime/proc.go:198 runtime/asm_amd64.s:2361]
```

* An event before acquiring the lock and an event after acquiring the lock.
* Timestamp for each event (accurate to a microsecond).
* The name of the function that performed this lock operation.
* How long did the goroutime block and wait before successfully acquiring the lock (in milliseconds).
* Compact stack trace of the call showing only file names and line numbers (call stack with last 4 callers only).

&nbsp;
## License
MIT
