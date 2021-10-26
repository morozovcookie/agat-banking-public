import React from 'react';
import { BrowserRouter as ReactRouter, Route, Switch } from 'react-router-dom';
import PrivateRoute from './PrivateRoute';
import { App, UserLogin } from '@pages';

export const Router = () => {
    return (
        <ReactRouter>
            <Switch>
                <PrivateRoute exact path='/' component={App} />
                <Route exact path='/login' component={UserLogin} />
            </Switch>
        </ReactRouter>
    );
}
