import {PropsWithChildren, createContext, useContext, useEffect, useRef, useState} from "react";
import {useChatContext} from "./Chat.context";
import {useConnection} from "./Connection.context";
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
    const {conn} = useConnection();

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
            const magicByte = new TextEncoder().encode("B");
            const file = new Blob(audioChunks, {type: "audio/webm"});
            const payload = new Blob([magicByte, file]);

            if (file.size > 0 && currentUser) {
                conn.send(payload);
                setAudioChunks([]);
                setIsPending(false);
            }
        }
    }, [isRecording, isPending, audioChunks, conn, currentUser]);

    return (
        <MediaContext.Provider
            value={{startRecording, stopRecording, isRecording, isPending, stream}}
        >
            {children}
        </MediaContext.Provider>
    );
}

export const useUserMedia = () => useContext(MediaContext);
