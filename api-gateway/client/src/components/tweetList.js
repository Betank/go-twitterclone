import React, {Component, PropTypes} from 'react'
import {connect} from 'react-redux'
import {Panel, ListGroup} from 'react-bootstrap'
import Tweet from './tweet'
import {fetchTweets} from '../actions/twitterActions'

class TweetList extends Component {

    componentDidMount() {
        this.props.dispatch(fetchTweets())
    }

    render() {
        return (
            <Panel collapsible defaultExpanded header="Tweets">
                <ListGroup fill>
                    {this.props.tweets.map(tweet =>
                        <Tweet key={tweet.id} {...tweet}
                        />
                    )}
                </ListGroup>
            </Panel>
        )
    }
}

TweetList.propTypes = {
    tweets: PropTypes.arrayOf(PropTypes.shape({
        id: PropTypes.string.isRequired,
        user: PropTypes.shape({
            name: PropTypes.string
        }).isRequired,
        text: PropTypes.string.isRequired
    }).isRequired).isRequired
}

const mapStateToProps = (state) => {
    return {
        tweets: state.tweets
    }
}


export default connect(mapStateToProps)(TweetList)