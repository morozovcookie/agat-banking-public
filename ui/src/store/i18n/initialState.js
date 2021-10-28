import { defaultLocale } from '@i18n';

// TODO: по умолчанию берем defaultLocale, но если что-то есть в браузере, то тянет оттуда
export const initialState = {
    locale: defaultLocale,
};
