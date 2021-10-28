import { defaultLocale } from '@i18n';
import store from 'store';

const initialState = {
    locale: defaultLocale,
};

let storedLocale = store.get('locale');
if (storedLocale === undefined) {
    storedLocale = initialState.locale;
    store.set('locale', storedLocale);
}

initialState.locale = storedLocale;

export { initialState }
