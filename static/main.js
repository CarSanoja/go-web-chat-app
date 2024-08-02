const ws = new WebSocket("ws://localhost:8080/ws");

ws.onopen = function() {
    const username = prompt("Enter your username");
    ws.send(JSON.stringify({ type: "username", data: username }));
};

ws.onmessage = function(event) {
    const msg = JSON.parse(event.data);
    const output = document.getElementById("output");

    if (msg.type === "text") {
        output.innerHTML += `<p><strong>${msg.username}:</strong> ${msg.text}</p>`;
    } else if (msg.type === "image") {
        output.innerHTML += `<p><strong>${msg.username}:</strong><br><img src="${msg.data}" /></p>`;
    }

    output.scrollTop = output.scrollHeight;
};

function sendMessage() {
    const messageInput = document.getElementById("message");
    const usernameInput = document.getElementById("username");
    const imageInput = document.getElementById("imageUpload");

    if (imageInput.files.length > 0) {
        const reader = new FileReader();
        reader.onload = function(e) {
            const imgData = e.target.result;
            const msg = {
                type: "image",
                username: usernameInput.value,
                data: imgData
            };
            ws.send(JSON.stringify(msg));
        };
        reader.readAsDataURL(imageInput.files[0]);
    } else {
        const msg = {
            type: "text",
            username: usernameInput.value,
            text: messageInput.value
        };
        ws.send(JSON.stringify(msg));
        messageInput.value = "";
    }
}
