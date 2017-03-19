import React, {Component, PropTypes} from 'react'
import {Media, ListGroupItem} from 'react-bootstrap'

class Tweet extends Component {
    render() {
        return (
            <ListGroupItem>
                <Media>
                    <Media.Body>
                        <Media.Heading>{this.props.user.name}</Media.Heading>
                        <p>{this.props.text}</p>
                    </Media.Body>
                </Media>
            </ListGroupItem>
        )
    }
}

Tweet.propTypes = {
    user: PropTypes.shape({
        name: PropTypes.string.isRequired
    }).isRequired,
    text: PropTypes.string.isRequired
}

export default Tweet