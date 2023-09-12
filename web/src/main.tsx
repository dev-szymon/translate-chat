import {StrictMode} from "react";
import {createRoot} from "react-dom/client";
import App from "./App.tsx";
import "./styles/index.css";
import {ConnectionContextProvider} from "./context/WebSocket.context.tsx";
import {UserContextProvider} from "./context/Chat.context.tsx";

createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <UserContextProvider>
            <ConnectionContextProvider>
                <App />
            </ConnectionContextProvider>
        </UserContextProvider>
    </StrictMode>
);
