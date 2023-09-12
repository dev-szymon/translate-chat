import {Dispatch, PropsWithChildren, createContext, useContext, useReducer} from "react";
import {array, number, object, string} from "yup";

const userSchema = object({
    id: string().required(),
    username: string().required(),
    language: string().required()
});
export type User = {
    id: string;
    username: string;
    language: string;
};

const messageSchema = object({
    id: string().required(),
    transcript: string().required(),
    confidence: number().required(),
    translation: string().defined().nullable(),
    senderId: string().required()
});
export type Message = {
    id: string;
    transcript: string;
    translation: string | null;
    confidence: number;
    senderId: string;
};

const roomSchema = object({
    id: string().required(),
    name: string().required(),
    users: array().required().of(userSchema)
});
export type Room = {
    id: string;
    name: string;
    users: Array<User>;
};

export const userJoinedPayloadSchema = object({
    newUser: userSchema.required(),
    room: roomSchema.required()
});
type UserJoinedAction = {
    type: "user-joined";
    payload: {newUser: User; room: Room};
};

export const newMessagePayloadSchema = object({
    message: messageSchema.required()
});
type NewMessageAction = {
    type: "new-message";
    payload: {message: Message};
};

export type ChatState = {
    currentUser: User | null;
    room: Room | null;
    messages: Array<Message>;
};
export type Action = UserJoinedAction | NewMessageAction;

function chatReducer(state: ChatState, action: Action): ChatState {
    switch (action.type) {
        case "user-joined": {
            if (!state.currentUser && !state.room) {
                return {...state, currentUser: action.payload.newUser, room: action.payload.room};
            }
            return {...state, room: action.payload.room};
        }
        case "new-message": {
            return {...state, messages: [...state.messages, action.payload.message]};
        }
        default:
            return state;
    }
}
const initialChatState: ChatState = {
    currentUser: null,
    room: null,
    messages: []
};
type ChatContextValue = {state: ChatState; dispatch: Dispatch<Action>};
const ChatContext = createContext<ChatContextValue>({} as ChatContextValue);

export function UserContextProvider({children}: PropsWithChildren) {
    const [state, dispatch] = useReducer(chatReducer, initialChatState);

    return (
        <ChatContext.Provider
            value={{
                state,
                dispatch
            }}
        >
            {children}
        </ChatContext.Provider>
    );
}

export const useChatContext = () => useContext(ChatContext);
