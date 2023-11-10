import clsx from "clsx";
import {Message, User} from "../context/Chat.context";
import {supportedLanguages} from "./LanguageSelector";
import SpeechBubble from "./SpeechBubble/SpeechBubble";

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

interface ChatMessageProps {
    sender: User;
    message: Message;
    isOwn: boolean;
}
const ChatMessage: React.FC<ChatMessageProps> = ({message, sender, isOwn}) => {
    const confidenceColor = heatmapColors[Math.floor(message.confidence * 10)];
    const language = supportedLanguages[sender.language];
    return (
        <div
            className={clsx(
                "flex flex-col max-w-[75%] p-1 text-theme-base",
                isOwn && "items-end self-end pl-2"
            )}
        >
            <div className="flex gap-2">
                <span className="flex-shrink-0 text-xs flex gap-1">
                    <span className="text-theme-base">confidence</span>
                    <span
                        style={{
                            color: confidenceColor
                        }}
                    >
                        {Math.floor(message.confidence * 100)}
                    </span>
                </span>
                <span className="text-xs">{`${language.icon} ${sender.username}`}</span>
            </div>
            <SpeechBubble
                text={decodeHtmlEntity(
                    message.translation?.length ? message.translation : message.transcript
                )}
            />
        </div>
    );
};
export default ChatMessage;
