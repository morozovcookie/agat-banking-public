import {configureStore} from '@reduxjs/toolkit';
import rootReducer from './rootReducer';
import { preloadedState } from './preloadedState';

const store = configureStore({
    reducer: rootReducer,
    preloadedState
});

export default store;
