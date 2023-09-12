import Chat from "./components/Chat";
import Home from "./components/Home";
import {useUser} from "./context/User.context";

function App() {
    const {room} = useUser();
    return <>{room ? <Chat /> : <Home />}</>;
}

export default App;
