import React, {Component, PropTypes} from 'react'
import {HashRouter, Route, Redirect} from 'react-router-dom'
import App from './../App';
import Register from './pages/register'
import AppNavbar from './navbar'
import Login from './pages/login'
import {loginUser} from '../actions/twitterActions'
import {connect} from 'react-redux'

class Root extends Component {
    render() {
        const {dispatch, isAuthenticated, errorMessage} = this.props
        return(
            <div>
                <AppNavbar dispatch={dispatch} isAuthenticated={isAuthenticated}/>
                <HashRouter>
                    <div>
                        <Route path="/register" component={Register}/>
                        <Route path="/login" render={(props) => (
                            <Login isAuthenticated={isAuthenticated} errorMessage={errorMessage} onLoginClick={ creds => dispatch(loginUser(creds))}/>
                        )}/>
                        <Route exact path="/" render={() => (
                            !isAuthenticated ? (
                                <Redirect to="/login"/>
                            ) : (
                                <App/>
                            )
                        )}/>
                    </div>
                </HashRouter>
            </div>
        )
    }
}
Root.propTypes = {
    dispatch: PropTypes.func.isRequired,
    isAuthenticated: PropTypes.bool.isRequired,
    errorMessage: PropTypes.string
}

function mapStateToProps(state) {

    const {auth} = state
    const {isAuthenticated, errorMessage} = auth

    return {
        isAuthenticated,
        errorMessage
    }
}


export default connect(mapStateToProps)(Root);