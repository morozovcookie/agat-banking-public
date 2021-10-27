import React from 'react';
import { Form, Input } from 'antd';

const PasswordInput = ({name, ...rest}) => {
    const formItemProps = {
        ...rest,
        required: true,
        label: 'Password',
        name: name,
    }

    return (
        <Form.Item {...formItemProps}>
            <Input.Password />
        </Form.Item>
    );
};

export default PasswordInput;
