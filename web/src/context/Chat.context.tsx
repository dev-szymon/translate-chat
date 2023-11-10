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
    sender: userSchema.required()
});
export type Message = {
    id: string;
    transcript: string;
    translation: string | null;
    confidence: number;
    sender: User;
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
    type: "user-joined-event";
    payload: {newUser: User; room: Room};
};

export const newMessagePayloadSchema = object({
    message: messageSchema.required()
});

type NewMessageAction = {
    type: "new-message-event";
    payload: {message: Message};
};

export type JoinRoomMessage = {
    type: "join-room-event";
    payload: {
        username: string;
        language: string;
        roomId?: string;
    };
};

export type ChatState = {
    currentUser: User | null;
    room: Room | null;
    roomUsers: Record<User["id"], User>;
    messages: Array<Message>;
};
export type Action = UserJoinedAction | NewMessageAction;

function chatReducer(state: ChatState, action: Action): ChatState {
    switch (action.type) {
        case "user-joined-event": {
            const roomUsers = action.payload.room.users.reduce(
                (users: ChatState["roomUsers"], curr: User) => {
                    return {...users, [curr.id]: curr};
                },
                {}
            );
            const room = action.payload.room;

            return {
                ...state,
                currentUser: !state.currentUser ? action.payload.newUser : state.currentUser,
                room,
                roomUsers
            };
        }
        case "new-message-event": {
            return {...state, messages: [...state.messages, action.payload.message]};
        }
        default:
            return state;
    }
}
const initialChatState: ChatState = {
    currentUser: null,
    room: null,
    roomUsers: {},
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
