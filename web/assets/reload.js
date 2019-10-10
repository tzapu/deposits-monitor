function connect() {
    let socket = new WebSocket("ws://localhost:8080/ws");

    socket.onclose = function(event) {
        console.log("Websocket connection closed or unable to connect; " +
            "starting reconnect timeout");

        // Allow the last socket to be cleaned up.
        socket = null;

        // Set an interval to continue trying to reconnect
        // periodically until we succeed.
        setTimeout(function() {
            connect();
        }, 5000)
    }

    socket.onmessage = function(event) {
        const data = JSON.parse(event.data);
        switch(data.type) {
            case "build_complete":
                socket.close(1000, "Reloading page after receiving build_complete");

                console.log("Reloading page after receiving build_complete");
                location.reload(true);

                break;

            default:
                console.log(`Don't know how to handle type '${data.type}'`);
        }
    }
}

connect();