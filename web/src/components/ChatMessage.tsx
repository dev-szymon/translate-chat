import clsx from "clsx";
import {Message, useChatContext} from "../context/Chat.context";
import {supportedLanguages} from "./LanguageSelector";

const heatmapColors = [
    "#FF0000",
    "#FF3300",
    "#FF6600",
    "#FF9900",
    "#FFCC00",
    "#FFFF00",
    "#CCFF00",
    "#99FF00",
    "#66FF00",
    "#33FF00",
    "#00FF00"
];

function decodeHtmlEntity(encodedString: string): string {
    const tempElement = document.createElement("div");
    tempElement.innerHTML = encodedString;
    return tempElement.textContent ?? "";
}

const ChatMessage: React.FC<{message: Message}> = ({message}) => {
    const {
        state: {currentUser, room}
    } = useChatContext();

    const sender = room?.users.find(({id}) => id === message.senderId);
    if (!room || !sender) return null;
    return (
        <div
            className={clsx(
                "flex flex-col max-w-[75%] p-1 text-slate-50",

                message.senderId === currentUser?.id && "items-end self-end pl-2"
            )}
        >
            <div className="flex gap-2">
                <span className="flex-shrink-0 text-xs flex gap-1">
                    <span className="text-slate-400">confidence</span>
                    <span
                        style={{
                            color: heatmapColors[Math.floor(message.confidence * 10)]
                        }}
                    >
                        {Math.floor(message.confidence * 100)}
                    </span>
                </span>
                <span className="text-xs">{`${supportedLanguages[sender.language].icon} ${
                    sender.username
                }`}</span>
            </div>
            <p className="flex-1 text-lg">
                {decodeHtmlEntity(
                    message.translation?.length ? message.translation : message.transcript
                )}
            </p>
        </div>
    );
};
export default ChatMessage;
