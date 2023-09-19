import {PropsWithChildren, createContext, useContext, useEffect, useRef, useState} from "react";
import {Language} from "../components/LanguageSelector";
import {useChatContext} from "./Chat.context";

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
interface MediaContextValue {
    isRecording: boolean;
    isPending: boolean;
    stream: MediaStream | null;

    startRecording: () => void;
    stopRecording: () => void;
}

const MediaContext = createContext<MediaContextValue>({} as MediaContextValue);

export function MediaContextProvider({children}: PropsWithChildren) {
    const [isRecording, setIsRecording] = useState<MediaContextValue["isRecording"]>(false);
    const mediaRecorder = useRef<MediaRecorder | null>(null);
    const [stream, setStream] = useState<MediaContextValue["stream"]>(null);
    const [audioChunks, setAudioChunks] = useState<Blob[]>([]);
    const [isPending, setIsPending] = useState<MediaContextValue["isPending"]>(false);
    const {
        state: {currentUser}
    } = useChatContext();

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
    }, [isRecording, isPending, audioChunks, currentUser]);

    return (
        <MediaContext.Provider
            value={{startRecording, stopRecording, isRecording, isPending, stream}}
        >
            {children}
        </MediaContext.Provider>
    );
}

export const useUserMedia = () => useContext(MediaContext);
