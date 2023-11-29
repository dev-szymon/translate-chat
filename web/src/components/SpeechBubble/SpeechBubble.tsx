import clsx from "clsx";

type SpeechBubbleProps = {
    className?: string;
    isOwn: boolean;
    text: string;
};
const SpeechBubble: React.FC<SpeechBubbleProps> = ({text, isOwn, className}) => {
    return (
        <div
            className={clsx(
                "text-inverted p-2 w-fit max-w-[60%]",
                isOwn
                    ? "self-end bg-theme-primary-light rounded-tl-md rounded-bl-md first-of-type:rounded-tr-md last-of-type:rounded-br-md"
                    : "self-start bg-theme-base-200 rounded-tr-md rounded-br-md first-of-type:rounded-tl-md last-of-type:rounded-bl-md",
                className
            )}
        >
            <p>{text}</p>
        </div>
    );
};

export default SpeechBubble;
