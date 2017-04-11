import React from 'react';
import ReactDOM from 'react-dom';
import {Router, Route, IndexRoute} from 'react-router';
import createHistory from 'history/createBrowserHistory'
import {Provider} from 'react-redux'
import configureStore from './configureStore';
import App from './App';
import LoginPage from './components/pages/login'
import './index.css';

let store = configureStore()

ReactDOM.render(
    <Provider store={store}>
        <App />
    </Provider>, document.getElementById('root'));
