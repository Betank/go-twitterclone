import React from 'react';
import ReactDOM from 'react-dom';
import {HashRouter, Route, Link} from 'react-router-dom'
import {Provider} from 'react-redux'
import configureStore from './configureStore';
import Root from './components/root'
import './index.css';

let store = configureStore()

ReactDOM.render(
    <Provider store={store}>
        <Root/>
    </Provider>, document.getElementById('root'));
