import React, {Component, PropTypes} from 'react'
import {Redirect} from 'react-router-dom'

class Login extends Component {

    handleClick(event) {
        const username = this.refs.username
        const password = this.refs.password
        const creds = {username: username.value.trim(), password: password.value.trim()}
        this.props.onLoginClick(creds)
        event.preventDefault()
    }

    render() {
        const {errorMessage, isAuthenticated} = this.props

        return (
            <div className="container">
                <div className="wrapper">
                    <form name="Registration_Form" className="form-signin">
                        <div>
                            <h3 className="form-signin-heading">Welcome Back! Please Sign In</h3>
                            <input type="text" className="form-control" ref="username" placeholder="Username"
                                   autoFocus=""/>
                            <input type="password" className="form-control" ref="password" placeholder="Password"/>
                            <button onClick={(event) => this.handleClick(event)}
                                    className="btn btn-lg btn-primary btn-block">
                                Login
                            </button>
                            <a href="#/register">Register</a> - <a href="#">Forgot Password</a>
                            {errorMessage &&
                            <p>{errorMessage}</p>
                            }
                        </div>
                    </form>
                </div>
                {isAuthenticated &&
                <Redirect to="/"/>
                }
            </div>
        )
    }
}

Login.propTypes = {
    onLoginClick: PropTypes.func.isRequired,
    errorMessage: PropTypes.string,
    isAuthenticated: PropTypes.bool.isRequired
}

export default Login