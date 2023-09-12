import {PropsWithChildren, createContext, useContext, useState} from "react";
import {ChatMessage, User} from "../service/ChatMessage";

type Room = {
    id: string;
    name: string;
    users: User[];
};
interface UserContext {
    user: User | null;
    setUser: (user: User) => void;
    room: Room | null;
    setRoom: (room: Room) => void;
    messages: ChatMessage[];
    addMessage: (message: ChatMessage) => void;
}

const UserContext = createContext<UserContext>({} as UserContext);

export function UserContextProvider({children}: PropsWithChildren) {
    const [user, setUser] = useState<User | null>(null);
    const [room, setRoom] = useState<Room | null>(null);
    const [messages, setMessages] = useState<ChatMessage[]>([]);

    return (
        <UserContext.Provider
            value={{
                user,
                setUser,
                room,
                setRoom,
                messages,
                addMessage: (message) =>
                    setMessages((current) => {
                        console.log("message", message);
                        return [...current, message];
                    })
            }}
        >
            {children}
        </UserContext.Provider>
    );
}

export const useUser = () => useContext(UserContext);
