import React from 'react';
import { Layout } from 'antd';
import { UserLoginForm, Logo } from './components';
import { LayoutStyle, HeaderStyle, ContentStyle } from './UserLogin.styles';

const UserLogin = () => {
    return (
        <Layout style={LayoutStyle}>
            <Layout.Header style={HeaderStyle}>
                <Logo />
            </Layout.Header>
            <Layout.Content style={ContentStyle}>
                <UserLoginForm />
            </Layout.Content>
            <Layout.Footer />
        </Layout>
    );
};

export default UserLogin;
