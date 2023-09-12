import {PropsWithChildren, createContext, useCallback, useContext, useEffect, useMemo} from "react";
import {
    User,
    eventSchema,
    newMessagePayloadSchema,
    userJoinedRoomPayloadSchema
} from "../service/ChatMessage";
import {useUser} from "./User.context";
import {supportedLanguages} from "../components/LanguageSelector";

interface ConnectionContextValue {
    conn: WebSocket;
}

const ConnectionContext = createContext<ConnectionContextValue>({} as ConnectionContextValue);

const useHandleEvent = () => {
    const {setUser, setRoom, addMessage, room} = useUser();

    const handleEvent = useCallback(
        async (data: string) => {
            const json = JSON.parse(data);
            const eventData = await eventSchema.validate(json);

            if (eventData) {
                const {type, payload} = eventData;
                switch (type) {
                    case "user-joined": {
                        const validPayload = await userJoinedRoomPayloadSchema.validate(payload);
                        const language = supportedLanguages.find(
                            ({tag}) => tag === validPayload?.language
                        );
                        if (payload && language) {
                            setRoom({
                                id: validPayload.roomId,
                                name: validPayload.roomName,
                                users: validPayload.users
                                    .map((user) => ({
                                        ...user,
                                        language: supportedLanguages.find(
                                            ({tag}) => tag === user.language
                                        )
                                    }))
                                    .filter((user): user is User => !!user.language)
                            });
                            setUser({
                                id: validPayload.userId,
                                language,
                                username: validPayload.username
                            });
                        }
                        return;
                    }
                    case "translated-message": {
                        const validPayload = await newMessagePayloadSchema.validate(payload);
                        const sender = room?.users.find((user) => user.id === validPayload.userId);
                        if (sender) {
                            addMessage({
                                transcript: validPayload.transcript,
                                confidence: validPayload.confidence,
                                translation: validPayload.translation ?? null,
                                sender
                            });
                        }
                        return;
                    }
                    case "error":
                    default:
                    // return await errorPayloadSchema.validate(payload);
                }
            }
        },
        [setUser, setRoom, addMessage, room?.users]
    );
    return {handleEvent};
};

export function ConnectionContextProvider({children}: PropsWithChildren) {
    const conn = useMemo(() => new WebSocket("ws://localhost:8055/ws"), []);

    const {handleEvent} = useHandleEvent();
    useEffect(() => {
        conn.onmessage = async ({data}) => handleEvent(data);
    }, [conn, handleEvent]);

    return <ConnectionContext.Provider value={{conn}}>{children}</ConnectionContext.Provider>;
}

export const useConnection = () => useContext(ConnectionContext);
