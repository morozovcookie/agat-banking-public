import { combineReducers } from 'redux';
import i18n from './i18n/reducer';

const rootReducer = combineReducers({
    i18n,
});

export default rootReducer;
