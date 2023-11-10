import {Menu} from "@headlessui/react";
import clsx from "clsx";

export type Language = {
    name: string;
    tag: string;
    icon: string;
};
export const supportedLanguages: Record<Language["tag"], Language> = {
    "it-IT": {
        name: "Italian",
        tag: "it-IT",
        icon: "🇮🇹"
    },
    "pl-PL": {name: "Polish", tag: "pl-PL", icon: "🇵🇱"},
    "en-US": {name: "English", tag: "en-US", icon: "🇺🇸"},
    "es-ES": {name: "Spanish", tag: "es-ES", icon: "🇪🇸"},
    "de-DE": {name: "German", tag: "de-DE", icon: "🇩🇪"}
};

interface SelectLanguageProps {
    className?: string;
    currentLanguage: Language | null;
    onLanguageChange: (language: Language) => void;
}

export default function LanguageSelector({
    className,
    currentLanguage,
    onLanguageChange
}: SelectLanguageProps) {
    return (
        <Menu
            as="div"
            className={clsx(
                "border border-theme-secondary relative flex h-fit items-center rounded",
                className
            )}
        >
            <Menu.Button className="w-full px-4 py-2 h-full flex items-center">
                {currentLanguage ? (
                    <div className="flex gap-2">
                        <span className="text-lg leading-6">{currentLanguage.icon}</span>
                        <span className="text-base">{currentLanguage.name}</span>
                    </div>
                ) : (
                    <span>Select language</span>
                )}
            </Menu.Button>
            <Menu.Items className="absolute bg-theme-inverted top-full border border-theme-secondary overflow-hidden mt-2 w-full rounded left-0">
                {Object.values(supportedLanguages).map((language) => (
                    <Menu.Item
                        key={language.tag}
                        className={clsx(
                            "cursor-pointer flex gap-1 py-2 px-2",
                            currentLanguage?.tag === language.tag
                                ? "bg-gray-200"
                                : "bg-theme-inverted"
                        )}
                        as="div"
                        onClick={() => onLanguageChange(language)}
                    >
                        <span className="text-lg">{language.icon}</span>
                        <span className="text-base">{language.name}</span>
                    </Menu.Item>
                ))}
            </Menu.Items>
        </Menu>
    );
}
