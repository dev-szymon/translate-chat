import {useChatContext} from "../context/Chat.context";
import UsersList from "./UsersList";
import {useUserMedia} from "../context/Media.context";
import ChatMessageGroup from "./ChatMessageGroup";
import RecordButton from "./RecordButton/RecordButton";

export default function Chat() {
    const {state} = useChatContext();
    const {currentUser, messages, room, roomUsers} = state;
    const {startRecording, stopRecording, isRecording, isPending} = useUserMedia();

    return (
        <div className="h-full min-h-screen max-h=-screen w-full text-theme-inverted flex-col-reverse items-center lg:flex-row flex gap-2 lg:gap-4 justify-center">
            <div className="mt-16 max-w-[200px] w-full">
                <span className="mb-4">Users: </span>
                {room && <UsersList users={room.users} className="w-full" />}
            </div>
            <div className="w-full h-full flex flex-col gap-4 px-2 py-4 max-w-[500px]">
                <div>
                    <div className="flex gap-2 items-center">
                        <span className="text-xs text-theme-inverted">room name:</span>
                        <span className="text-xl font-medium">{room?.name}</span>
                    </div>
                    <div className="flex gap-2 items-center">
                        <span className="text-xs text-theme-base">room id:</span>
                        <span className="text-xl font-medium">{room?.id}</span>
                    </div>
                </div>
                <div className="flex h-full flex-col gap-4 runded bg-gray-50 min-h-[400px] p-2">
                    {messages.map((message) => (
                        <ChatMessageGroup
                            key={message[0].id}
                            messages={message}
                            sender={roomUsers[message[0]?.sender.id]}
                            isOwn={message[0].sender.id === currentUser?.id}
                        />
                    ))}
                </div>
                <RecordButton
                    onClick={isRecording ? stopRecording : startRecording}
                    variant={isPending ? "pending" : isRecording ? "recording" : "idle"}
                />
            </div>
        </div>
    );
}
