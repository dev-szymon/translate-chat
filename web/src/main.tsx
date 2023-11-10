import {StrictMode} from "react";
import {createRoot} from "react-dom/client";
import App from "./App.tsx";
import "./styles/index.css";
import {ConnectionContextProvider} from "./context/Connection.context.tsx";
import {UserContextProvider} from "./context/Chat.context.tsx";
import {MediaContextProvider} from "./context/Media.context.tsx";

createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <UserContextProvider>
            <ConnectionContextProvider>
                <MediaContextProvider>
                    <App />
                </MediaContextProvider>
            </ConnectionContextProvider>
        </UserContextProvider>
    </StrictMode>
);
