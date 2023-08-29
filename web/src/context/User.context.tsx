import {PropsWithChildren, createContext, useContext, useState} from "react";
import {User} from "../service/ChatMessage";

interface UserContext {
    user: User | null;
    setUser: (user: User) => void;
}

const UserContext = createContext<UserContext>({} as UserContext);

export function UserContextProvider({children}: PropsWithChildren) {
    const [user, setUser] = useState<User | null>(null);

    return <UserContext.Provider value={{user, setUser}}>{children}</UserContext.Provider>;
}

export const useUser = () => useContext(UserContext);
