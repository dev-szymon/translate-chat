import clsx from "clsx";
import {MouseEventHandler} from "react";

type RecordButtonProps = {
    onClick: MouseEventHandler<HTMLButtonElement>;
    disabled?: boolean;
    className?: string;
    variant: "idle" | "recording" | "pending";
};

const recordButtonText: Record<RecordButtonProps["variant"], string> = {
    idle: "speak",
    pending: "translating...",
    recording: "stop"
};

const RecordButton: React.FC<RecordButtonProps> = ({className, onClick, disabled, variant}) => {
    return (
        <button
            type="button"
            onClick={onClick}
            disabled={disabled || variant == "pending"}
            className={clsx(
                "px-4 py-2 border text-white rounded w-full disabled:bg-gray-600 disabled:border-gray-200",
                variant === "recording" ? "bg-red-800 border-red-400" : "bg-sky-700 border-sky-400",
                className
            )}
        >
            {recordButtonText[variant]}
        </button>
    );
};
export default RecordButton;
