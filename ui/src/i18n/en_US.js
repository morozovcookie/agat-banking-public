import enUS from 'antd/lib/locale/en_US';
import { version } from './version';

const localeValues = {
    version: version,

    ...enUS,

    LogIn: {
        cardMetaTitle: 'Log in',
        cardMetaDescription: 'Log in into your account',
        usernameInputLabel: 'Username or email',
        passwordInputLabel: 'Password',
        logInButtonText: 'Log in',
        forgotPasswordLinkText: 'Forgot password?',
        footerText: 'Unable to access your account? To request an account, please contact your site administrators.',
    }
};

export default localeValues;
