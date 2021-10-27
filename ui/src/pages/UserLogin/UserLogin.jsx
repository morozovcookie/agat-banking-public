import React from 'react';
import { Layout, List, Button } from 'antd';
import { LocaleStore } from '@stores';
import { UserLoginForm, HeaderLogo } from './components';
import { LayoutStyle, HeaderStyle, ContentStyle } from './UserLogin.styles';

const UserLogin = () => {
    const items = [
        {
            locale: 'en_US',
            caption: 'English'
        },
        {
            locale: 'ru_RU',
            caption: 'Русский'
        }
    ]

    const onChangeLocale = (locale: String) => {
        return (event) => {
            LocaleStore.dispatch({type: locale});
        }
    }

    const listProps = {
        grid: {
            gutter: 12,
            columns: items.length,
        },
        dataSource: items,
        renderItem: item => (
            <List.Item style={{margin: 0}}>
                <Button type='link' onClick={onChangeLocale(item.locale)}>
                    {item.caption}
                </Button>
            </List.Item>
        )
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
