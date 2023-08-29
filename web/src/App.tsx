import Chat from "./components/Chat";
import Home from "./components/Home";
import {useConnection} from "./context/WebSocket.context";

function App() {
    const {roomId} = useConnection();
    return <div>{roomId ? <Chat /> : <Home />}</div>;
}

export default App;
