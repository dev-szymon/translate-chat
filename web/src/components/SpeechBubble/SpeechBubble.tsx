import clsx from "clsx";

type SpeechBubbleProps = {
    className?: string;
    text: string;
};
const SpeechBubble: React.FC<SpeechBubbleProps> = ({text, className}) => {
    return (
        <div className={clsx("bg-theme-primary text-theme-base p-2 rounded-md", className)}>
            <p>{text}</p>
        </div>
    );
};

export default SpeechBubble;
