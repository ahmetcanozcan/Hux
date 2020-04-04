
![Test Automation](http://ahmetcanozcan.github.io/assets/img/huxlogo.png)


Hux is a channel based event abstraction for web-sockets

## Installation

Use `go get` commadn to install  Hux

```bash
go get -u github.com/ahmetcanozcan/hux
```

## Usage

Hux provides both server-side and client-side libraries.

In your go file, add these blocks
```go


// blocks of code
func main() {
  // block of codes
  hux.Initialize() // Initialize hux
  go handleHub() // Then start handling sockets in another goroutine
  http.ListenAndServe(":8080", nil)
}

func handleHub() {
  h := hux.GetHub() // Get main hub.
  for {
    select {
    //When a web socket connected, this block executes
    case sck := <-h.SocketConnection: 
      go handleSocket(sck)//Handle socket
    case sck := <-h.SocketDisconnection:
      go handleDisconnection(sck)
    }
  }
}

// blocks of code

func handleSocket(sck *hux.Socket) {
  fmt.Println("Socket connected.")
  for {
    select {
    // When the socket sends Hello event,
    // this block will be executed
    case data := <-sck.GetEvent("Hello"):
      fmt.Println(data) // Print data
      sck.Emit("Hello", "Hello There!")//Send response to client
    }
  }
}


```

```html

<!-- Add this block before your script -->
<script src="https://unpkg.com/hux-client@1.0.0/hux.minifiy.js"></script>
<script>
  //Initialize hux 
  var hux = new Hux();
  // When hux connection is established, open event will be invoked.
  hux.on('open', () => {
  console.log('Connection established');
    // Listen hello events from server
    hux.on("Hello", () => console.log("GOT MESSAGE"));
    // Send Hello message to server
    hux.emit("Hello", "Hi");

  })
</script>
```


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)