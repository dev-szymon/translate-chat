import Chat from "./components/Chat";
import Home from "./components/Home";
import {useChatContext} from "./context/Chat.context";

function App() {
    const {
        state: {room, currentUser}
    } = useChatContext();

    if (!room || !currentUser) return <Home />;

    return <Chat />;
}

export default App;
