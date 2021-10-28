import React from 'react';
import { useDispatch } from 'react-redux';
import { Layout, List, Button } from 'antd';
import { switchOnEnglishLocale, switchOnRussianLocale } from '@store/i18n/actions';
import { UserLoginForm, HeaderLogo } from './components';
import { LayoutStyle, HeaderStyle, ContentStyle } from './UserLogin.styles';

const UserLogin = () => {
    const dispatch = useDispatch();

    const onChangeLocale = (dispatch, switcher) => () => dispatch(switcher());

    const items = [
        <Button type='link' onClick={onChangeLocale(dispatch, switchOnEnglishLocale)}>English</Button>,
        <Button type='link' onClick={onChangeLocale(dispatch, switchOnRussianLocale)}>Русский</Button>
    ]

    const listProps = {
        grid: {
            columns: items.length,
        },
        dataSource: items,
        renderItem: item => <List.Item style={{margin: 0}}>{item}</List.Item>
    }


    return (
        <Layout style={LayoutStyle}>
            <Layout.Header style={HeaderStyle}>
                <HeaderLogo />
            </Layout.Header>
            <Layout.Content style={ContentStyle}>
                <UserLoginForm />
            </Layout.Content>
            <Layout.Footer>
                <List {...listProps} style={{width: '203px', left: '50%', marginLeft: '-101.5px'}} />
            </Layout.Footer>
        </Layout>
    );
};

export default UserLogin;
