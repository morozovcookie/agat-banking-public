import React from 'react';
import { Card, Form, Divider, Typography } from 'antd';
import { UserLoginFormCardStyle, UserLoginFormStyle, ForgotPasswordStyles } from './UserLoginForm.styles';
import { UserNameInput, PasswordInput, LogInButton } from './../';

const UserLoginForm = () => {
    const onFinish = (values: any) => {
        console.log(values);
    };

    const cardMetaProps = {
        title: 'Log in',
        description: 'Log in into your account'
    }

    const formProps = {
        layout: 'vertical',
        requiredMark: false,
        style: UserLoginFormStyle,
        onFinish: onFinish,
    }

    const forgotPasswordProps = {
        href: '',
        style: ForgotPasswordStyles
    }

    return (
        <Card style={UserLoginFormCardStyle}>
            <Card.Meta {...cardMetaProps} />
            <Form {...formProps} >
                <UserNameInput name='username' />
                <PasswordInput name='password' />

                <Form.Item>
                    <LogInButton />
                    <Typography.Link {...forgotPasswordProps} >
                        Forgot password?
                    </Typography.Link>
                </Form.Item>
            </Form>
            <Divider />
            <Typography.Text type='secondary'>
                Unable to access your account? To request an account, please contact your site administrators.
            </Typography.Text>
        </Card>
    );
};

export default UserLoginForm;
