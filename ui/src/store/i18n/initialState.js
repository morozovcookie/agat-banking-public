import { defaultLocale, version as lastVersion } from '@i18n';
import store from 'store';

const initialState = {
    locale: defaultLocale,
};

const storeLocale = (initial, stored) => {
    stored = initial.locale;
    store.set('locale', stored);
}

let storedLocale = store.get('locale');
if (storedLocale === undefined) {
    storeLocale(initialState, storedLocale);
}

if (storedLocale.version === undefined) {
    storeLocale(initialState, storedLocale);
}

if (storedLocale.version < lastVersion) {
    storeLocale(initialState, storedLocale);
}

initialState.locale = storedLocale;

export { initialState }
