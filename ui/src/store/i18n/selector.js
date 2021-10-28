import { createSelector } from '@reduxjs/toolkit';

export const localeSelector = createSelector(
    state => state.i18n.locale,
    locale => locale
);