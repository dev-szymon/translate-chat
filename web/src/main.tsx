import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./styles/index.css";
import {ConnectionContextProvider} from "./context/WebSocket.context.tsx";
import {UserContextProvider} from "./context/User.context.tsx";

ReactDOM.createRoot(document.getElementById("root")!).render(
    <React.StrictMode>
        <ConnectionContextProvider>
            <UserContextProvider>
                <App />
            </UserContextProvider>
        </ConnectionContextProvider>
    </React.StrictMode>
);
