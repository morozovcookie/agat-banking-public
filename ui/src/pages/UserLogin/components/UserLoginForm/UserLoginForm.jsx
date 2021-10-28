import React, { useEffect, useState } from 'react';
import { Card, Form, Divider, Typography } from 'antd';
import { useSelector } from 'react-redux';
import {localeSelector} from '@store/i18n/selector';
import { UserLoginFormCardStyle, UserLoginFormStyle, ForgotPasswordStyles } from './UserLoginForm.styles';
import { UserNameInput, PasswordInput, LogInButton } from './../';

const UserLoginForm = () => {
    const locale = useSelector(localeSelector);
    const [currentLocale, setCurrentLocale] = useState(locale.LogIn);

    useEffect(() => setCurrentLocale(locale.LogIn), [locale]);

    const cardMetaProps = {
        title: currentLocale.cardMetaTitle,
        description: currentLocale.cardMetaDescription,
    }

    const submitForm = (values: any) => {
        console.log(values);
    };

    const formProps = {
        layout: 'vertical',
        requiredMark: false,
        style: UserLoginFormStyle,
        onFinish: submitForm,
    }

    const forgotPasswordProps = {
        href: '',
        style: ForgotPasswordStyles
    }

    return (
        <Card style={UserLoginFormCardStyle}>
            <Card.Meta {...cardMetaProps} />
            <Form {...formProps}>
                <UserNameInput name='username' />
                <PasswordInput name='password' />

                <Form.Item>
                    <LogInButton />
                    <Typography.Link {...forgotPasswordProps}>{currentLocale.forgotPasswordLinkText}</Typography.Link>
                </Form.Item>
            </Form>
            <Divider />
            <Typography.Text type='secondary'>{currentLocale.footerText}</Typography.Text>
        </Card>
    );
};

export default UserLoginForm;
