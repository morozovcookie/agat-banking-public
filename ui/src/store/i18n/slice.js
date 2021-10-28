import { createSlice } from '@reduxjs/toolkit';
import { enUS, ruRU } from '@i18n';
import { initialState } from './initialState';

export const slice = createSlice({
    name: 'i18n',
    initialState: initialState,
    reducers: {
        switchOnEnglishLocale: (state) => {
            state.locale = enUS;
        },
        switchOnRussianLocale: (state) => {
            state.locale = ruRU;
        }
    }
});
