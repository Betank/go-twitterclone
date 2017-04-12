import React, {Component, PropTypes} from 'react'
import {Navbar} from 'react-bootstrap'
import {logoutUser} from '../actions/twitterActions'
import Logout from './pages/logout'

class AppNavbar extends Component {
    render() {
        const {dispatch, isAuthenticated, errorMessage} = this.props

        return (
            <div>
                <Navbar inverse collapseOnSelect fixedTop fluid>
                    <div className='container-fluid'>
                        <Navbar.Header>
                            <Navbar.Brand>
                                <a href="#">GoTwitterClone</a>
                            </Navbar.Brand>
                            <Navbar.Toggle />
                        </Navbar.Header>
                        <form className="form-inline">
                            {isAuthenticated &&
                            <Logout onLogoutClick={() => dispatch(logoutUser())}/>
                            }
                        </form>
                    </div>
                </Navbar>
            </div>
        )
    }
}

AppNavbar.propTypes = {
    dispatch: PropTypes.func.isRequired,
    isAuthenticated: PropTypes.bool.isRequired,
    errorMessage: PropTypes.string
}

export default AppNavbar