import React from 'react';
import { Typography } from 'antd';
import { BankTwoTone } from '@ant-design/icons';
import { TitleStyle, IconStyle } from './Logo.styles';

const Logo = () => {
    const titleProps = {
        level: 4,
        style: TitleStyle
    }

    const iconProps = {
        twoToneColor: '#fff',
        style: IconStyle
    }

    return (
        <Typography.Title {...titleProps}>
            <BankTwoTone {...iconProps} />
            Bank
        </Typography.Title>
    );
};

export default Logo;
