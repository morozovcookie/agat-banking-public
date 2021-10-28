import { createSlice } from '@reduxjs/toolkit';
import { enUS, ruRU } from '@i18n';
import store from 'store';
import { initialState } from './initialState';

const storeWrapper = fn => {
    return state => {
        fn(state);
        store.set('locale', state.locale);
    }
}

export const slice = createSlice({
    name: 'i18n',
    initialState: initialState,
    reducers: {
        switchOnEnglishLocale: storeWrapper(state => state.locale = enUS),
        switchOnRussianLocale: storeWrapper(state => state.locale = ruRU),
    }
});
