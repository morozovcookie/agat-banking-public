import React, { useEffect, useState } from 'react';
import { Form, Input } from 'antd';
import { useSelector } from 'react-redux';
import { localeSelector } from '@store/i18n/selector';

const PasswordInput = ({name, ...rest}) => {
    const locale = useSelector(localeSelector);
    const [currentLocale, setCurrentLocale] = useState(locale.LogIn);

    useEffect(() => setCurrentLocale(locale.LogIn), [locale]);

    const formItemProps = {
        ...rest,
        required: true,
        label: currentLocale.passwordInputLabel,
        name: name,
    }

    return (
        <Form.Item {...formItemProps}>
            <Input.Password />
        </Form.Item>
    );
};

export default PasswordInput;
