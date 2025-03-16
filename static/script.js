const clientId = localStorage.getItem("clientId") || crypto.randomUUID();
localStorage.setItem("clientId", clientId);

let ws;
const sentTimestamps = {};
const reconnectInterval = 5000;

function connectWebSocket() {
    const wsUrl = `wss://spectrum.cecil-personal.site/ws?clientId=${clientId}`;
    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        console.log(`✅ Connected as ${clientId}`);
        updateStatus("🟢 Connected");

        const initMessage = {
            room: "default",
            sender: clientId,
            content: "",
            type: "join"
        };
        ws.send(JSON.stringify(initMessage));
    };

    ws.onmessage = (event) => {
        const receivedTime = new Date();

        let msgObj;
        try {
            msgObj = JSON.parse(event.data);
        } catch (err) {
            console.error("❌ Failed to parse message:", err);
            return;
        }

        const msgText = `[${msgObj.sender}] ${msgObj.content}`;

        let latencyText = "";
        if (sentTimestamps[msgObj.content]) {
            const sentTime = sentTimestamps[msgObj.content];
            const latency = receivedTime - sentTime;
            latencyText = ` (latency: ${latency} ms)`;
            delete sentTimestamps[msgObj.content];
        }

        displayMessage(`[received: ${formatTime(receivedTime)}] ${msgText}${latencyText}`);
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
    return `${date.getHours()}:${String(date.getMinutes()).padStart(2, "0")}:${String(date.getSeconds()).padStart(2, "0")}.${date.getMilliseconds()}`;
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
