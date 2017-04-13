import React, {Component, PropTypes} from 'react';
import {connect} from 'react-redux'
import {Grid, Row, Col} from 'react-bootstrap'
import TweetList from './components/tweetList'
import Profile from './components/profile'
import AppNavbar from './components/navbar'
import {loginUser} from './actions/twitterActions'
import Login from './components/pages/login'

class App extends Component {
    render() {
        const {dispatch, isAuthenticated, errorMessage} = this.props
        return (
            <div>
                <AppNavbar dispatch={dispatch} isAuthenticated={isAuthenticated} errorMessage={errorMessage}/>
                {!isAuthenticated &&
                <div className="container">
                    <div className="wrapper">
                        <form name="Login_Form" className="form-signin">
                            <Login
                                errorMessage={errorMessage}
                                onLoginClick={ creds => dispatch(loginUser(creds)) }
                            />
                        </form>
                    </div>
                </div>
                }
                { isAuthenticated &&
                <Grid fluid>
                    <Row className="show-grid">
                        <Col sm={3} md={2}>
                            <Profile/>
                        </Col>
                        <Col sm={9} md={10}>
                            <TweetList/>
                        </Col>
                    </Row>
                </Grid>
                }
            </div>
        );
    }
}
App.propTypes = {
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


export default connect(mapStateToProps)(App);
