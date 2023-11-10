import clsx from "clsx";
import {useChatContext} from "../context/Chat.context";
import UsersList from "./UsersList";
import {useUserMedia} from "../context/Media.context";
import ChatMessage from "./ChatMessage";

export default function Chat() {
    const {state} = useChatContext();
    const {currentUser, messages, room, roomUsers} = state;
    const {startRecording, stopRecording, isRecording, isPending} = useUserMedia();

    return (
        <div className="bg-slate-900 h-full min-h-screen w-full flex-col-reverse items-center lg:flex-row flex gap-2 lg:gap-4 justify-center">
            <div className="mt-16 max-w-[200px] text-slate-400 w-full">
                <span className="mb-4">Users: </span>
                {room && <UsersList users={room.users} className="w-full" />}
            </div>
            <div className="w-full h-full flex flex-col gap-4 px-2 py-4 max-w-[500px]">
                <div>
                    <div className="flex gap-2 items-center">
                        <span className="text-xs text-slate-300">room name:</span>
                        <span className="text-xl font-medium text-slate-50">{room?.name}</span>
                    </div>
                    <div className="flex gap-2 items-center">
                        <span className="text-xs text-slate-300">room id:</span>
                        <span className="text-xl font-medium text-slate-50">{room?.id}</span>
                    </div>
                </div>
                <div className="flex h-full flex-col gap-4 runded bg-slate-800 min-h-[400px] p-2">
                    {messages.map((message) => (
                        <ChatMessage
                            key={message.id}
                            message={message}
                            sender={roomUsers[message.sender.id]}
                            isOwn={message.sender.id === currentUser?.id}
                        />
                    ))}
                </div>
                <button
                    type="button"
                    onClick={isRecording ? stopRecording : startRecording}
                    disabled={isPending}
                    className={clsx(
                        "px-4 py-2 border text-white rounded w-full disabled:bg-gray-600 disabled:border-gray-200",
                        isRecording ? "bg-red-800 border-red-400" : "bg-sky-700 border-sky-400"
                    )}
                >
                    {isPending ? "translating..." : isRecording ? "stop" : "speak"}
                </button>
            </div>
        </div>
    );
}
