let mediaRecorder;
let mediaStream;
let chunks = [];

async function startRecording() {
    try {
        const stream = await navigator.mediaDevices.getUserMedia({
            audio: true,
        });
        mediaStream = stream;
        mediaRecorder = new MediaRecorder(stream, { mimeType: "audio/webm" });

        mediaRecorder.addEventListener("dataavailable", (e) => {
            chunks.push(e.data);
        });

        mediaRecorder.start();
    } catch (err) {
        console.error("Error starting recording: ", err);
    }
}

document
    .getElementById("record-button")
    .addEventListener("click", startRecording);
