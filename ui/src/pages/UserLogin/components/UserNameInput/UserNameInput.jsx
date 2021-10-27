import React from 'react';
import { Form, Input } from 'antd';

const UserNameInput = ({name, ...rest}) => {
    const formItemProps = {
        ...rest,
        required: true,
        label: 'Email Address / Username',
        name: name,
    }

    return (
        <Form.Item {...formItemProps} >
            <Input />
        </Form.Item>
    );
};

export default UserNameInput;
