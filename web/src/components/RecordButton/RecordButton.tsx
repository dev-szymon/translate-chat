import {MicrophoneIcon} from "@heroicons/react/24/outline";
import clsx from "clsx";
import {MouseEventHandler} from "react";

type RecordButtonProps = {
    onClick: MouseEventHandler<HTMLButtonElement>;
    disabled?: boolean;
    className?: string;
    variant: "idle" | "recording" | "pending";
};

const RecordButton: React.FC<RecordButtonProps> = ({className, onClick, disabled, variant}) => {
    return (
        <button
            type="button"
            onClick={onClick}
            disabled={disabled || variant == "pending"}
            className={clsx(
                "flex min-w-[100px] gap-2 items-center relative justify-center py-2 text-theme-base rounded-md disabled:bg-theme-disabled",
                variant === "recording" ? "bg-theme-primary-light" : "bg-theme-primary",
                className
            )}
        >
            <div className="w-2" />
            <MicrophoneIcon className="h-4 w-4" />
            <div className="w-2">
                {variant === "recording" && (
                    <div className="relative w-2 h-2 rounded-full bg-red-500">
                        <div className="absolute animate-ping rounded-full bg-red-500 opacity-70 inset-0" />
                    </div>
                )}
            </div>
        </button>
    );
};
export default RecordButton;
