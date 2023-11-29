import clsx from "clsx";
import {Message, User} from "../context/Chat.context";
import {supportedLanguages} from "./LanguageSelector";
import SpeechBubble from "./SpeechBubble/SpeechBubble";

// const heatmapColors = [
//     "#FF0000",
//     "#FF3300",
//     "#FF6600",
//     "#FF9900",
//     "#FFCC00",
//     "#FFFF00",
//     "#CCFF00",
//     "#99FF00",
//     "#66FF00",
//     "#33FF00",
//     "#00FF00"
// ];

function decodeHtmlEntity(encodedString: string): string {
    const tempElement = document.createElement("div");
    tempElement.innerHTML = encodedString;
    return tempElement.textContent ?? "";
}

interface ChatMessageGroupProps {
    sender: User;
    messages: Message[];
    isOwn: boolean;
}
const ChatMessageGroup: React.FC<ChatMessageGroupProps> = ({messages, sender, isOwn}) => {
    const language = supportedLanguages[sender.language];
    return (
        <div className={clsx("flex flex-col w-full p-1 text-theme-inverted")}>
            <div className="flex gap-2">
                <span className="text-xs">{`${language.icon} ${sender.username}`}</span>
            </div>
            <div className="flex w-full flex-col gap-1">
                {messages.map((message) => (
                    <SpeechBubble
                        key={message.id}
                        isOwn={isOwn}
                        text={decodeHtmlEntity(
                            message.translation?.length ? message.translation : message.transcript
                        )}
                    />
                ))}
            </div>
        </div>
    );
};
export default ChatMessageGroup;
