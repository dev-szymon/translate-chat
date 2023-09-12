import clsx from "clsx";
import {useEffect, useRef, useState} from "react";
import LanguageSelector, {Language, supportedLanguages} from "./LanguageSelector";
import {useConnection} from "../context/WebSocket.context";
import {useChatContext} from "../context/Chat.context";

type FetchTranslationArgs = {
    file: Blob;
    sourceLanguage: Language["tag"];
    userId: string;
};
const postAudioFile = async ({file, sourceLanguage, userId}: FetchTranslationArgs) => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("sourceLanguage", sourceLanguage);
    formData.append("userId", userId);

    const postFileRequest = await fetch("http://localhost:8055/translate-file", {
        body: formData,
        method: "POST"
    });
    return await postFileRequest.json();
};

const heatmapColors = [
    "#FF0000",
    "#FF3300",
    "#FF6600",
    "#FF9900",
    "#FFCC00",
    "#FFFF00",
    "#CCFF00",
    "#99FF00",
    "#66FF00",
    "#33FF00",
    "#00FF00"
];

export default function Chat() {
    const [isRecording, setIsRecording] = useState<boolean>(false);
    const mediaRecorder = useRef<MediaRecorder | null>(null);
    const [stream, setStream] = useState<MediaStream | null>(null);
    const [audioChunks, setAudioChunks] = useState<Blob[]>([]);
    const [isPending, setIsPending] = useState<boolean>(false);
    const {state} = useChatContext();

    const {currentUser, messages, room} = state;
    function decodeHtmlEntity(encodedString: string): string {
        const tempElement = document.createElement("div");
        tempElement.innerHTML = encodedString;
        return tempElement.textContent ?? "";
    }
    useEffect(() => {
        if (navigator && !stream) {
            navigator.mediaDevices
                .getUserMedia({
                    audio: true
                })
                .then((streamData) => setStream(streamData));
        }
    }, [stream]);

    function startRecording() {
        setIsRecording(true);
        if (stream) {
            if (!mediaRecorder.current) {
                mediaRecorder.current = new MediaRecorder(stream, {mimeType: "audio/webm"});
                mediaRecorder.current.ondataavailable = ({data}) => {
                    if (data?.size > 0) {
                        setAudioChunks((current) => [...current, data]);
                    }
                };
            }

            mediaRecorder.current.start();
        }
    }

    function stopRecording() {
        if (stream && mediaRecorder.current && mediaRecorder.current.state !== "inactive") {
            mediaRecorder.current.stop();
            // stream.getTracks().forEach((track) => track.stop());
        }
        setIsPending(true);
        setIsRecording(false);
    }

    const {conn} = useConnection();
    useEffect(() => {
        if (isPending && !isRecording && currentUser) {
            const file = new Blob(audioChunks, {type: "audio/webm"});

            if (file.size > 0 && currentUser) {
                postAudioFile({
                    file,
                    sourceLanguage: currentUser.language,
                    userId: currentUser.id
                })
                    .then((response) => {
                        if (response) {
                            return;
                        }
                    })
                    .finally(() => {
                        setAudioChunks([]);
                        setIsPending(false);
                    });
            }
        }
    }, [isRecording, isPending, audioChunks, currentUser, conn]);

    return (
        <div className="w-full flex flex-col items-center">
            <h1 className="text-xl font-medium">{room?.name}</h1>
            <div className="flex flex-col p-2 gap-2 max-w-xl w-full">
                <div className="flex flex-col justify-end gap-2 h-96 rounded-md border border-gray-300 p-2">
                    {messages.map((message) => (
                        <div
                            key={message.transcript}
                            className={clsx(
                                "flex flex-col max-w-[75%] p-1 rounded-sm border border-gray-100",
                                "items-end self-end pl-2"
                            )}
                        >
                            <p className={"flex-1"}>
                                {decodeHtmlEntity(
                                    message.translation?.length
                                        ? message.translation
                                        : message.transcript
                                )}
                            </p>
                            <span className="flex-shrink-0 text-xs flex gap-1">
                                <span>confidence</span>
                                <span
                                    style={{
                                        color: heatmapColors[Math.floor(message.confidence * 10)]
                                    }}
                                >
                                    {Math.floor(message.confidence * 100)}
                                </span>
                            </span>
                        </div>
                    ))}
                </div>
                <div className="flex justify-between">
                    <div className="flex items-center gap-2">
                        <LanguageSelector
                            currentLanguage={
                                currentUser ? supportedLanguages[currentUser.language] : null
                            }
                            onLanguageChange={(language) => console.log(language)} //TODO
                        />
                        <button
                            type="button"
                            onClick={isRecording ? stopRecording : () => startRecording()}
                            disabled={!currentUser?.language || isPending || !stream}
                            className={clsx(
                                "px-4 py-2 border text-white rounded  disabled:bg-gray-600 disabled:border-gray-200",
                                isRecording
                                    ? "bg-red-700 border-red-400"
                                    : "bg-sky-700 border-sky-400"
                            )}
                        >
                            {isRecording ? "stop" : "speak"}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
