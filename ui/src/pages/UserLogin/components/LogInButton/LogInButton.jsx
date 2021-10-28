import React, { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { Button, Form } from 'antd';
import { localeSelector } from '@store/i18n/selector';
import { LoginButtonFormStyle } from './LoginButton.styles'

const LogInButton = () => {
    const locale = useSelector(localeSelector);
    const [currentLocale, setCurrentLocale] = useState(locale.LogIn);

    useEffect(() => setCurrentLocale(locale.LogIn), [locale]);

    const formProps = {
        style: LoginButtonFormStyle,
    }

    const buttonProps = {
        type: 'primary',
        htmlType: 'submit',
    }

    return (
        <Form.Item {...formProps}>
            <Button {...buttonProps}>
                {currentLocale.logInButtonText}
            </Button>
        </Form.Item>
    );
};

export default LogInButton;
