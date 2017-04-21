import React, {Component, PropTypes} from 'react'
import {Redirect} from 'react-router-dom'

class Register extends Component {

    handleClick(event) {
        const username = this.refs.username
        const password = this.refs.password
        const creds = {username: username.value.trim(), password: password.value.trim()}
        this.props.onRegistrationClick(creds)
        event.preventDefault()
    }

    render() {
        const {errorMessage, isAuthenticated} = this.props
        return (
            <div className="container">
                <div className="wrapper">
                    <form name="Registration_Form" className="form-signin">
                        <div>
                            <h3 className="form-signin-heading">Registration</h3>
                            <input type="text" className="form-control" ref="username" placeholder="Username"
                                   autoFocus=""/>
                            <input type="password" className="form-control" ref="password" placeholder="Password"/>
                            <button onClick={(event) => this.handleClick(event)}
                                    className="btn btn-lg btn-primary btn-block">
                                Register
                            </button>
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

Register.propTypes = {
    onRegistrationClick: PropTypes.func.isRequired,
    errorMessage: PropTypes.string,
    isAuthenticated: PropTypes.bool.isRequired
}

export default Register;