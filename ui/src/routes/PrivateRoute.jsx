import React from 'react';
import { Route, Redirect } from 'react-router-dom';

const PrivateRoute = ({ component: Component, ...rest }) => {
    const fakeAuth = {
        isAuthenticated: false
    }

    const render = () => {
        return fakeAuth.isAuthenticated === true ? <Component /> : <Redirect to='/login' />;
    };

    return <Route {...rest} render={render} />;
};

export default PrivateRoute;
