import {Language} from "../components/LanguageSelector";

export type User = {
    id: string;
    username: string;
    language: Language;
};

export type ChatMessage = {
    from_user: User;
    translation: string | null;
    transcript: string;
    confidence: number;
};
