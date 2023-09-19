import Chat from "./components/Chat";
import Home from "./components/Home";
import {useChatContext} from "./context/Chat.context";

function App() {
    const {
        state: {room}
    } = useChatContext();

    if (!room) return <Home />;

    return <Chat />;
}

export default App;
