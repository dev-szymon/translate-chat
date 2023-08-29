import {PropsWithChildren, createContext, useContext, useEffect, useMemo, useState} from "react";

interface ConnectionContextValue {
    conn: WebSocket;
    roomId: string | null;
    setRoomId: (roomId: string) => void;
    messages: string[];
}

const ConnectionContext = createContext<ConnectionContextValue>({} as ConnectionContextValue);

export function ConnectionContextProvider({children}: PropsWithChildren) {
    const [roomId, setRoomId] = useState<string | null>(null);
    const [messages, setMessages] = useState<string[]>([]);
    const conn = useMemo(() => new WebSocket("ws://localhost:8055/ws"), []);

    useEffect(() => {
        conn.onmessage = ({data}) => {
            console.log(data);
            const roomIdPrefix = /^roomId:/;
            if (typeof data === "string") {
                if (roomIdPrefix.test(data)) {
                    setRoomId(data.replace(roomIdPrefix, ""));
                } else {
                    setMessages((current) => [...current, data]);
                }
            }
        };
    }, [conn]);

    return (
        <ConnectionContext.Provider value={{conn, roomId, setRoomId, messages}}>
            {children}
        </ConnectionContext.Provider>
    );
}

export const useConnection = () => useContext(ConnectionContext);
