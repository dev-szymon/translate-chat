import clsx from "clsx";
import {useEffect, useRef, useState} from "react";

interface ApiResponse {
    transcript: string;
    translation: string;
}

const handleTranslate = async (file: Blob) => {
    const formData = new FormData();
    formData.append("file", file);

    const postFileRequest = await fetch("http://localhost:8055/transcribe", {
        body: formData,
        method: "POST"
    });
    return (await postFileRequest.json()) as ApiResponse | null;
};

export default function Chat() {
    const [isRecording, setIsRecording] = useState<boolean>(false);
    const mediaRecorder = useRef<MediaRecorder | null>(null);
    const [stream, setStream] = useState<MediaStream | null>(null);
    const [audioChunks, setAudioChunks] = useState<Blob[]>([]);
    const [isPending, setIsPending] = useState<boolean>(false);
    const [responses, setResponses] = useState<ApiResponse[]>([]);

    useEffect(() => {
        if (navigator && !stream) {
            navigator.mediaDevices
                .getUserMedia({
                    audio: true
                })
                .then((streamData) => setStream(streamData));
        }
    }, [stream]);

    async function startRecording() {
        if (stream) {
            setIsRecording(true);
            mediaRecorder.current = new MediaRecorder(stream, {mimeType: "audio/webm"});

            mediaRecorder.current.start();

            mediaRecorder.current.ondataavailable = ({data}) => {
                if (data?.size > 0) {
                    setAudioChunks((current) => [...current, data]);
                }
            };
        }
    }

    async function stopRecording() {
        if (mediaRecorder.current && mediaRecorder.current.state !== "inactive" && stream) {
            mediaRecorder.current.stop();
            stream.getTracks().forEach((track) => track.stop());
        }
        setIsPending(true);
        setIsRecording(false);
    }

    useEffect(() => {
        if (isPending && !isRecording) {
            const file = new Blob(audioChunks, {type: "audio/webm"});

            handleTranslate(file)
                .then((response) => {
                    if (response) {
                        setResponses((current) => [...current, response]);
                    }
                })
                .finally(() => {
                    setAudioChunks([]);
                    setIsPending(false);
                });
        }
    }, [isRecording, isPending, audioChunks]);

    return (
        <div className="flex flex-col">
            <div className="flex flex-col-reverse justify-start gap-2 h-96">
                {responses.map(({translation}) => (
                    <p key={translation}>{translation}</p>
                ))}
            </div>
            <button
                type="button"
                onClick={isRecording ? stopRecording : startRecording}
                disabled={isPending || !stream}
                className={clsx(
                    "px-4 py-2 border text-white rounded  disabled:bg-gray-600 disabled:border-gray-200",
                    isRecording ? "bg-red-700 border-red-400" : "bg-sky-700 border-sky-400"
                )}
            >
                {isRecording ? "stop" : "record"}
            </button>
        </div>
    );
}
