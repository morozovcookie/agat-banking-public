import ruRU from 'antd/lib/locale/ru_RU';
import { version } from './version';

const localeValues = {
    version: version,

    ...ruRU,

    LogIn: {
        cardMetaTitle: 'Авторизация',
        cardMetaDescription: 'Войдите в свою учетную запись',
        usernameInputLabel: 'Имя пользователя или email',
        passwordInputLabel: 'Пароль',
        logInButtonText: 'Авторизация',
        forgotPasswordLinkText: 'Забыли пароль?',
        footerText: 'Еще не зарегистрированы? Чтобы получить аккаунт, пожалуйста, обратитесь к администраторам.',
    },
};

export default localeValues;
