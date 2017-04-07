import React, {Component, PropTypes} from 'react'
import {connect} from 'react-redux'
import {addTweet} from '../actions/twitterActions'

class TweetInput extends Component {

    constructor(props) {
        super(props)
        this.state = {text: ''}

        this.handleChange = this.handleChange.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
    }

    handleChange(event) {
        this.setState({text: event.target.value})
    }

    handleSubmit(event) {
        const text = this.state.text.trim()
        this.props.dispatch(addTweet(text))
        this.setState({text: ''})
    }

    render() {
        return (
            <form onSubmit={this.handleSubmit}>
                <input
                    className="form-control"
                    type="text"
                    placeholder="Compose new tweet..."
                    value={this.state.text}
                    onChange={this.handleChange}/>
            </form>
        )
    }
}

export default connect()(TweetInput)