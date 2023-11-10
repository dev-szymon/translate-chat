import LanguageSelector, {Language} from "./LanguageSelector";
import {useConnection} from "../context/Connection.context";
import {useState} from "react";

type JoinRoomMessage = {
    type: "join-room-event";
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
            type: "join-room-event",
            payload: {
                username,
                roomId,
                language: language.tag
            }
        };
        conn.send(JSON.stringify(message));
    };

    return (
        <div className="w-full bg-slate-900 text-slate-50 min-h-[100vh] flex items-center justify-center p-4">
            <form onSubmit={handleSubmit} className="bg-slate-700 rounded p-4 max-w-xl w-full">
                <div className="flex gap-2">
                    <div className="flex-1">
                        <div className="text-base gap-2 flex mb-2 flex-col">
                            <label htmlFor="username">Username</label>
                            <input
                                className="rounded border border-slate-400 bg-slate-600 px-4 py-2"
                                placeholder="John Doe"
                                type="text"
                                id="username"
                                name="username"
                            />
                        </div>
                        <div className="text-base gap-2 mb-2 flex flex-col">
                            <label htmlFor="roomId">{`Room id or name (optional)`}</label>
                            <input
                                className="rounded border border-slate-400 bg-slate-600 px-4 py-2"
                                placeholder="room-id-or-name"
                                type="text"
                                id="roomId"
                                name="roomId"
                            />
                        </div>
                    </div>

                    <LanguageSelector
                        className="mt-8 min-w-[200px]"
                        currentLanguage={language}
                        onLanguageChange={(language) => setLanguage(language)}
                    />
                </div>
                <button
                    type="submit"
                    disabled={!language}
                    className="flex px-4 py-2 cursor-pointer mt-4 rounded border-none bg-sky-700 text-white"
                >
                    join room
                </button>
            </form>
        </div>
    );
}
