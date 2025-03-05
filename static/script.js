const ws = new WebSocket("wss://www.spectrum.cecil-personal.site/ws");

const sentTimestamps = {};

ws.onopen = () => {
    console.log("Connected to WebSocket server");
};

ws.onmessage = (event) => {
    const receivedTime = new Date();
    const msgText = event.data;

    // message sent from client to message received by server
    let latencyText = "";
    if (sentTimestamps[msgText]) {
        const sentTime = sentTimestamps[msgText];
        const latency = receivedTime - sentTime;
        latencyText = ` (latency: ${latency} ms)`;
        delete sentTimestamps[msgText];
    }

    const msgDiv = document.createElement("div");
    msgDiv.textContent = `[message received: ${formatTime(receivedTime)}] ${msgText}${latencyText}`;
    document.getElementById("messages").appendChild(msgDiv);
};

function sendMessage() {
    const input = document.getElementById("messageInput");
    if (input.value.trim() !== "") {
        const sentTime = new Date();
        sentTimestamps[input.value] = sentTime;

        const msgDiv = document.createElement("div");
        msgDiv.textContent = `[message sent: ${formatTime(sentTime)}] ${input.value}`;
        document.getElementById("messages").appendChild(msgDiv);

        ws.send(input.value);
        input.value = "";
    }
}

function formatTime(date) {
    return `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}.${date.getMilliseconds()}`;
}