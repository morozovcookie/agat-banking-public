import { createStore } from 'redux';

const LocaleStore = createStore((state = {locale: 'en_US'}, action) => {
    const newLocale = action.type;

    if (newLocale === 'en_US') {
        return {locale: 'en_US'}
    }

    if (newLocale === 'ru_RU') {
        return {locale: 'ru_RU'}
    }

    return state
});

export default LocaleStore;
