import {object, string, array, number} from "yup";
import {Language} from "../components/LanguageSelector";

export type User = {
    id: string;
    username: string;
    language: Language;
};

export type ChatMessage = {
    sender: User;
    translation: string | null;
    transcript: string;
    confidence: number;
};

export const eventSchema = object({type: string().required(), payload: object().required()});

export const errorPayloadSchema = object({
    message: string()
});

export const userJoinedRoomPayloadSchema = object({
    userId: string().required(),
    username: string().required(),
    language: string().required(),
    roomId: string().required(),
    roomName: string().required(),
    users: array()
        .required()
        .of(
            object({
                id: string().required(),
                username: string().required(),
                language: string().required()
            })
        )
});

export const newMessagePayloadSchema = object({
    transcript: string().required(),
    confidence: number().required(),
    translation: string().nullable(),
    userId: string(),
    roomId: string()
});
