import clsx from "clsx";
import {useEffect, useRef, useState} from "react";
import LanguageSelector, {Language} from "./LanguageSelector";

type ApiResponse = {
    transcript: string;
    confidence: number;
    translation: string;
    target_language: string;
    source_language: string;
};

type FetchTranslationArgs = {
    file: Blob;
    sourceLanguage: Language["tag"];
    targetLanguage: Language["tag"];
};
const fetchTranslation = async ({file, sourceLanguage, targetLanguage}: FetchTranslationArgs) => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("sourceLanguage", sourceLanguage);
    formData.append("targetLanguage", targetLanguage);

    const postFileRequest = await fetch("http://localhost:8055/transcribe", {
        body: formData,
        method: "POST"
    });
    return (await postFileRequest.json()) as ApiResponse | null;
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
    const [responses, setResponses] = useState<ApiResponse[]>([]);
    const [leftLanguage, setLeftLanguage] = useState<Language | null>(null);
    const [rightLanguage, setRightLanguage] = useState<Language | null>(null);
    const [currentSpeaker, setCurrentSpeaker] = useState<"left" | "right" | null>(null);

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

    function startRecording(speaker: "left" | "right") {
        setCurrentSpeaker(speaker);
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

    useEffect(() => {
        if (isPending && !isRecording) {
            const file = new Blob(audioChunks, {type: "audio/webm"});

            if (file.size > 0 && leftLanguage && rightLanguage) {
                console.log(currentSpeaker);
                fetchTranslation({
                    file,
                    sourceLanguage:
                        currentSpeaker === "left" ? leftLanguage.tag : rightLanguage.tag,
                    targetLanguage: currentSpeaker === "left" ? rightLanguage.tag : leftLanguage.tag
                })
                    .then((response) => {
                        if (response) {
                            setResponses((current) => [...current, response]);
                        }
                    })
                    .finally(() => {
                        setAudioChunks([]);
                        setIsPending(false);
                        setCurrentSpeaker(null);
                    });
            }
        }
    }, [isRecording, isPending, audioChunks, leftLanguage, rightLanguage, currentSpeaker]);

    console.log(responses);
    return (
        <div className="w-full flex justify-center">
            <div className="flex flex-col p-2 gap-2 max-w-xl w-full">
                <div className="flex flex-col justify-end gap-2 h-96 rounded-md border border-gray-300 p-2">
                    {responses.map(({translation, confidence, source_language}) => (
                        <div
                            key={translation}
                            className={clsx(
                                "flex flex-col max-w-[75%] p-1 rounded-sm border border-gray-100",
                                leftLanguage && source_language === leftLanguage.tag
                                    ? "items-start self-start pr-2"
                                    : "items-end self-end pl-2"
                            )}
                        >
                            <p className={"flex-1"}>{decodeHtmlEntity(translation)}</p>
                            <span className="flex-shrink-0 text-xs flex gap-1">
                                <span>confidence</span>
                                <span style={{color: heatmapColors[Math.floor(confidence * 10)]}}>
                                    {Math.floor(confidence * 100)}
                                </span>
                            </span>
                        </div>
                    ))}
                </div>
                <div className="flex justify-between">
                    <div className="flex items-center gap-2">
                        <LanguageSelector
                            currentLanguage={leftLanguage}
                            onLanguageChange={(language) => setLeftLanguage(language)}
                        />
                        <button
                            type="button"
                            onClick={
                                isRecording && currentSpeaker === "left"
                                    ? stopRecording
                                    : () => startRecording("left")
                            }
                            disabled={
                                !leftLanguage || currentSpeaker === "right" || isPending || !stream
                            }
                            className={clsx(
                                "px-4 py-2 border text-white rounded  disabled:bg-gray-600 disabled:border-gray-200",
                                isRecording
                                    ? "bg-red-700 border-red-400"
                                    : "bg-sky-700 border-sky-400"
                            )}
                        >
                            {isRecording && currentSpeaker === "left" ? "stop" : "speak"}
                        </button>
                    </div>

                    <div className="flex items-center gap-2">
                        <button
                            type="button"
                            onClick={
                                isRecording && currentSpeaker === "right"
                                    ? stopRecording
                                    : () => startRecording("right")
                            }
                            disabled={
                                !rightLanguage || currentSpeaker === "left" || isPending || !stream
                            }
                            className={clsx(
                                "px-4 py-2 border text-white rounded  disabled:bg-gray-600 disabled:border-gray-200",
                                isRecording
                                    ? "bg-red-700 border-red-400"
                                    : "bg-sky-700 border-sky-400"
                            )}
                        >
                            {isRecording && currentSpeaker === "right" ? "stop" : "speak"}
                        </button>
                        <LanguageSelector
                            currentLanguage={rightLanguage}
                            onLanguageChange={(language) => setRightLanguage(language)}
                        />
                    </div>
                </div>
            </div>
        </div>
    );
}
