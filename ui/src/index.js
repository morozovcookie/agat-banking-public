import './index.css';

import React, {useEffect, useState} from 'react';
import ReactDOM from 'react-dom';
import {Provider as ReduxProvider, useSelector} from 'react-redux';
import { Router } from '@routes';
import { store } from '@store';
import { localeSelector } from '@store/i18n/selector';
import reportWebVitals from './reportWebVitals';
import {ConfigProvider} from 'antd';

const Root = () => {
    const locale = useSelector(localeSelector);
    const [currentLocale, setCurrentLocale] = useState(locale);

    useEffect(() => setCurrentLocale(locale), [locale])

    return (
        <ConfigProvider locale={currentLocale}>
            <Router />
        </ConfigProvider>
    );
};

ReactDOM.render(
    <React.StrictMode>
        <ReduxProvider store={store}>
            <Root />
        </ReduxProvider>
    </React.StrictMode>,
    document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
