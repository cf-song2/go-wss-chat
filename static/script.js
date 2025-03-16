const wsUrl = "wss://www.spectrum.cecil-personal.site/ws";
let ws;
const sentTimestamps = {};
let reconnectInterval = 5000;

function connectWebSocket() {
    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        console.log("✅ Connected to WebSocket server");
        updateStatus("🟢 Connected");
    };

    ws.onmessage = (event) => {
        const receivedTime = new Date();
        const msgText = event.data;

        let latencyText = "";
        if (sentTimestamps[msgText]) {
            const sentTime = sentTimestamps[msgText];
            const latency = receivedTime - sentTime;
            latencyText = ` (latency: ${latency} ms)`;
            delete sentTimestamps[msgText];
        }

        displayMessage(`[message received: ${formatTime(receivedTime)}] ${msgText}${latencyText}`);
    };

    ws.onclose = () => {
        console.warn("⚠️ WebSocket Disconnected. Attempting to reconnect...");
        updateStatus("🔴 Disconnected. Reconnecting...");
        setTimeout(connectWebSocket, reconnectInterval);
    };

    ws.onerror = (error) => {
        console.error("WebSocket Error:", error);
        ws.close();
    };
}

function sendMessage() {
    const input = document.getElementById("messageInput");
    if (input.value.trim() !== "" && ws.readyState === WebSocket.OPEN) {
        const sentTime = new Date();
        sentTimestamps[input.value] = sentTime;

        displayMessage(`[message sent: ${formatTime(sentTime)}] ${input.value}`);
        ws.send(input.value);
        input.value = "";
    } else {
        displayMessage("⚠️ Cannot send message, WebSocket is disconnected.");
    }
}

function formatTime(date) {
    return `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}.${date.getMilliseconds()}`;
}

function displayMessage(message) {
    const msgDiv = document.createElement("div");
    msgDiv.textContent = message;
    document.getElementById("messages").appendChild(msgDiv);
}

function updateStatus(status) {
    document.getElementById("status").textContent = status;
}

connectWebSocket();
