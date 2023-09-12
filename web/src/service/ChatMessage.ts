import {object, string} from "yup";

export const eventSchema = object({type: string().required(), payload: object().required()});

export const errorPayloadSchema = object({
    message: string()
});
