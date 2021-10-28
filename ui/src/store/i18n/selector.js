import { createSelector } from '@reduxjs/toolkit';

// TODO: значение сбрасывается при каждом нажатии F5 - нужно хранить выбранный язык в сессии или тупо писать в браузер.
export const localeSelector = createSelector(
    state => state.i18n.locale,
    locale => locale
);