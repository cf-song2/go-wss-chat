const ws = new WebSocket("wss://www.spectrum.cecil-personal.site/ws");

ws.onopen = () => {
    console.log("Connected to WebSocket server");
};

ws.onmessage = (event) => {
    const msgDiv = document.createElement("div");
    msgDiv.textContent = event.data;
    document.getElementById("messages").appendChild(msgDiv);
};

function sendMessage() {
    const input = document.getElementById("messageInput");
    if (input.value.trim() !== "") {
        ws.send(input.value);
        input.value = "";
    }
}

