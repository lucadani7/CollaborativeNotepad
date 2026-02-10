const docID = "demo1";
const protocol = window.location.protocol === "https:" ? "wss://" : "ws://";
const ws = new WebSocket(protocol + window.location.host + "/ws?doc=" + docID);
const editor = document.getElementById("editor");
const status = document.getElementById("status");
let lastContent = "";

ws.onopen = () => {
    status.textContent = "Connected!";
    status.className = "connected";
    editor.disabled = false;
};

ws.onclose = () => {
    status.textContent = "Disconnected!";
    status.className = "disconnected";
    editor.disabled = true;
};

ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    const start = editor.selectionStart;
    const scrollTop = editor.scrollTop;
    if (msg.type === "INIT") {
        editor.value = msg.content;
        lastContent = msg.content;
    } else if (msg.type === "INSERT") {
        const text = editor.value;
        const newVal = text.substring(0, msg.pos) + msg.content + text.substring(msg.pos);
        editor.value = newVal;
        lastContent = newVal;
        let newCursor = start;
        if (msg.pos <= start) {
            newCursor += msg.content.length;
        }
        editor.setSelectionRange(newCursor, newCursor);
    } else if (msg.type === "DELETE") {
        const text = editor.value;
        const newVal = text.substring(0, msg.pos) + text.substring(msg.pos + msg.len);
        editor.value = newVal;
        lastContent = newVal;
        let newCursor = start;
        if (msg.pos < start) {
            const amountDeletedBeforeCursor = Math.min(msg.len, start - msg.pos);
            newCursor -= amountDeletedBeforeCursor;
        }
        editor.setSelectionRange(newCursor, newCursor);
    }
    editor.scrollTop = scrollTop;
};


function getDiff(oldText, newText) {
    let start = 0;
    while (start < oldText.length && start < newText.length && oldText[start] === newText[start]) {
        start++;
    }
    let endOld = oldText.length;
    let endNew = newText.length;
    while (endOld > start && endNew > start && oldText[endOld - 1] === newText[endNew - 1]) {
        endOld--;
        endNew--;
    }
    const lenInsert = endNew - start;
    const lenDelete = endOld - start;
    if (lenInsert > 0) {
        return {
            type: "INSERT",
            pos: start,
            content: newText.substring(start, endNew)
        };
    }
    if (lenDelete > 0) {
        return {
            type: "DELETE",
            pos: start,
            len: lenDelete
        };
    }
    return null;
}

editor.addEventListener("input", () => {
    const currentContent = editor.value;
    const change = getDiff(lastContent, currentContent);
    if (change) {
        ws.send(JSON.stringify({
            type: change.type,
            doc_id: docID,
            pos: change.pos,
            content: change.content || "",
            len: change.len || 0
        }));
        lastContent = currentContent;
    }
});