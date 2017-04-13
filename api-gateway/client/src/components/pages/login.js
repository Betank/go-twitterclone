import React, {Component, PropTypes} from 'react'

class Login extends Component {

    handleClick(event) {
        const username = this.refs.username
        const password = this.refs.password
        const creds = {username: username.value.trim(), password: password.value.trim()}
        this.props.onLoginClick(creds)
        event.preventDefault()
    }

    render() {
        const {errorMessage} = this.props

        return (
            <div>
                <h3 className="form-signin-heading">Welcome Back! Please Sign In</h3>
                <input type="text" className="form-control" ref="username" placeholder="Username"
                       autoFocus=""/>
                <input type="password" className="form-control" ref="password" placeholder="Password"/>
                <button onClick={(event) => this.handleClick(event)} className="btn btn-lg btn-primary btn-block">
                    Login
                </button>
                {errorMessage &&
                <p>{errorMessage}</p>
                }
            </div>

        )
    }
}

Login.propTypes = {
    onLoginClick: PropTypes.func.isRequired,
    errorMessage: PropTypes.string
}

export default Login