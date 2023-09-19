import {PropsWithChildren, createContext, useCallback, useContext, useEffect, useMemo} from "react";
import {newMessagePayloadSchema, useChatContext, userJoinedPayloadSchema} from "./Chat.context";
import {object, string} from "yup";

export const eventSchema = object({type: string().required(), payload: object().required()});

export const errorPayloadSchema = object({
    message: string(),
    error: string()
});

interface ConnectionContextValue {
    conn: WebSocket;
}

const ConnectionContext = createContext<ConnectionContextValue>({} as ConnectionContextValue);

export function ConnectionContextProvider({children}: PropsWithChildren) {
    const conn = useMemo(() => new WebSocket("ws://localhost:8055/ws"), []);

    const {dispatch} = useChatContext();

    const handleEvent = useCallback(
        async (data: string) => {
            const json = JSON.parse(data);
            const eventData = await eventSchema.validate(json);

            if (eventData) {
                const {type, payload} = eventData;
                switch (type) {
                    case "user-joined": {
                        const validPayload = await userJoinedPayloadSchema.validate(payload);
                        if (validPayload) {
                            return dispatch({type: "user-joined", payload: validPayload});
                        }
                        return;
                    }
                    case "new-message": {
                        const validPayload = await newMessagePayloadSchema.validate(payload);
                        if (validPayload) {
                            return dispatch({type: "new-message", payload: validPayload});
                        }
                        return;
                    }
                    case "error":
                    default:
                    // return await errorPayloadSchema.validate(payload);
                }
            }
        },
        [dispatch]
    );
    useEffect(() => {
        conn.onmessage = async ({data}) => handleEvent(data);
    }, [conn, handleEvent]);

    return <ConnectionContext.Provider value={{conn}}>{children}</ConnectionContext.Provider>;
}

export const useConnection = () => useContext(ConnectionContext);
