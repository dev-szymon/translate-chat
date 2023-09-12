import {Menu} from "@headlessui/react";
import clsx from "clsx";

export type Language = {
    name: string;
    tag: string;
    icon: string;
};
export const supportedLanguages: Language[] = [
    {
        name: "Italian",
        tag: "it-IT",
        icon: "🇮🇹"
    },
    {name: "Polish", tag: "pl-PL", icon: "🇵🇱"},
    {name: "English", tag: "en-US", icon: "🇺🇸"},
    {name: "Spanish", tag: "es-ES", icon: "🇪🇸"},
    {name: "German", tag: "de-DE", icon: "🇩🇪"}
];

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
        <Menu as="div" className={clsx("", className)}>
            <Menu.Button>
                {currentLanguage ? (
                    <div className="flex gap-1">
                        <span className="text-lg">{currentLanguage.icon}</span>
                        <span className="text-base">{currentLanguage.name}</span>
                    </div>
                ) : (
                    <span>Select language</span>
                )}
            </Menu.Button>
            <Menu.Items>
                {supportedLanguages.map((language) => (
                    <Menu.Item
                        key={language.tag}
                        className="cursor-pointer flex gap-1"
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
