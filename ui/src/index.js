import React, {useEffect, useState} from 'react';
import ReactDOM from 'react-dom';
import { ConfigProvider } from 'antd';
import { Router } from '@routes';
import { LocaleStore } from '@stores';
import './index.css';
import reportWebVitals from './reportWebVitals';


const LocaleProvider = () => {
    const [locale, setLocale] = useState('en_US')

    useEffect(() => {
        LocaleStore.subscribe(() => setLocale(LocaleStore.getState().locale));
    }, [])

    return (
        <ConfigProvider locale={locale}>
            <Router />
        </ConfigProvider>
    );
};

ReactDOM.render(
  <React.StrictMode>
      <LocaleProvider />
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
