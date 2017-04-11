import React, {Component, PropTypes} from 'react'
import {Navbar} from 'react-bootstrap'
import { loginUser, logoutUser } from '../actions/twitterActions'
import Login from './pages/login'
import Logout from './pages/logout'

class AppNavbar extends Component {
    render() {
        const {dispatch, isAuthenticated, errorMessage} = this.props

        return (
            <Navbar inverse collapseOnSelect fixedTop fluid>
                <div className='container-fluid'>
                    <Navbar.Header>
                        <Navbar.Brand>
                            <a href="#">GoTwitterClone</a>
                        </Navbar.Brand>
                        <Navbar.Toggle />
                    </Navbar.Header>
                    <form className="form-inline">
                        {!isAuthenticated &&
                        <Login
                            errorMessage={errorMessage}
                            onLoginClick={ creds => dispatch(loginUser(creds)) }
                        />
                        }

                        {isAuthenticated &&
                        <Logout onLogoutClick={() => dispatch(logoutUser())} />
                        }
                    </form>
                </div>
            </Navbar>
        )
    }
}

AppNavbar.propTypes = {
    dispatch: PropTypes.func.isRequired,
    isAuthenticated: PropTypes.bool.isRequired,
    errorMessage: PropTypes.string
}

export default AppNavbar