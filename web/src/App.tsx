import Chat from "./components/Chat";
import Home from "./components/Home";
import {useChatContext} from "./context/Chat.context";

function App() {
    const {state} = useChatContext();
    return <>{state.room ? <Chat /> : <Home />}</>;
}

export default App;
