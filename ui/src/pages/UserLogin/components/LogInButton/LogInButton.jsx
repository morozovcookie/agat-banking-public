import React from 'react';
import { Button, Form } from 'antd';
import { LoginButtonFormStyle } from './LoginButton.styles'

const LogInButton = () => {
    const formProps = {
        style: LoginButtonFormStyle,
    }

    const buttonProps = {
        type: 'primary',
        htmlType: 'submit'
    }

    return (
        <Form.Item {...formProps} >
            <Button {...buttonProps} >
                Log in
            </Button>
        </Form.Item>
    );
};

export default LogInButton;
