import LanguageSelector, {Language} from "./LanguageSelector";
import {useConnection} from "../context/WebSocket.context";
import {useState} from "react";

type JoinRoomMessage = {
    type: "join-room";
    payload: {
        username: string;
        language: string;
        roomId?: string;
    };
};

export default function Home() {
    const {conn} = useConnection();
    const [language, setLanguage] = useState<Language | null>(null);

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const formData = new FormData(event.target as HTMLFormElement);
        const username = formData.get("username");
        if (typeof username !== "string" || !language) {
            return;
        }
        const roomIdValue = formData.get("roomId") ?? undefined;
        const roomId = typeof roomIdValue === "string" ? roomIdValue : undefined;

        const message: JoinRoomMessage = {
            type: "join-room",
            payload: {
                username,
                roomId,
                language: language.tag
            }
        };
        conn.send(JSON.stringify(message));
    };

    return (
        <div>
            <form onSubmit={handleSubmit}>
                <div>
                    <LanguageSelector
                        currentLanguage={language}
                        onLanguageChange={(language) => setLanguage(language)}
                    />
                    <input className="rounded border border-gray-300" type="text" name="username" />
                </div>
                <input className="rounded border border-gray-300" type="text" name="roomId" />
                <button
                    type="submit"
                    disabled={!language}
                    className="flex px-4 py-2 rounded bg-sky-800 text-white"
                >
                    join room
                </button>
            </form>
        </div>
    );
}
